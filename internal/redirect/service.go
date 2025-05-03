package redirect

import (
	"errors"
	"time"

	redisDB "github.com/mellgit/shorturl/internal/db"

	goredis "github.com/redis/go-redis/v9"
)

type Service interface {
	ResolveAndTrack(alias, ip, userAgent string) (string, error)
}
type RedirectService struct {
	repo Repository
	rdb  *goredis.Client
}

func NewService(repo Repository, rdb *goredis.Client) Service {
	return &RedirectService{repo: repo, rdb: rdb}
}

func (s *RedirectService) ResolveAndTrack(alias, ip, userAgent string) (string, error) {
	// 1. Check Redis
	cached, err := redisDB.Get(s.rdb, "short:"+alias)
	if err == nil {
		go s.repo.SaveClick(&Click{Alias: alias, IP: ip, UserAgent: userAgent})
		return cached, nil
	}

	// 2. Fallback to Postgres
	original, expiresAt, err := s.repo.FindOriginalByAlias(alias)
	if err != nil {
		return "", errors.New("link not found")
	}

	if time.Now().After(expiresAt) {
		return "", errors.New("link expired")
	}

	// 3. Save to Redis
	ttl := time.Until(expiresAt)
	_ = redisDB.SetWithTTL(s.rdb, "short:"+alias, original, ttl)

	// 4. Track click
	go s.repo.SaveClick(&Click{Alias: alias, IP: ip, UserAgent: userAgent})

	return original, nil
}
