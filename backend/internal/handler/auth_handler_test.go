package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/internal/config"
	"backend/internal/domain"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
)

type mockUserRepo struct {
	users map[string]*domain.User
}

func (m *mockUserRepo) Create(ctx interface{}, user *domain.User) error {
	user.ID = 1
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepo) GetByEmail(ctx interface{}, email string) (*domain.User, error) {
	return m.users[email], nil
}

func (m *mockUserRepo) GetByID(ctx interface{}, id int64) (*domain.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepo) List(ctx interface{}) ([]domain.User, error) {
	var list []domain.User
	for _, u := range m.users {
		list = append(list, *u)
	}
	return list, nil
}

func setupTestRouter() (*gin.Engine, *AuthHandler) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	cfg := &config.Config{JWTSecret: "test-secret"}
	repo := &mockUserRepo{users: make(map[string]*domain.User)}
	svc := service.NewAuthService(cfg, nil)
	handler := NewAuthHandler(svc)
	
	return router, handler
}

func TestRegisterHandler(t *testing.T) {
	router, handler := setupTestRouter()
	router.POST("/register", handler.register)

	body := map[string]string{
		"name":     "Test User",
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if response["token"] == nil {
		t.Error("Expected token in response")
	}
}

func TestLoginHandler(t *testing.T) {
	router, handler := setupTestRouter()
	router.POST("/register", handler.register)
	router.POST("/login", handler.login)

	// Register
	regBody := map[string]string{
		"name":     "Test User",
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonReg, _ := json.Marshal(regBody)
	reqReg := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonReg))
	reqReg.Header.Set("Content-Type", "application/json")
	wReg := httptest.NewRecorder()
	router.ServeHTTP(wReg, reqReg)

	// Login
	loginBody := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonLogin, _ := json.Marshal(loginBody)
	reqLogin := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonLogin))
	reqLogin.Header.Set("Content-Type", "application/json")
	wLogin := httptest.NewRecorder()
	router.ServeHTTP(wLogin, reqLogin)

	if wLogin.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", wLogin.Code)
	}
}

func TestLoginInvalidCredentials(t *testing.T) {
	router, handler := setupTestRouter()
	router.POST("/login", handler.login)

	body := map[string]string{
		"email":    "nonexistent@example.com",
		"password": "wrongpassword",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}
