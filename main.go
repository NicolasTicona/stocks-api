package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/nicolasticona/stocks-api/db"
	"github.com/nicolasticona/stocks-api/factory"
	"github.com/nicolasticona/stocks-api/models"
	"github.com/nicolasticona/stocks-api/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file")
	}

	db.RedisConnection()
	db.DbConnection()

	db.DB.AutoMigrate(models.Stock{})
	db.DB.AutoMigrate(models.StockRecommendation{})

	// Initialize components using the factory function
	components := factory.InitializeComponents()

	// Set up the router
	router := mux.NewRouter()

	router.HandleFunc("/sync", components.StocksSyncHandler.SyncDbHandler).Methods("POST")
	router.HandleFunc("/stocks", components.StocksHandler.GetStocksHandler).Methods("GET")
	router.HandleFunc("/recommendations", components.RecommendationsHandler.GetRecommendationsHandler).Methods("GET")
	router.HandleFunc("/analyze", routes.GetStockAnalysisHandler).Methods("GET")
	router.HandleFunc("/test-analyze", routes.GetRedisHandler).Methods("GET")
	router.HandleFunc("/remove-redis-key", components.StocksSyncHandler.RemoveRedisKeyPattern).Methods("DELETE")

	var corsOptions handlers.CORSOption
	isDev := os.Getenv("IS_DEV")

	if isDev == "true" {
		corsOptions = handlers.AllowedOrigins([]string{"http://localhost:5173"})
	} else {
		corsOptions = handlers.AllowedOrigins([]string{"https://stocks-ui.vercel.app"})
	}

	corsHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	corsMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})

	http.ListenAndServe(":8000", handlers.CORS(corsOptions, corsHeaders, corsMethods)(router))
}
