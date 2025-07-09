package stockssync

import (
	"errors"
	"fmt"

	"github.com/nicolasticona/stocks-api/apis"
)

type StocksSyncController struct {
	Repository *StocksSyncRepository
}

func (controller *StocksSyncController) InsertStockRatings() error {
	stocks, err := apis.FetchAllStockRatings()
	if err != nil {
		fmt.Println("Failed to fetch stock ratings:", err)
		return errors.New("failed to fetch stock ratings")
	}

	if err := controller.Repository.InsertStockRatings(stocks); err != nil {
		fmt.Println("Failed to insert stock ratings:", err)
		return errors.New("failed to insert stock ratings")
	}

	if err := controller.Repository.InsertStockRecommendations(); err != nil {
		fmt.Println("Failed to insert stock recommendations:", err)
		return errors.New("failed to insert stock recommendations")
	}

	return nil
}
