package factory

import (
	"github.com/nicolasticona/stocks-api/recommendations"
	"github.com/nicolasticona/stocks-api/stocks"
	stockssync "github.com/nicolasticona/stocks-api/stocks-sync"
)

type AppComponents struct {
	RecommendationsHandler *recommendations.RecommendationsHandler
	StocksHandler          *stocks.StocksHandler
	StocksSyncHandler      *stockssync.StocksSyncHandler
}

func InitializeComponents() *AppComponents {
	// Initialize repositories
	recommendationsRepository := &recommendations.RecommendationRepository{}
	stocksRepository := &stocks.StocksRepository{}
	stocksSyncRepository := &stockssync.StocksSyncRepository{}

	// Initialize controllers
	recommendationsController := &recommendations.RecommendationsController{
		Repository: recommendationsRepository,
	}
	stocksController := &stocks.StocksController{
		Repository: stocksRepository,
	}
	stocksSyncController := &stockssync.StocksSyncController{
		Repository: stocksSyncRepository,
	}

	// Initialize handlers
	recommendationsHandler := &recommendations.RecommendationsHandler{
		Controller: recommendationsController,
	}
	stocksHandler := &stocks.StocksHandler{
		Controller: stocksController,
	}
	stocksSyncHandler := &stockssync.StocksSyncHandler{
		Controller: stocksSyncController,
	}

	return &AppComponents{
		RecommendationsHandler: recommendationsHandler,
		StocksHandler:          stocksHandler,
		StocksSyncHandler:      stocksSyncHandler,
	}
}
