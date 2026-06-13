package shortener

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) Save(url *URL) error {
	args := m.Called(url)
	return args.Error(0)
}

func (m *mockRepository) IsAliasTaken(alias string) (bool, error) {
	args := m.Called(alias)
	return args.Bool(0), args.Error(1)
}

func (m *mockRepository) Stats(alias string) (int, error) {
	args := m.Called(alias)
	return args.Int(0), args.Error(1)
}

func (m *mockRepository) List() (*[]URL, error) {
	args := m.Called()
	return args.Get(0).(*[]URL), args.Error(1)
}

func (m *mockRepository) Delete(alias string) error {
	args := m.Called(alias)
	return args.Error(0)
}

func (m *mockRepository) UpdateAlias(alias, newAlias string) error {
	args := m.Called(alias, newAlias)
	return args.Error(0)
}

func (m *mockRepository) GetUrlFromAlias(alias string) (string, error) {
	args := m.Called(alias)
	return args.String(0), args.Error(1)
}

func TestShortenerService_CreateShortURL_AutoAlias(t *testing.T) {
	// Arrange
	repo := new(mockRepository)
	service := NewService(repo)
	userID := uuid.New()

	repo.On("IsAliasTaken", mock.AnythingOfType("string")).Return(false, nil).Once()
	repo.On("Save", mock.AnythingOfType("*shortener.URL")).Return(nil).Once()

	// Act
	url, err := service.CreateShortURL(userID, "https://example.com", "", 24)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, userID, url.UserID)
	assert.Equal(t, "https://example.com", url.Original)
	assert.Len(t, url.Alias, aliasLength)
	assert.WithinDuration(t, time.Now().Add(24*time.Hour), url.ExpiresAt, time.Second)
	repo.AssertExpectations(t)
}

func TestShortenerService_CreateShortURL_CustomAlias_Available(t *testing.T) {
	// Arrange
	repo := new(mockRepository)
	service := NewService(repo)
	userID := uuid.New()

	repo.On("IsAliasTaken", "custom").Return(false, nil).Once()
	repo.On("Save", mock.AnythingOfType("*shortener.URL")).Return(nil).Once()

	// Act
	url, err := service.CreateShortURL(userID, "https://example.com", "custom", 12)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "custom", url.Alias)
	assert.WithinDuration(t, time.Now().Add(12*time.Hour), url.ExpiresAt, time.Second)
	repo.AssertExpectations(t)
}

func TestShortenerService_CreateShortURL_CustomAlias_Taken(t *testing.T) {
	// Arrange
	repo := new(mockRepository)
	service := NewService(repo)
	userID := uuid.New()

	repo.On("IsAliasTaken", "custom").Return(true, nil).Once()

	// Act
	url, err := service.CreateShortURL(userID, "https://example.com", "custom", 24)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "custom alias is already used")
	assert.Nil(t, url)
	repo.AssertExpectations(t)
}

func TestShortenerService_CreateShortURL_IsAliasTaken_Error(t *testing.T) {
	// Arrange
	repo := new(mockRepository)
	service := NewService(repo)
	userID := uuid.New()

	repo.On("IsAliasTaken", "custom").Return(false, errors.New("db error")).Once()

	// Act
	url, err := service.CreateShortURL(userID, "https://example.com", "custom", 24)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, url)
	repo.AssertExpectations(t)
}

func TestShortenerService_CreateShortURL_Save_Error(t *testing.T) {
	// Arrange
	repo := new(mockRepository)
	service := NewService(repo)
	userID := uuid.New()

	repo.On("IsAliasTaken", mock.AnythingOfType("string")).Return(false, nil).Once()
	repo.On("Save", mock.AnythingOfType("*shortener.URL")).Return(errors.New("save error")).Once()

	// Act
	url, err := service.CreateShortURL(userID, "https://example.com", "", 24)

	// Assert
	assert.Error(t, err)
	assert.NotNil(t, url)
	repo.AssertExpectations(t)
}

func TestShortenerService_CreateShortURL_CollisionRetry(t *testing.T) {
	// Arrange
	repo := new(mockRepository)
	service := NewService(repo)
	userID := uuid.New()

	// seed random for determinism in this test
	rand.Seed(1)

	repo.On("IsAliasTaken", mock.AnythingOfType("string")).Return(true, nil).Once()
	repo.On("IsAliasTaken", mock.AnythingOfType("string")).Return(false, nil).Once()
	repo.On("Save", mock.AnythingOfType("*shortener.URL")).Return(nil).Once()

	// Act
	url, err := service.CreateShortURL(userID, "https://example.com", "", 24)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, url)
	assert.Len(t, url.Alias, aliasLength)
	repo.AssertExpectations(t)
}

func TestShortenerService_Stats(t *testing.T) {
	// Arrange
	repo := new(mockRepository)
	service := NewService(repo)

	repo.On("Stats", "abc123").Return(42, nil)

	// Act
	count, err := service.Stats("abc123")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 42, count)
	repo.AssertExpectations(t)
}

func TestShortenerService_Stats_Error(t *testing.T) {
	// Arrange
	repo := new(mockRepository)
	service := NewService(repo)

	repo.On("Stats", "abc123").Return(0, errors.New("db error"))

	// Act
	count, err := service.Stats("abc123")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error getting stats")
	assert.Equal(t, 0, count)
	repo.AssertExpectations(t)
}

func TestShortenerService_List(t *testing.T) {
	// Arrange
	repo := new(mockRepository)
	service := NewService(repo)

	urls := &[]URL{{Alias: "abc123", Original: "https://example.com"}}
	repo.On("List").Return(urls, nil)

	// Act
	result, err := service.List()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, urls, result)
	repo.AssertExpectations(t)
}

func TestShortenerService_List_Error(t *testing.T) {
	// Arrange
	repo := new(mockRepository)
	service := NewService(repo)

	repo.On("List").Return(&[]URL{}, errors.New("db error"))

	// Act
	result, err := service.List()

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestShortenerService_Delete(t *testing.T) {
	// Arrange
	repo := new(mockRepository)
	service := NewService(repo)

	repo.On("Delete", "abc123").Return(nil)

	// Act
	err := service.Delete("abc123")

	// Assert
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestShortenerService_UpdateAlias(t *testing.T) {
	// Arrange
	repo := new(mockRepository)
	service := NewService(repo)

	repo.On("UpdateAlias", "abc123", "new123").Return(nil)

	// Act
	err := service.UpdateAlias("abc123", "new123")

	// Assert
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestShortenerService_GenerateQRCode(t *testing.T) {
	// Arrange
	repo := new(mockRepository)
	service := NewService(repo)

	repo.On("GetUrlFromAlias", "abc123").Return("https://example.com", nil)

	// Act
	qr, err := service.GenerateQRCode("abc123")

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, qr)
	repo.AssertExpectations(t)
}

func TestShortenerService_GenerateQRCode_RepositoryError(t *testing.T) {
	// Arrange
	repo := new(mockRepository)
	service := NewService(repo)

	repo.On("GetUrlFromAlias", "abc123").Return("", errors.New("not found"))

	// Act
	qr, err := service.GenerateQRCode("abc123")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error getting url from alias")
	assert.Nil(t, qr)
	repo.AssertExpectations(t)
}
