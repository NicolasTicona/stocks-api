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

func RedisDelete(key string) error {
	searchPattern := key + "*"

	var foundedRecordCount int = 0
	iter := db.RedisClient.Scan(context.Background(), 0, searchPattern, 0).Iterator()
	fmt.Printf("YOUR SEARCH PATTERN= %s\n", searchPattern)
	for iter.Next(context.Background()) {
		fmt.Printf("Deleted= %s\n", iter.Val())
		db.RedisClient.Del(context.Background(), iter.Val())
		foundedRecordCount++
	}

	if err := iter.Err(); err != nil {
		return err
	}

	fmt.Printf("Deleted Count %d\n", foundedRecordCount)

	return nil
}
