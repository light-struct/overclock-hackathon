package handler

import (
	"net/http"

	"backend/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/register", h.register)
	r.POST("/login", h.login)
	r.POST("/logout", h.logout)
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	res, err := h.svc.Register(c.Request.Context(), service.RegisterInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		println("[LOGIN] Failed to bind JSON:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	println("[LOGIN] Attempting login for email:", req.Email)
	res, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		println("[LOGIN] Login failed:", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	println("[LOGIN] Login successful for user:", res.User.Email)
	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) logout(c *gin.Context) {
	// В JWT‑подходе на сервере нет состояния — клиент просто забывает токен.
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

