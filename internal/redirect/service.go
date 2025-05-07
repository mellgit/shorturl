package redirect

import (
	"fmt"
	"time"
)

type Service interface {
	ResolveAndTrack(alias, ip, userAgent string) (string, error)
}
type RedirectService struct {
	postgresRepo PostgresRepository
	redisRepo    RedisRepository
}

func NewService(postgresRepo PostgresRepository, redisRepo RedisRepository) Service {
	return &RedirectService{postgresRepo: postgresRepo, redisRepo: redisRepo}
}

func (s *RedirectService) ResolveAndTrack(alias, ip, userAgent string) (string, error) {
	// 1. Check Redis
	cached, err := s.redisRepo.Get("short:" + alias)
	if err == nil {
		go s.postgresRepo.SaveClick(&Click{Alias: alias, IP: ip, UserAgent: userAgent})
		return cached, nil
	}

	// 2. Fallback to Postgres
	original, expiresAt, err := s.postgresRepo.FindOriginalByAlias(alias)
	if err != nil {
		return "", fmt.Errorf("failed to find original url for alias %s: %w", alias, err)
	}

	if time.Now().After(expiresAt) {
		return "", fmt.Errorf("link expired %s: %w", alias, err)
	}

	// 3. Save to Redis
	ttl := time.Until(expiresAt)
	_ = s.redisRepo.SetWithTTL("short:"+alias, original, ttl)

	// 4. Track click
	go s.postgresRepo.SaveClick(&Click{Alias: alias, IP: ip, UserAgent: userAgent})

	return original, nil
}
