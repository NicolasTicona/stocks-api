package models

import (
	"time"
)

type Stock struct {
	Id           int          `gorm:"primaryKey;unique;autoIncrement;not null" json:"id"`
	Ticker       string       `gorm:"type:varchar;not null" json:"ticker"`
	TargetFrom   float64      `gorm:"not null" json:"target_from"`
	TargetTo     float64      `gorm:"not null" json:"target_to"`
	Company      string       `gorm:"type:varchar;not null" json:"company"`
	Action       string       `gorm:"type:varchar;not null" json:"action"`
	Brokerage    string       `gorm:"type:varchar;not null" json:"brokerage"`
	RatingFrom   string       `gorm:"type:varchar;not null" json:"rating_from"`
	RatingTo     string       `gorm:"type:varchar;not null" json:"rating_to"`
	Time         time.Time    `gorm:"not null" json:"time"`
	RatingScore  float64      `gorm:"not null" json:"rating_score"`
	TargetScore  float64      `gorm:"not null" json:"target_score"`
	CurrentQuote CurrentQuote `gorm:"-" json:"currentQuote,omitempty"`
	TotalCount   uint16       `gorm:"null" json:"total_count"`
}
