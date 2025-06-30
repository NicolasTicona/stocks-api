package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/nicolasticona/stocks-api/apis"
	"github.com/nicolasticona/stocks-api/db"
	"github.com/nicolasticona/stocks-api/models"
)

func GetStocksHandler(w http.ResponseWriter, r *http.Request) {
	var stocks []models.Stock
	var totalStocks int64
	page := r.URL.Query().Get("page")
	pageLimit := r.URL.Query().Get("limit")
	filter := r.URL.Query().Get("filter")

	if page == "" {
		page = "0"
	}
	if pageLimit == "" {
		pageLimit = "10"
	}

	db.DB.Raw(`SELECT COUNT(*) FROM stocks WHERE ($1 = '' OR ticker = $1);`, filter).Scan(&totalStocks)

	if totalStocks == 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"total":      0,
			"page":       0,
			"limit":      0,
			"totalPages": 0,
			"stocks":     []models.Stock{},
		})
		return
	}

	rows := db.DB.Raw(`SELECT * FROM stocks
            WHERE ($1 = '' OR ticker = $1)
            ORDER BY time DESC
            LIMIT $2 OFFSET $3;`, filter, pageLimit, page).Scan(&stocks)

	w.Header().Set("Content-Type", "application/json")

	if rows.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to fetch stocks from the database",
		})
		return
	}

	client := &http.Client{}
	for i, stock := range stocks {
		apiURL := fmt.Sprintf("https://finnhub.io/api/v1/quote?symbol=%s&token=%s", stock.Ticker, os.Getenv("FINNHUB_API_KEY"))
		resp, err := client.Get(apiURL)
		if err != nil {
			fmt.Printf("Error fetching quote for stock %s: %v\n", stock.Ticker, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Non-200 response for stock %s: %d\n", stock.Ticker, resp.StatusCode)
			continue
		}

		var stockQuote struct {
			CurrentPrice     float64 `json:"c"`
			Change           float64 `json:"d"`
			ChangePercentage float64 `json:"dp"`
			High             float64 `json:"h"`
			Low              float64 `json:"l"`
			Open             float64 `json:"o"`
			PreviousClose    float64 `json:"pc"`
		}

		err = json.NewDecoder(resp.Body).Decode(&stockQuote)
		if err != nil {
			fmt.Printf("Error decoding quote for stock %s: %v\n", stock.Ticker, err)
			continue
		}

		// Map the response to the currentQuote field
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

	limit, _ := strconv.Atoi(pageLimit)
	pageNum, _ := strconv.Atoi(page)
	totalPages := 0

	if len(stocks) > 0 {
		totalPages = int((totalStocks + int64(limit) - 1) / int64(limit))
	} else {
		stocks = []models.Stock{}
	}

	// Success response structure
	response := map[string]interface{}{
		"total":      totalStocks,
		"page":       pageNum,
		"limit":      limit,
		"totalPages": totalPages,
		"stocks":     stocks,
	}

	json.NewEncoder(w).Encode(response)
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
