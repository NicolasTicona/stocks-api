package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/nicolasticona/stocks-api/db"
)

func RedisSave(param string, obj []byte, timeToCache time.Duration) error {
	err := db.RedisClient.Set(context.Background(), param, obj, timeToCache).Err()

	if err != nil {
		fmt.Println("Error setting stock in Redis:", err)
	}

	fmt.Println("Setting in Redis:", param)

	return err
}

func RedisGet(param string) (string, error) {
	value, err := db.RedisClient.Get(context.Background(), param).Result()

	if err != nil {
		fmt.Println("Error getting param from redis:", err)
	}

	return value, err
}
