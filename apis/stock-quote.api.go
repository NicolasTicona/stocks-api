package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/nicolasticona/stocks-api/models"
)

type StockQuoteAPI struct {
	CurrentPrice     float64 `json:"c"`
	Change           float64 `json:"d"`
	ChangePercentage float64 `json:"dp"`
	High             float64 `json:"h"`
	Low              float64 `json:"l"`
	Open             float64 `json:"o"`
	PreviousClose    float64 `json:"pc"`
}

func GetStockQuote(stock string) (models.CurrentQuote, error) {
	apiURL := fmt.Sprintf("https://finnhub.io/api/v1/quote?symbol=%s&token=%s", stock, os.Getenv("FINNHUB_API_KEY"))

	client := &http.Client{}
	resp, err := client.Get(apiURL)

	if err != nil {
		fmt.Printf("Error fetching quote for stock %s: %v\n", stock, err)
		return models.CurrentQuote{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Non-200 response for stock %s: %d\n", stock, resp.StatusCode)
		return models.CurrentQuote{}, fmt.Errorf("non-200 response for stock %s: %d", stock, resp.StatusCode)
	}

	var stockQuote StockQuoteAPI
	err = json.NewDecoder(resp.Body).Decode(&stockQuote)

	if err != nil {
		fmt.Printf("Error decoding quote for stock %s: %v\n", stock, err)
		return models.CurrentQuote{}, err
	}

	fmt.Println("Stock quote fetched successfully:", stockQuote)

	return models.CurrentQuote{
		CurrentPrice:     stockQuote.CurrentPrice,
		Change:           stockQuote.Change,
		ChangePercentage: stockQuote.ChangePercentage,
		High:             stockQuote.High,
		Low:              stockQuote.Low,
		Open:             stockQuote.Open,
		PreviousClose:    stockQuote.PreviousClose,
	}, nil
}
