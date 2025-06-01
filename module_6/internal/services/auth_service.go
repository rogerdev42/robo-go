package services

import (
	"context"
	"errors"
	"fmt"
	"module_6/internal/config"
	"module_6/internal/database/repository"
	"module_6/internal/logger"
	"module_6/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  repository.UserRepository
	jwtSecret string
	jwtExpire time.Duration
	logger    logger.Logger
}

func NewAuthService(
	userRepo repository.UserRepository,
	cfg *config.Config,
	log logger.Logger,
) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: cfg.JWTSecret,
		jwtExpire: time.Duration(cfg.JWTExpireHours) * time.Hour,
		logger:    log,
	}
}

func (s *AuthService) SignUp(ctx context.Context, req *models.UserCreate) (*models.User, string, error) {
	s.logger.Info("Attempting to sign up user",
		logger.String("email", req.Email),
		logger.String("name", req.Name))

	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		s.logger.Error("Failed to check existing user by email", logger.Error(err))
		return nil, "", err
	}
	if existingUser != nil {
		s.logger.Warn("User with email already exists", logger.String("email", req.Email))
		return nil, "", models.ErrAlreadyExists
	}

	// Проверяем, существует ли пользователь с таким именем
	existingUser, err = s.userRepo.GetByName(ctx, req.Name)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		s.logger.Error("Failed to check existing user by name", logger.Error(err))
		return nil, "", err
	}
	if existingUser != nil {
		s.logger.Warn("User with name already exists", logger.String("name", req.Name))
		return nil, "", models.ErrAlreadyExists
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash password", logger.Error(err))
		return nil, "", err
	}

	// Создаем пользователя
	user := &models.User{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: string(hashedPassword),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error("Failed to create user", logger.Error(err))
		return nil, "", err
	}

	// Генерируем JWT токен
	token, err := s.GenerateToken(user.ID)
	if err != nil {
		s.logger.Error("Failed to generate token", logger.Error(err))
		return nil, "", err
	}

	s.logger.Info("User signed up successfully",
		logger.Int("user_id", user.ID),
		logger.String("email", user.Email))

	return user, token, nil
}

// SignIn вход пользователя
func (s *AuthService) SignIn(ctx context.Context, req *models.UserLogin) (*models.User, string, error) {
	s.logger.Info("Attempting to sign in user", logger.String("email", req.Email))

	// Получаем пользователя по email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			s.logger.Warn("User not found", logger.String("email", req.Email))
			return nil, "", models.ErrInvalidCredentials
		}
		s.logger.Error("Failed to get user", logger.Error(err))
		return nil, "", err
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		s.logger.Warn("Invalid password", logger.String("email", req.Email))
		return nil, "", models.ErrInvalidCredentials
	}

	// Генерируем JWT токен
	token, err := s.GenerateToken(user.ID)
	if err != nil {
		s.logger.Error("Failed to generate token", logger.Error(err))
		return nil, "", err
	}

	s.logger.Info("User signed in successfully",
		logger.Int("user_id", user.ID),
		logger.String("email", user.Email))

	return user, token, nil
}

// ValidateToken проверяет JWT токен и возвращает user_id
func (s *AuthService) ValidateToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(float64)
		if !ok {
			return 0, errors.New("invalid token claims")
		}
		return int(userID), nil
	}

	return 0, models.ErrUnauthorized
}

// GenerateToken генерирует JWT токен для пользователя (ПУБЛИЧНЫЙ МЕТОД)
func (s *AuthService) GenerateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(s.jwtExpire).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}