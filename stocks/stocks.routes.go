package stocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type StocksHandler struct {
	Controller *StocksController
}

func (handler *StocksHandler) GetStocksHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	filter := r.URL.Query().Get("filter")
	sortBy := r.URL.Query().Get("sortBy")

	page, err := strconv.ParseUint(pageStr, 10, 16)
	if err != nil {
		page = 0
	}

	limit, err := strconv.ParseUint(limitStr, 10, 16)
	if err != nil || limit < 1 {
		limit = 10
	}

	if limit > 50 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Limit cannot be greater than 50"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response, err := handler.Controller.GetStocks(w, uint16(page), uint16(limit), sortBy, filter)

	if err != nil {
		fmt.Printf("Error fetching recommendations: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch stocks"})
		return
	}

	json.NewEncoder(w).Encode(response)
}
