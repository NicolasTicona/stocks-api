package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type StockQuote struct {
	CurrentPrice     float64 `json:"c"`
	Change           float64 `json:"d"`
	ChangePercentage float64 `json:"dp"`
	High             float64 `json:"h"`
	Low              float64 `json:"l"`
	Open             float64 `json:"o"`
	PreviousClose    float64 `json:"pc"`
}

func GetStockQuote(stock string) (StockQuote, error) {
	apiURL := fmt.Sprintf("https://finnhub.io/api/v1/quote?symbol=%s&token=%s", stock, os.Getenv("FINNHUB_API_KEY"))

	client := &http.Client{}
	resp, err := client.Get(apiURL)

	if err != nil {
		fmt.Printf("Error fetching quote for stock %s: %v\n", stock, err)
		return StockQuote{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Non-200 response for stock %s: %d\n", stock, resp.StatusCode)
		return StockQuote{}, fmt.Errorf("non-200 response for stock %s: %d", stock, resp.StatusCode)
	}

	var stockQuote StockQuote
	err = json.NewDecoder(resp.Body).Decode(&stockQuote)

	if err != nil {
		fmt.Printf("Error decoding quote for stock %s: %v\n", stock, err)
		return StockQuote{}, err
	}

	fmt.Println("Stock quote fetched successfully:", stockQuote)

	return stockQuote, nil
}
