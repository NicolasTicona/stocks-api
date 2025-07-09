package stocks

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/nicolasticona/stocks-api/apis"
	"github.com/nicolasticona/stocks-api/models"
	"github.com/nicolasticona/stocks-api/utils"
)

type StocksController struct {
	Repository *StocksRepository
}

type StocksResponse struct {
	Total      uint16         `json:"total"`
	Page       uint16         `json:"page"`
	Limit      uint16         `json:"limit"`
	TotalPages uint16         `json:"totalPages"`
	Stocks     []models.Stock `json:"stocks"`
}

func (controller *StocksController) GetStocks(w http.ResponseWriter, page uint16, limit uint16, sortBy string, filter string) (StocksResponse, error) {
	var stocks []models.Stock
	var totalStocks uint16

	cacheKey := fmt.Sprintf("STOCKS_PAGE_%d_LIMIT_%d_SORT_%s_FILTER_%s", page, limit, sortBy, filter)
	cachedResponse, err := utils.RedisGet(cacheKey)

	if err == nil {
		fmt.Println("Cache hit for key:", cacheKey)

		var cachedData StocksResponse

		err = json.Unmarshal([]byte(cachedResponse), &cachedData)
		if err != nil {
			return StocksResponse{
				Stocks: []models.Stock{},
			}, err
		}

		return cachedData, nil
	}

	stocks, err = controller.Repository.FindAll(page, limit, filter, sortBy)

	if err != nil {
		fmt.Printf("Failed to fetch stocks from the database: %v\n", err)
		return StocksResponse{
			Stocks: []models.Stock{},
		}, err
	}

	if len(stocks) == 0 {
		return StocksResponse{
			Stocks: []models.Stock{},
		}, nil
	}

	totalStocks = stocks[0].TotalCount

	if len(stocks) > 50 {
		fmt.Printf("Warning: totalStocks is greater than 50, which may lead to performance issues.\n")
		return StocksResponse{
			Stocks: []models.Stock{},
		}, errors.New("executing query with more than 50 stocks is not allowed")
	}

	var wg sync.WaitGroup
	for i, stock := range stocks {
		wg.Add(1)
		go func(i int, stock models.Stock) {
			defer wg.Done()
			value, err := utils.RedisGet("QUOTE_" + stock.Ticker)
			if err != nil {
				fmt.Printf("Cache miss for QUOTE_%s: %v\n", stock.Ticker, err)

				stockQuote, err := apis.GetStockQuote(stock.Ticker)
				if err != nil {
					fmt.Printf("Error fetching quote for stock %s: %v\n", stock.Ticker, err)
					return
				}

				quoteJSON, err := json.Marshal(stockQuote)

				if err != nil {
					fmt.Printf("Error marshalling quote for stock %s: %v\n", stock.Ticker, err)
					return
				}

				stocks[i].CurrentQuote = stockQuote

				err = utils.RedisSave("QUOTE_"+stock.Ticker, quoteJSON, time.Hour)

				if err != nil {
					fmt.Println("Error saving quote to cache:", err)
				}

				return
			}

			fmt.Println("Cache hit for QUOTE_" + stock.Ticker)
			var cacheResult models.CurrentQuote
			err = json.Unmarshal([]byte(value), &cacheResult)
			if err != nil {
				fmt.Printf("Error unmarshalling cache for QUOTE_%s: %v\n", stock.Ticker, err)
				return
			}
			stocks[i].CurrentQuote = cacheResult
		}(i, stock)
	}
	wg.Wait()

	totalPages := (totalStocks + limit - 1) / limit

	response := StocksResponse{
		Total:      totalStocks,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		Stocks:     stocks,
	}

	responseJSON, err := json.Marshal(response)

	if err != nil {
		fmt.Printf("Error marshalling response: %v\n", err)
		return StocksResponse{
			Stocks: []models.Stock{},
		}, err
	}

	if len(stocks) > 0 {
		utils.RedisSave(cacheKey, responseJSON, time.Minute*5)
	}

	return response, nil
}
