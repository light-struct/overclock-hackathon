package service

import (
	"context"
	"testing"

	"backend/internal/config"
	"backend/internal/domain"
	"backend/internal/repository"
)

type mockUserRepo struct {
	users map[string]*domain.User
}

func (m *mockUserRepo) Create(ctx context.Context, user *domain.User) error {
	user.ID = 1
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return m.users[email], nil
}

func (m *mockUserRepo) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepo) List(ctx context.Context) ([]domain.User, error) {
	var list []domain.User
	for _, u := range m.users {
		list = append(list, *u)
	}
	return list, nil
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*domain.User)}
}

func TestRegister(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	repo := newMockUserRepo()
	svc := NewAuthService(cfg, (*repository.UserRepository)(nil))
	svc.userRepo = (*repository.UserRepository)(repo)

	result, err := svc.Register(context.Background(), RegisterInput{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	})

	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	if result.Token == "" {
		t.Error("Expected token, got empty string")
	}
	if result.User.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", result.User.Email)
	}
}

func TestLogin(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	repo := newMockUserRepo()
	svc := NewAuthService(cfg, (*repository.UserRepository)(nil))
	svc.userRepo = (*repository.UserRepository)(repo)

	// Register first
	_, err := svc.Register(context.Background(), RegisterInput{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Login
	result, err := svc.Login(context.Background(), "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if result.Token == "" {
		t.Error("Expected token")
	}
}

func TestLoginInvalidPassword(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	repo := newMockUserRepo()
	svc := NewAuthService(cfg, (*repository.UserRepository)(nil))
	svc.userRepo = (*repository.UserRepository)(repo)

	svc.Register(context.Background(), RegisterInput{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	})

	_, err := svc.Login(context.Background(), "test@example.com", "wrongpassword")
	if err == nil {
		t.Error("Expected error for invalid password")
	}
}
