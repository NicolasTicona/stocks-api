package models

import (
	"time"
)

type StockRecommendation struct {
	Id           int          `gorm:"primaryKey;unique;autoIncrement;not null" json:"id"`
	Ticker       string       `gorm:"type:varchar;not null" json:"ticker"`
	TotalScore   float64      `gorm:"not null" json:"total_score"`
	CreatedAt    time.Time    `gorm:"not null;defautl:current_timestamp" json:"created_at"`
	CurrentQuote CurrentQuote `gorm:"-" json:"currentQuote,omitempty"`
}
