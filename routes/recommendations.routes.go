package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/nicolasticona/stocks-api/db"
	"github.com/nicolasticona/stocks-api/models"
)

func GetRecommendationsHandler(w http.ResponseWriter, r *http.Request) {
	var stocks []models.StockRecommendation
	rows := db.DB.Find(&stocks)
	w.Header().Set("Content-Type", "application/json")

	if rows.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to fetch stocks from the database",
		})
		return
	}

	if len(stocks) == 0 {
		stocks = []models.StockRecommendation{}
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

	response := map[string]interface{}{
		"stocks": stocks,
	}

	json.NewEncoder(w).Encode(response)
}
