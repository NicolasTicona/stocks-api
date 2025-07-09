package stockssync

import (
	"fmt"
	"time"

	"github.com/nicolasticona/stocks-api/db"
	"github.com/nicolasticona/stocks-api/models"
)

type StocksSyncRepository struct {
}

func (repository *StocksSyncRepository) InsertStockRatings(items []map[string]interface{}) error {
	if err := db.DB.Exec("DELETE FROM stocks").Error; err != nil {
		return err
	}

	var stocks []models.Stock
	for _, item := range items {
		timeStr := item["time"].(string)
		parsedTime, err := time.Parse(time.RFC3339, timeStr)

		if err != nil {
			return err
		}

		stock := models.Stock{
			Ticker:      item["ticker"].(string),
			Company:     item["company"].(string),
			Action:      item["action"].(string),
			Brokerage:   item["brokerage"].(string),
			RatingFrom:  item["rating_from"].(string),
			RatingTo:    item["rating_to"].(string),
			TargetFrom:  item["target_from"].(float64),
			TargetTo:    item["target_to"].(float64),
			Time:        parsedTime,
			RatingScore: item["rating_score"].(float64),
			TargetScore: item["target_score"].(float64),
		}
		stocks = append(stocks, stock)
	}

	if err := db.DB.Create(&stocks).Error; err != nil {
		return err
	}

	fmt.Println("Successfully inserted stock ratings into the database.")

	return nil
}

func (r *StocksSyncRepository) InsertStockRecommendations() error {
	if err := db.DB.Exec("DELETE FROM stock_recommendations").Error; err != nil {
		return err
	}

	queryTopItems := `
        SELECT ticker, (target_score + rating_score) AS total_score
        FROM stocks
        WHERE time >= NOW() - INTERVAL '7 days'
        ORDER BY total_score DESC
        LIMIT 10;
    `

	var topItems []struct {
		Ticker     string
		TotalScore float64
	}

	rows := db.DB.Raw(queryTopItems).Scan(&topItems)

	if rows.Error != nil {
		return rows.Error
	}

	if len(topItems) == 0 {
		fmt.Println("No items found to insert into stock_recommendations.")
		return nil
	}

	var recommendations []models.StockRecommendation
	for _, item := range topItems {
		recommendation := models.StockRecommendation{
			Ticker:     item.Ticker,
			TotalScore: item.TotalScore,
			CreatedAt:  time.Now(),
		}
		recommendations = append(recommendations, recommendation)
	}

	if err := db.DB.Create(&recommendations).Error; err != nil {
		return err
	}

	fmt.Println("Successfully inserted stock recommendations into the database.")

	return nil
}
