package db

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var RedisClient *redis.Client

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

func RedisConnection() error {
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		fmt.Println("Error converting REDIS_DB to integer:", err)
		return err
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})

	fmt.Println("Redis connection established successfully")

	return nil
}
