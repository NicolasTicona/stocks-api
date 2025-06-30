package models

import (
	"fmt"
	"log"
	"time"

	"github.com/nicolasticona/stocks-api/db"
	"github.com/nicolasticona/stocks-api/utils"
)

type StockRecommendation struct {
	Id           int         `gorm:"primaryKey;unique;autoIncrement;not null" json:"id"`
	Ticker       string      `gorm:"type:varchar;not null" json:"ticker"`
	TotalScore   float64     `gorm:"not null" json:"total_score"`
	CreatedAt    time.Time   `gorm:"not null;defautl:current_timestamp" json:"created_at"`
	CurrentQuote utils.JSONB `gorm:"type:jsonb null" json:"currentQuote,omitempty"`
}

func InsertStockRecommendations() error {
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
		log.Println("No items found to insert into stock_recommendations.")
		return nil
	}

	var recommendations []StockRecommendation
	for _, item := range topItems {
		recommendation := StockRecommendation{
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
