package shortener

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const aliasLength = 6

type Service interface {
	CreateShortURL(userID int64, original, customAlias string, ttlHours int) (*URL, error)
	Stats(alias string) (int, error)
	List() (*[]URL, error)
	Delete(alias string) error
	UpdateAlias(alias, newAlias string) error
}
type ShortenerService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &ShortenerService{repo}
}

func (s *ShortenerService) CreateShortURL(userID int64, original, customAlias string, ttlHours int) (*URL, error) {
	var alias string

	// custom alias
	if customAlias != "" {
		exists, err := s.repo.IsAliasTaken(customAlias)
		if err != nil {
			return nil, fmt.Errorf("error checking if alias exists: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("custom alias is already used")
		}
		alias = customAlias
	} else {
		// generate uniq alias
		for {
			alias = generateRandomString(aliasLength)
			exists, err := s.repo.IsAliasTaken(alias)
			if err != nil {
				return nil, err
			}
			if !exists {
				break
			}
		}
	}

	url := &URL{
		UserID:    userID,
		Original:  original,
		Alias:     alias,
		ExpiresAt: time.Now().Add(time.Duration(ttlHours) * time.Hour),
	}

	err := s.repo.Save(url)
	return url, err
}

func (s *ShortenerService) Stats(alias string) (int, error) {
	count, err := s.repo.Stats(alias)
	if err != nil {
		return 0, fmt.Errorf("error getting stats: %w", err)
	}
	return count, nil
}

func (s *ShortenerService) List() (*[]URL, error) {
	urls, err := s.repo.List()
	if err != nil {
		return nil, fmt.Errorf("error listing urls: %w", err)
	}
	return urls, nil
}

func (s *ShortenerService) Delete(alias string) error {
	return s.repo.Delete(alias)
}

func (s *ShortenerService) UpdateAlias(alias, newAlias string) error {
	return s.repo.UpdateAlias(alias, newAlias)
}

func generateRandomString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}
