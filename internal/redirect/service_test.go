package redirect

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPostgresRepo struct {
	mock.Mock
	saveClickCalled chan struct{}
}

func (m *mockPostgresRepo) FindOriginalByAlias(alias string) (string, time.Time, error) {
	args := m.Called(alias)
	return args.String(0), args.Get(1).(time.Time), args.Error(2)
}

func (m *mockPostgresRepo) SaveClick(c *Click) error {
	args := m.Called(c)
	if m.saveClickCalled != nil {
		m.saveClickCalled <- struct{}{}
	}
	return args.Error(0)
}

func waitForSaveClick(t *testing.T, ch chan struct{}) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatal("SaveClick was not called in time")
	}
}

type mockRedisRepo struct {
	mock.Mock
}

func (m *mockRedisRepo) SetWithTTL(key, value string, ttl time.Duration) error {
	args := m.Called(key, value, ttl)
	return args.Error(0)
}

func (m *mockRedisRepo) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func TestRedirectService_ResolveAndTrack_CacheHit(t *testing.T) {
	// Arrange
	postgresRepo := &mockPostgresRepo{saveClickCalled: make(chan struct{}, 1)}
	redisRepo := new(mockRedisRepo)
	service := NewService(postgresRepo, redisRepo)

	redisRepo.On("Get", "short:abc123").Return("https://example.com", nil)
	postgresRepo.On("SaveClick", mock.AnythingOfType("*redirect.Click")).Return(nil)

	// Act
	original, err := service.ResolveAndTrack("abc123", "127.0.0.1", "Mozilla/5.0")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", original)
	waitForSaveClick(t, postgresRepo.saveClickCalled)
	redisRepo.AssertExpectations(t)
	postgresRepo.AssertExpectations(t)
}

func TestRedirectService_ResolveAndTrack_CacheMiss_Success(t *testing.T) {
	// Arrange
	postgresRepo := &mockPostgresRepo{saveClickCalled: make(chan struct{}, 1)}
	redisRepo := new(mockRedisRepo)
	service := NewService(postgresRepo, redisRepo)

	expiresAt := time.Now().Add(1 * time.Hour)

	redisRepo.On("Get", "short:abc123").Return("", errors.New("cache miss"))
	postgresRepo.On("FindOriginalByAlias", "abc123").Return("https://example.com", expiresAt, nil)
	redisRepo.On("SetWithTTL", "short:abc123", "https://example.com", mock.AnythingOfType("time.Duration")).Return(nil)
	postgresRepo.On("SaveClick", mock.AnythingOfType("*redirect.Click")).Return(nil)

	// Act
	original, err := service.ResolveAndTrack("abc123", "127.0.0.1", "Mozilla/5.0")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", original)
	waitForSaveClick(t, postgresRepo.saveClickCalled)
	redisRepo.AssertExpectations(t)
	postgresRepo.AssertExpectations(t)
}

func TestRedirectService_ResolveAndTrack_CacheMiss_LinkExpired(t *testing.T) {
	// Arrange
	postgresRepo := new(mockPostgresRepo)
	redisRepo := new(mockRedisRepo)
	service := NewService(postgresRepo, redisRepo)

	expiresAt := time.Now().Add(-1 * time.Hour)

	redisRepo.On("Get", "short:abc123").Return("", errors.New("cache miss"))
	postgresRepo.On("FindOriginalByAlias", "abc123").Return("https://example.com", expiresAt, nil)

	// Act
	original, err := service.ResolveAndTrack("abc123", "127.0.0.1", "Mozilla/5.0")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "link expired")
	assert.Empty(t, original)
	redisRepo.AssertExpectations(t)
	postgresRepo.AssertExpectations(t)
	redisRepo.AssertNotCalled(t, "SetWithTTL")
	postgresRepo.AssertNotCalled(t, "SaveClick")
}

func TestRedirectService_ResolveAndTrack_PostgresError(t *testing.T) {
	// Arrange
	postgresRepo := new(mockPostgresRepo)
	redisRepo := new(mockRedisRepo)
	service := NewService(postgresRepo, redisRepo)

	redisRepo.On("Get", "short:abc123").Return("", errors.New("cache miss"))
	postgresRepo.On("FindOriginalByAlias", "abc123").Return("", time.Time{}, errors.New("db error"))

	// Act
	original, err := service.ResolveAndTrack("abc123", "127.0.0.1", "Mozilla/5.0")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to find original url")
	assert.Empty(t, original)
	redisRepo.AssertExpectations(t)
	postgresRepo.AssertExpectations(t)
}

func TestRedirectService_ResolveAndTrack_CacheMiss_RedisSetError_Ignored(t *testing.T) {
	// Arrange
	postgresRepo := &mockPostgresRepo{saveClickCalled: make(chan struct{}, 1)}
	redisRepo := new(mockRedisRepo)
	service := NewService(postgresRepo, redisRepo)

	expiresAt := time.Now().Add(1 * time.Hour)

	redisRepo.On("Get", "short:abc123").Return("", errors.New("cache miss"))
	postgresRepo.On("FindOriginalByAlias", "abc123").Return("https://example.com", expiresAt, nil)
	redisRepo.On("SetWithTTL", "short:abc123", "https://example.com", mock.AnythingOfType("time.Duration")).Return(errors.New("redis set error"))
	postgresRepo.On("SaveClick", mock.AnythingOfType("*redirect.Click")).Return(nil)

	// Act
	original, err := service.ResolveAndTrack("abc123", "127.0.0.1", "Mozilla/5.0")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", original)
	waitForSaveClick(t, postgresRepo.saveClickCalled)
	redisRepo.AssertExpectations(t)
	postgresRepo.AssertExpectations(t)
}
