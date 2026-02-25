package service

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"exam-system/internal/models"
	"exam-system/internal/repository"
)

type AuthService interface {
	Signup(ctx context.Context, email, password, username, role string) (*models.User, error)
	Login(ctx context.Context, email, password string) (string, *models.User, error)
}

type authService struct {
	repo      repository.UserRepository
	jwtSecret []byte
}

func NewAuthService(repo repository.UserRepository) AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret"
	}
	return &authService{
		repo:      repo,
		jwtSecret: []byte(secret),
	}
}

func (s *authService) Signup(ctx context.Context, email, password, username, role string) (*models.User, error) {
	if role == "" {
		role = "student"
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        email,
		Username:     username,
		PasswordHash: string(hash),
		Role:         role,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, *models.User, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", nil, errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"sub":      user.ID,
		"email":    user.Email,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", nil, err
	}

	return signed, user, nil
}

