package models

import (
	"fmt"
	"time"

	"github.com/nicolasticona/stocks-api/db"
	"github.com/nicolasticona/stocks-api/utils"
)

type Stock struct {
	Id           int         `gorm:"primaryKey;unique;autoIncrement;not null" json:"id"`
	Ticker       string      `gorm:"type:varchar;not null" json:"ticker"`
	TargetFrom   float64     `gorm:"not null" json:"target_from"`
	TargetTo     float64     `gorm:"not null" json:"target_to"`
	Company      string      `gorm:"type:varchar;not null" json:"company"`
	Action       string      `gorm:"type:varchar;not null" json:"action"`
	Brokerage    string      `gorm:"type:varchar;not null" json:"brokerage"`
	RatingFrom   string      `gorm:"type:varchar;not null" json:"rating_from"`
	RatingTo     string      `gorm:"type:varchar;not null" json:"rating_to"`
	Time         time.Time   `gorm:"not null" json:"time"`
	RatingScore  float64     `gorm:"not null" json:"rating_score"`
	TargetScore  float64     `gorm:"not null" json:"target_score"`
	CurrentQuote utils.JSONB `gorm:"type:jsonb null" json:"currentQuote,omitempty"`
}

func InsertStockRatings(items []map[string]interface{}) error {
	if err := db.DB.Exec("DELETE FROM stocks").Error; err != nil {
		return err
	}

	var stocks []Stock
	for _, item := range items {
		timeStr := item["time"].(string)
		parsedTime, err := time.Parse(time.RFC3339, timeStr)

		if err != nil {
			return err
		}

		stock := Stock{
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
