package stockssync

import (
	"encoding/json"
	"net/http"

	"github.com/nicolasticona/stocks-api/utils"
)

type StocksSyncHandler struct {
	Controller *StocksSyncController
}

func (handler *StocksSyncHandler) SyncDbHandler(w http.ResponseWriter, r *http.Request) {
	err := handler.Controller.InsertStockRatings()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		println("Failed to insert stock ratings:", err)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to insert stock ratings"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Stock ratings inserted successfully"})
}

func (handler *StocksSyncHandler) RemoveRedisKeyPattern(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Key parameter is required"})
		return
	}
	err := utils.RedisDelete(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		println("Failed to delete key from Redis:", err)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete key from Redis"})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Key deleted successfully"})
}
