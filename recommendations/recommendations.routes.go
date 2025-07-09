package recommendations

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RecommendationsHandler struct {
	Controller *RecommendationsController
}

func (handler *RecommendationsHandler) GetRecommendationsHandler(w http.ResponseWriter, r *http.Request) {
	stocks, err := handler.Controller.GetRecommendations(w)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		fmt.Printf("Error fetching recommendations: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch recommendations"})
		return
	}

	response := map[string]interface{}{
		"stocks": stocks,
	}

	json.NewEncoder(w).Encode(response)
}
