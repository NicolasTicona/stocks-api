package recommendations

import (
	"github.com/nicolasticona/stocks-api/models"

	"github.com/nicolasticona/stocks-api/db"
)

type RecommendationRepository struct{}

func (r *RecommendationRepository) FindAll() ([]models.StockRecommendation, error) {
	var stocks []models.StockRecommendation
	rows := db.DB.Find(&stocks)

	if rows.Error != nil {
		return nil, rows.Error
	}

	if len(stocks) == 0 {
		stocks = []models.StockRecommendation{}
	}

	return stocks, nil
}
