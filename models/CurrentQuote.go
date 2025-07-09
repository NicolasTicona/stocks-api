package models

type CurrentQuote struct {
	CurrentPrice     float64 `json:"currentPrice"`
	Change           float64 `json:"change"`
	ChangePercentage float64 `json:"changePercentage"`
	High             float64 `json:"high"`
	Low              float64 `json:"low"`
	Open             float64 `json:"open"`
	PreviousClose    float64 `json:"previousClose"`
}
