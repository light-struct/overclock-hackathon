package handler

import (
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
	userID, ok := userIDVal.(int64)
	if !ok {
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

func (h *ExamHandler) saveAttempt(c *gin.Context) {
	userIDVal, ok := c.Get(auth.ContextUserIDKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no user in context"})
		return
	}
	userID, ok := userIDVal.(int64)
	if !ok {
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

