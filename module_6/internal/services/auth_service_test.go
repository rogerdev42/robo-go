package services_test

import (
	"context"
	"errors"
	"module_6/internal/config"
	"module_6/internal/logger"
	"module_6/internal/models"
	"module_6/internal/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository мок репозитория пользователей
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	if args.Error(0) == nil {
		// Симулируем успешное создание
		user.ID = 1
		user.CreatedAt = time.Now()
	}
	return args.Error(0)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByName(ctx context.Context, name string) (*models.User, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupAuthService(t *testing.T) (*services.AuthService, *MockUserRepository) {
	mockRepo := &MockUserRepository{}
	cfg := &config.Config{
		JWTSecret:      "test-secret-key",
		JWTExpireHours: 24,
	}
	
	log, err := logger.NewSlogLogger("debug", "text", "stdout")
	require.NoError(t, err)
	
	service := services.NewAuthService(mockRepo, cfg, log)
	return service, mockRepo
}

func TestAuthService_SignUp_Success(t *testing.T) {
	service, mockRepo := setupAuthService(t)

	req := &models.UserCreate{
		Email:    "test@example.com",
		Name:     "testuser",
		Password: "password123",
	}

	// Настраиваем моки
	mockRepo.On("GetByEmail", mock.Anything, "test@example.com").
		Return(nil, models.ErrNotFound)
	mockRepo.On("GetByName", mock.Anything, "testuser").
		Return(nil, models.ErrNotFound)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).
		Return(nil)

	// Выполняем тест
	user, token, err := service.SignUp(context.Background(), req)

	// Проверки
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, token)
	assert.Equal(t, req.Email, user.Email)
	assert.Equal(t, req.Name, user.Name)
	assert.NotEmpty(t, user.PasswordHash)

	// Проверяем что пароль захеширован
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_SignUp_EmailExists(t *testing.T) {
	service, mockRepo := setupAuthService(t)

	req := &models.UserCreate{
		Email:    "existing@example.com",
		Name:     "testuser",
		Password: "password123",
	}

	existingUser := &models.User{
		ID:    1,
		Email: "existing@example.com",
		Name:  "existinguser",
	}

	mockRepo.On("GetByEmail", mock.Anything, "existing@example.com").
		Return(existingUser, nil)

	user, token, err := service.SignUp(context.Background(), req)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, models.ErrAlreadyExists))
	assert.Nil(t, user)
	assert.Empty(t, token)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_SignUp_NameExists(t *testing.T) {
	service, mockRepo := setupAuthService(t)

	req := &models.UserCreate{
		Email:    "test@example.com",
		Name:     "existinguser",
		Password: "password123",
	}

	existingUser := &models.User{
		ID:   1,
		Name: "existinguser",
	}

	mockRepo.On("GetByEmail", mock.Anything, "test@example.com").
		Return(nil, models.ErrNotFound)
	mockRepo.On("GetByName", mock.Anything, "existinguser").
		Return(existingUser, nil)

	user, token, err := service.SignUp(context.Background(), req)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, models.ErrAlreadyExists))
	assert.Nil(t, user)
	assert.Empty(t, token)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_SignIn_Success(t *testing.T) {
	service, mockRepo := setupAuthService(t)

	password := "password123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	testUser := &models.User{
		ID:           1,
		Email:        "test@example.com",
		Name:         "testuser",
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
	}

	req := &models.UserLogin{
		Email:    "test@example.com",
		Password: password,
	}

	mockRepo.On("GetByEmail", mock.Anything, "test@example.com").
		Return(testUser, nil)

	user, token, err := service.SignIn(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, token)
	assert.Equal(t, testUser.Email, user.Email)
	assert.Equal(t, testUser.Name, user.Name)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_SignIn_UserNotFound(t *testing.T) {
	service, mockRepo := setupAuthService(t)

	req := &models.UserLogin{
		Email:    "notfound@example.com",
		Password: "password123",
	}

	mockRepo.On("GetByEmail", mock.Anything, "notfound@example.com").
		Return(nil, models.ErrNotFound)

	user, token, err := service.SignIn(context.Background(), req)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, models.ErrInvalidCredentials))
	assert.Nil(t, user)
	assert.Empty(t, token)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_SignIn_WrongPassword(t *testing.T) {
	service, mockRepo := setupAuthService(t)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	require.NoError(t, err)

	testUser := &models.User{
		ID:           1,
		Email:        "test@example.com",
		Name:         "testuser",
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
	}

	req := &models.UserLogin{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	mockRepo.On("GetByEmail", mock.Anything, "test@example.com").
		Return(testUser, nil)

	user, token, err := service.SignIn(context.Background(), req)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, models.ErrInvalidCredentials))
	assert.Nil(t, user)
	assert.Empty(t, token)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_ValidateToken_Success(t *testing.T) {
	service, _ := setupAuthService(t)

	// Создаем валидный токен
	token, err := service.GenerateToken(123)
	require.NoError(t, err)

	userID, err := service.ValidateToken(token)

	assert.NoError(t, err)
	assert.Equal(t, 123, userID)
}

func TestAuthService_ValidateToken_Invalid(t *testing.T) {
	service, _ := setupAuthService(t)

	userID, err := service.ValidateToken("invalid.token.here")

	assert.Error(t, err)
	assert.Zero(t, userID)
}

func TestAuthService_ValidateToken_Empty(t *testing.T) {
	service, _ := setupAuthService(t)

	userID, err := service.ValidateToken("")

	assert.Error(t, err)
	assert.Zero(t, userID)
}

func TestAuthService_GenerateToken(t *testing.T) {
	service, _ := setupAuthService(t)

	token, err := service.GenerateToken(456)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Проверяем что токен можно валидировать
	userID, err := service.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, 456, userID)
}