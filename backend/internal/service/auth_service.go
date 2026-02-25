package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"backend/internal/auth"
	"backend/internal/config"
	"backend/internal/domain"
	"backend/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	cfg      *config.Config
	userRepo *repository.UserRepository
}

func NewAuthService(cfg *config.Config, userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		cfg:      cfg,
		userRepo: userRepo,
	}
}

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

type AuthResult struct {
	Token string       `json:"token"`
	User  *domain.User `json:"user"`
}

func (s *AuthService) Register(ctx context.Context, in RegisterInput) (*AuthResult, error) {
	email := strings.TrimSpace(strings.ToLower(in.Email))
	if email == "" || in.Password == "" {
		return nil, errors.New("email and password are required")
	}

	// Проверяем, что пользователя еще нет
	existing, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:        email,
		Username:     in.Name,
		PasswordHash: string(hash),
		Role:         "student",
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	token, err := auth.GenerateToken(s.cfg.JWTSecret, user.ID, user.Role, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*AuthResult, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	println("[AUTH_SERVICE] Looking up user with email:", email)
	
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		println("[AUTH_SERVICE] Database error:", err.Error())
		return nil, err
	}
	if user == nil {
		println("[AUTH_SERVICE] User not found")
		return nil, errors.New("invalid credentials")
	}

	println("[AUTH_SERVICE] User found, checking password")
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		println("[AUTH_SERVICE] Password mismatch:", err.Error())
		return nil, errors.New("invalid credentials")
	}

	println("[AUTH_SERVICE] Password correct, generating token")
	token, err := auth.GenerateToken(s.cfg.JWTSecret, user.ID, user.Role, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		Token: token,
		User:  user,
	}, nil
}

