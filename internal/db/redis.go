package db

import (
	"fmt"
	"github.com/mellgit/shorturl/internal/config"
	"github.com/redis/go-redis/v9"
)

func RedisClient(envCfg config.EnvConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", envCfg.RedisHost, envCfg.RedisPort),
		Password: "",
		DB:       envCfg.RedisDB,
	})
}
