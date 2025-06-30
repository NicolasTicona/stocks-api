package apis

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func FetchAllStockRatings() ([]map[string]interface{}, error) {
	client := &http.Client{}
	baseURL := "https://8j5baasof2.execute-api.us-west-2.amazonaws.com/production/swechallenge/list"
	allItems := []map[string]interface{}{}
	var nextPage string
	pageCount := 0

	ratings := map[string]float64{
		"Strong-Buy": 5, "Top Pick": 5, "Speculative Buy": 4, "Buy": 4, "Moderate Buy": 4,
		"Overweight": 3, "Positive": 3, "Outperform": 3, "Market Outperform": 3, "Outperformer": 3,
		"Hold": 2, "Neutral": 2, "Market Perform": 2, "Equal Weight": 2, "In-Line": 2,
		"Peer Perform": 2, "Sector Perform": 2, "Sector Weight": 2, "Unchanged": 1,
		"Reduce": 1, "Underweight": 1, "Underperform": 1, "Sector Underperform": 1,
		"Sell": 0, "Strong Sell": 0,
	}

	for {
		url := baseURL
		if nextPage != "" {
			url = fmt.Sprintf("%s?next_page=%s", baseURL, nextPage)
		}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		apiKey := os.Getenv("SWECHALLENGE_API_KEY")

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+apiKey)

		res, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch data from external API: %w", err)
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		var responseData struct {
			Items    []map[string]interface{} `json:"items"`
			NextPage string                   `json:"next_page"`
		}
		if err := json.Unmarshal(body, &responseData); err != nil {
			return nil, fmt.Errorf("failed to parse JSON response: %w", err)
		}

		// Process items
		for _, item := range responseData.Items {
			// Remove $ symbol and convert target_from and target_to to floats
			if targetFrom, ok := item["target_from"].(string); ok {
				item["target_from"], _ = strconv.ParseFloat(strings.ReplaceAll(targetFrom, "$", ""), 64)
			}
			if targetTo, ok := item["target_to"].(string); ok {
				item["target_to"], _ = strconv.ParseFloat(strings.ReplaceAll(targetTo, "$", ""), 64)
			}

			// Add target_score property
			if targetFrom, ok := item["target_from"].(float64); ok {
				if targetTo, ok := item["target_to"].(float64); ok {
					if targetTo > targetFrom {
						item["target_score"] = float64(1)
					} else {
						item["target_score"] = float64(0)
					}
				}
			}

			// Add rating_score property
			ratingFrom, okFrom := item["rating_from"].(string)
			ratingTo, okTo := item["rating_to"].(string)
			if okFrom && okTo {
				ratingFromScore := ratings[ratingFrom]
				ratingToScore := ratings[ratingTo]
				item["rating_score"] = ratingToScore - ratingFromScore
			}
		}

		allItems = append(allItems, responseData.Items...)

		if responseData.NextPage == "" {
			break
		}

		fmt.Printf("Fetched page %d with %d items\n", pageCount+1, len(responseData.Items))

		nextPage = responseData.NextPage
		pageCount++
	}

	return allItems, nil
}
