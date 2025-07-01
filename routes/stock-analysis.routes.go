package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func GetStockAnalysisHandler(w http.ResponseWriter, r *http.Request) {
	stock := r.URL.Query().Get("stock")
	if stock == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Stock parameter is required",
		})
		return
	}

	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		fmt.Println("Error converting REDIS_DB to integer:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Error connecting to Redis",
		})
		return
	}

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})

	value, err := client.Get(context.Background(), stock).Result()
	if err == nil {
		var cacheResult map[string]interface{}
		err = json.Unmarshal([]byte(value), &cacheResult)
		if err == nil {
			fmt.Println("Cache hit for stock:", stock)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(cacheResult)
			return
		}
	}

	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to resolve script path",
		})
		return
	}

	scriptPath := filepath.Join(workingDir, "stock-analysis", "analysis.py")

	fmt.Println("Executing Python script at path:", scriptPath)

	out, err := exec.Command("python3", scriptPath, stock).Output()
	if err != nil {
		fmt.Println("Error executing Python script:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "There was an error analysing the stock",
		})
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal(out, &result)
	if err != nil {
		fmt.Println("Error parsing Python response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "There was an error returning the stock analysis",
		})
		return
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Error serializing result to JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to cache the stock analysis",
		})
		return
	}

	err = client.Set(context.Background(), stock, resultJSON, time.Hour).Err()
	if err != nil {
		fmt.Println("Error setting stock in Redis:", err)
	}

	fmt.Println("Setting stock in Redis:", stock, "with value:", string(resultJSON))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func GetRedisHandler(w http.ResponseWriter, r *http.Request) {
	stock := r.URL.Query().Get("stock")
	if stock == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Stock parameter is required",
		})
		return
	}

	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		fmt.Println("Error converting REDIS_DB to integer:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Error connecting to Redis",
		})
		return
	}

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})

	value, err := client.Get(context.Background(), stock).Result()
	if err != nil {
		fmt.Println("Error getting stock from Redis:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]bool{
			"value": false,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"value": value})
}
