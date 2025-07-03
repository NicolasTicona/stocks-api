package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

	// Check if the stock analysis is cached in Redis
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

	apiURL := "https://faas-nyc1-2ef2e6cc.doserverless.co/api/v1/namespaces/fn-4296c404-6588-442a-9180-d58afa975070/actions/ai/analysis?blocking=true&result=true"
	reqBody, _ := json.Marshal(map[string]string{
		"symbol": stock,
	})

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to create request to external analysis service",
		})
		return
	}

	fmt.Println("TOKEN: ", os.Getenv("CLOUD_FUNC_AI_ANALYSIS_TOKEN"))

	req.Header.Set("Authorization", "Basic "+os.Getenv("CLOUD_FUNC_AI_ANALYSIS_TOKEN"))
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Error executing HTTP request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to call external analysis service",
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Non-200 response from external service: %d\n", resp.StatusCode)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "External analysis service returned an error",
		})
		return
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error parsing response from external service:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to parse response from external analysis service",
		})
		return
	}

	// Cache the result in Redis for 1 hour
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
