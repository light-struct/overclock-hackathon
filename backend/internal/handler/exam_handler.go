package handler

import (
	"fmt"
	"net/http"

	"backend/internal/auth"
	"backend/internal/domain"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ExamHandler struct {
	svc *service.ExamService
}

func NewExamHandler(svc *service.ExamService) *ExamHandler {
	return &ExamHandler{svc: svc}
}

func (h *ExamHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/quiz/generate", h.generateQuiz)
	r.GET("/attempts/:id", h.getAttemptByID) // более специфичный маршрут первым
	r.GET("/attempts", h.listAttempts)
	r.POST("/attempts", h.saveAttempt)
}

func (h *ExamHandler) generateQuiz(c *gin.Context) {
	var req service.GenerateQuizInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	questions, err := h.svc.GenerateQuiz(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"questions": questions})
}

func (h *ExamHandler) listAttempts(c *gin.Context) {
	userIDVal, ok := c.Get(auth.ContextUserIDKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no user in context"})
		return
	}
	var userID int64
	switch v := userIDVal.(type) {
	case int64:
		userID = v
	case float64:
		userID = int64(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id type"})
		return
	}

	// Check if user is admin
	roleVal, _ := c.Get(auth.ContextRoleKey)
	role, _ := roleVal.(string)
	
	var attempts []domain.TestAttempt
	var err error
	
	if role == "admin" {
		attempts, err = h.svc.ListAllAttempts(c.Request.Context())
	} else {
		attempts, err = h.svc.ListAttemptsForUser(c.Request.Context(), userID)
	}
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"attempts": attempts})
}

func (h *ExamHandler) getAttemptByID(c *gin.Context) {
	userIDVal, ok := c.Get(auth.ContextUserIDKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no user in context"})
		return
	}
	var userID int64
	switch v := userIDVal.(type) {
	case int64:
		userID = v
	case float64:
		userID = int64(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id type"})
		return
	}
	idParam := c.Param("id")
	var attemptID int64
	if _, err := fmt.Sscanf(idParam, "%d", &attemptID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid attempt id"})
		return
	}
	roleVal, _ := c.Get(auth.ContextRoleKey)
	role, _ := roleVal.(string)
	isAdmin := role == "admin" || role == "teacher"

	var att *domain.TestAttempt
	var err error
	if isAdmin {
		att, err = h.svc.GetAttemptByIDAny(c.Request.Context(), attemptID)
	} else {
		att, err = h.svc.GetAttemptByID(c.Request.Context(), attemptID, userID)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if att == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "attempt not found"})
		return
	}
	c.JSON(http.StatusOK, att)
}

func (h *ExamHandler) saveAttempt(c *gin.Context) {
	userIDVal, ok := c.Get(auth.ContextUserIDKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no user in context"})
		return
	}
	var userID int64
	switch v := userIDVal.(type) {
	case int64:
		userID = v
	case float64:
		userID = int64(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id type"})
		return
	}

	var req service.SaveAttemptInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	req.UserID = userID

	att, err := h.svc.SaveAttempt(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, att)
}

