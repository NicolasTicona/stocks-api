package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/nicolasticona/stocks-api/db"
	"github.com/nicolasticona/stocks-api/models"
	"github.com/nicolasticona/stocks-api/routes"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file")
	}

	db.DbConnection()

	db.DB.AutoMigrate(models.Stock{})
	db.DB.AutoMigrate(models.StockRecommendation{})

	router := mux.NewRouter()

	router.HandleFunc("/sync", routes.SyncDbHandler).Methods("POST")
	router.HandleFunc("/stocks", routes.GetStocksHandler).Methods("GET")
	router.HandleFunc("/recommendations", routes.GetRecommendationsHandler).Methods("GET")
	router.HandleFunc("/analyze", routes.GetStockAnalysisHandler).Methods("GET")
	router.HandleFunc("/test-analyze", routes.GetRedisHandler).Methods("GET")

	// Add CORS middleware
	corsOptions := handlers.AllowedOrigins([]string{"stocks-ui.vercel.app"})
	corsHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	corsMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})

	http.ListenAndServe(":8000", handlers.CORS(corsOptions, corsHeaders, corsMethods)(router))
}
