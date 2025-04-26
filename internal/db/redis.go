package db

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"), // "localhost:6379"
		Password: "",                      // or use os.Getenv("REDIS_PASS")
		DB:       0,
	})
}

func SetWithTTL(rdb *redis.Client, key, value string, ttl time.Duration) error {
	return rdb.Set(Ctx, key, value, ttl).Err()
}

func Get(rdb *redis.Client, key string) (string, error) {
	return rdb.Get(Ctx, key).Result()
}
