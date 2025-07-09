package stocks

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/nicolasticona/stocks-api/db"
	"github.com/nicolasticona/stocks-api/models"
)

type StocksRepository struct {
}

func (r *StocksRepository) FindAll(page uint16, limit uint16, filter string, sortBy string) ([]models.Stock, error) {
	var sortBySql string
	var stocks []models.Stock

	if limit > 50 {
		return nil, errors.New("limit cannot be greater than 50")
	}

	isValidFilter := regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(filter)
	if filter != "" && !isValidFilter {
		return nil, errors.New("filter must contain only alphabetic characters")
	}

	offset := page * limit

	switch sortBy {
	case "time":
		sortBySql = "time"
	case "rating":
		sortBySql = "rating_score"
	case "target":
		sortBySql = "target_to"
	default:
		sortBySql = "time"
	}

	rows := db.DB.Raw(fmt.Sprintf(`
		WITH total_count AS (
			SELECT COUNT(*) AS count FROM stocks WHERE ($1 = '' OR ticker = $1)
		)
		SELECT *, (SELECT count FROM total_count) AS total_count
		FROM stocks
		WHERE ($1 = '' OR ticker = $1)
		ORDER BY %s DESC
		LIMIT $2 OFFSET $3;
	`, sortBySql), filter, limit, offset).Scan(&stocks)

	if rows.Error != nil {
		return nil, rows.Error
	}

	if len(stocks) == 0 {
		stocks = []models.Stock{}
	}

	return stocks, nil
}
