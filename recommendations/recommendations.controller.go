package recommendations

import (
	"fmt"
	"net/http"

	"github.com/nicolasticona/stocks-api/apis"
	"github.com/nicolasticona/stocks-api/models"
)

type RecommendationsController struct {
	Repository *RecommendationRepository
}

func (controller *RecommendationsController) GetRecommendations(w http.ResponseWriter) ([]models.StockRecommendation, error) {
	stocks, err := controller.Repository.FindAll()

	if err != nil {
		fmt.Printf("Failed to fetch stocks from the database: %v\n", err)
		return nil, err
	}

	if len(stocks) == 0 {
		return stocks, nil
	}

	for i, stock := range stocks {
		response, err := apis.GetStockQuote(stock.Ticker)

		if err != nil {
			fmt.Printf("Failed to fetch stock quote for %s: %v\n", stock.Ticker, err)
			return nil, err
		}

		stocks[i].CurrentQuote = response
	}

	return stocks, nil
}
