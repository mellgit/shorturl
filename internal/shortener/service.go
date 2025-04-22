package shortener

import (
	"errors"
	"math/rand"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const aliasLength = 6

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo}
}

func (s *Service) CreateShortURL(userID int64, original, customAlias string, ttlHours int) (*URL, error) {
	var alias string

	// custom alias
	if customAlias != "" {
		exists, err := s.repo.IsAliasTaken(customAlias)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("alias already taken")
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

func generateRandomString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}
