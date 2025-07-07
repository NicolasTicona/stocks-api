package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/nicolasticona/stocks-api/apis"
	"github.com/nicolasticona/stocks-api/db"
	"github.com/nicolasticona/stocks-api/models"
	"github.com/nicolasticona/stocks-api/utils"
)

func GetStocksHandler(w http.ResponseWriter, r *http.Request) {
	var stocks []models.Stock
	var totalStocks int64

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	filter := r.URL.Query().Get("filter")
	sortBy := r.URL.Query().Get("sortBy")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 0
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := page * limit

	cacheKey := fmt.Sprintf("STOCKS_PAGE_%d_LIMIT_%d_SORT_%s_FILTER_%s", page, limit, sortBy, filter)
	cachedResponse, err := utils.RedisGet(cacheKey)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")

		fmt.Println("Cache hit for key:", cacheKey)
		w.Write([]byte(cachedResponse))
		return
	}

	var sortBySql string

	switch sortBy {
	case "time":
		sortBySql = "time DESC"
	case "rating":
		sortBySql = "rating_score DESC"
	case "target":
		sortBySql = "target_to DESC"
	default:
		sortBySql = "time DESC"
	}

	rows := db.DB.Raw(fmt.Sprintf(`
		WITH total_count AS (
			SELECT COUNT(*) AS count FROM stocks WHERE ($1 = '' OR ticker = $1)
		)
		SELECT *, (SELECT count FROM total_count) AS total_count
		FROM stocks
		WHERE ($1 = '' OR ticker = $1)
		ORDER BY %s
		LIMIT $2 OFFSET $3;
	`, sortBySql), filter, limit, offset).Scan(&stocks)

	if rows.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to fetch stocks from the database",
		})
		return
	}

	if len(stocks) > 0 {
		println("Fetched stocks from the database:", stocks[0].TotalCount)
		totalStocks = stocks[0].TotalCount
	} else {
		stocks = []models.Stock{}
		totalStocks = 0
	}

	var wg sync.WaitGroup
	for i, stock := range stocks {
		wg.Add(1)
		go func(i int, stock models.Stock) {
			defer wg.Done()
			value, err := utils.RedisGet("QUOTE_" + stock.Ticker)
			if err == nil {
				var cacheResult map[string]interface{}
				if json.Unmarshal([]byte(value), &cacheResult) == nil {
					stocks[i].CurrentQuote = cacheResult
					return
				}
			}

			stockQuote, err := apis.GetStockQuote(stock.Ticker)
			if err == nil {
				stocks[i].CurrentQuote = map[string]interface{}{
					"currentPrice":     stockQuote.CurrentPrice,
					"change":           stockQuote.Change,
					"changePercentage": stockQuote.ChangePercentage,
					"high":             stockQuote.High,
					"low":              stockQuote.Low,
					"open":             stockQuote.Open,
					"previousClose":    stockQuote.PreviousClose,
				}

				quoteJSON, _ := json.Marshal(stocks[i].CurrentQuote)
				utils.RedisSave("QUOTE_"+stock.Ticker, quoteJSON, time.Hour)
			} else {
				stocks[i].CurrentQuote = map[string]interface{}{
					"currentPrice":     stockQuote.CurrentPrice,
					"change":           stockQuote.Change,
					"changePercentage": stockQuote.ChangePercentage,
					"high":             stockQuote.High,
					"low":              stockQuote.Low,
					"open":             stockQuote.Open,
					"previousClose":    stockQuote.PreviousClose,
				}
			}
		}(i, stock)
	}
	wg.Wait()

	totalPages := int((totalStocks + int64(limit) - 1) / int64(limit))

	response := map[string]interface{}{
		"total":      totalStocks,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
		"stocks":     stocks,
	}

	responseJSON, _ := json.Marshal(response)

	if len(stocks) > 0 {
		utils.RedisSave(cacheKey, responseJSON, time.Minute*5)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func SyncDbHandler(w http.ResponseWriter, r *http.Request) {
	stocks, err := apis.FetchAllStockRatings()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if err := models.InsertStockRatings(stocks); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if err := models.InsertStockRecommendations(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Return success response with all items
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Data synced successfully",
	})

}
