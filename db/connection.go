package db

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DbConnection() {
	var error error
	var DSN = os.Getenv("POSTGRES_DSN")

	DB, error = gorm.Open(postgres.Open(DSN), &gorm.Config{})

	if error != nil {
		log.Fatal("Error connecting to the database:", error)
		return
	}

	log.Println("Database connection established successfully")
}
