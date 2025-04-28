package db

import (
	"context"
	"fmt"
	"github.com/mellgit/shorturl/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func RedisClient(envCfg config.EnvConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", envCfg.RedisHost, envCfg.RedisPort),
		Password: "",
		DB:       envCfg.RedisDB,
	})
}

func SetWithTTL(rdb *redis.Client, key, value string, ttl time.Duration) error {
	return rdb.Set(Ctx, key, value, ttl).Err()
}

func Get(rdb *redis.Client, key string) (string, error) {
	return rdb.Get(Ctx, key).Result()
}
