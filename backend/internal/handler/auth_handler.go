package handler

import (
	"net/http"

	"backend/internal/auth"
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

func (h *AuthHandler) Me(c *gin.Context) {
	uidVal, exists := c.Get(auth.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	uid, ok := uidVal.(int64)
	if !ok {
		// gin may store numeric claims as float64 depending on parsing
		if f, okf := uidVal.(float64); okf {
			uid = int64(f)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id in context"})
			return
		}
	}

	user, err := h.svc.GetUserByID(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// return minimal public fields
	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"name":  user.Username,
		"email": user.Email,
	})
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
