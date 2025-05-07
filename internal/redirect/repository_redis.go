package redirect

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisRepository interface {
	SetWithTTL(key, value string, ttl time.Duration) error
	Get(key string) (string, error)
}

type redisRepository struct {
	ctx    context.Context
	client *redis.Client
}

func NewRedisRepo(client *redis.Client) RedisRepository {
	return &redisRepository{ctx: context.Background(), client: client}
}

func (r *redisRepository) SetWithTTL(key, value string, ttl time.Duration) error {
	return r.client.Set(r.ctx, key, value, ttl).Err()
}
func (r *redisRepository) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}
