package handlers

import (
	"html"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"exam-system/internal/models"
	"exam-system/internal/service"
)

type Handler struct {
	services *service.Services
}

func NewHandler(s *service.Services) *Handler {
	return &Handler{services: s}
}

func (h *Handler) RegisterRoutes(r *gin.Engine, jwtSecret string) {
	api := r.Group("/api")
	{
		api.POST("/signup", h.signup)
		api.POST("/login", h.login)
		api.POST("/logout", h.logout)

		api.POST("/profiles", h.createProfile)
		api.GET("/profiles/:id", h.getProfile)

		api.POST("/test/generate", h.GenerateTestHandler)
		api.POST("/test/submit", h.submitTest)

		api.GET("/analytics/group", h.GetGroupAnalytics)

		api.POST("/attempts", h.createAttempt)
		api.GET("/attempts/user/:userID", h.listAttemptsByUser)
	}
}

func (h *Handler) createProfile(c *gin.Context) {
	var req models.Profile
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": sanitizeError(err.Error())})
		return
	}
	if err := h.services.Profiles.CreateProfile(c.Request.Context(), &req); err != nil {
		log.Printf("Profile creation error: %v", sanitizeForLog(err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create profile"})
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *Handler) getProfile(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	profile, err := h.services.Profiles.GetProfile(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}
	c.JSON(http.StatusOK, profile)
}

type GenerateRequest struct {
	Subject    string `json:"subject" binding:"required"`
	Topic      string `json:"topic" binding:"required"`
	Difficulty string `json:"difficulty" binding:"required"`
	Lang       string `json:"lang" binding:"required"`
}

func (h *Handler) GenerateTestHandler(c *gin.Context) {
	var req GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": sanitizeError(err.Error())})
		return
	}

	if h.services.Gemini == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI service unavailable"})
		return
	}

	test, err := h.services.Gemini.GenerateTest(req.Subject, req.Topic, req.Lang, req.Difficulty)
	if err != nil {
		log.Printf("AI generation error: %v", sanitizeForLog(err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate test"})
		return
	}
	c.JSON(http.StatusOK, test)
}

type submitTestQuestion struct {
	ID            int    `json:"id"`
	Text          string `json:"text"`
	CorrectAnswer string `json:"correct_answer"`
	UserAnswer    string `json:"user_answer"`
}

type submitTestRequest struct {
	UserID    uint                 `json:"user_id" binding:"required"`
	Subject   string               `json:"subject" binding:"required"`
	Topic     string               `json:"topic" binding:"required"`
	Language  string               `json:"language" binding:"required"`
	Questions []submitTestQuestion `json:"questions" binding:"required"`
}

func (h *Handler) submitTest(c *gin.Context) {
	var req submitTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": sanitizeError(err.Error())})
		return
	}

	var correct int
	answerSummaries := make([]service.AnswerSummary, 0, len(req.Questions))

	for _, q := range req.Questions {
		userAns := strings.TrimSpace(strings.ToLower(q.UserAnswer))
		correctAns := strings.TrimSpace(strings.ToLower(q.CorrectAnswer))
		isCorrect := userAns != "" && userAns == correctAns

		if isCorrect {
			correct++
		}

		answerSummaries = append(answerSummaries, service.AnswerSummary{
			Question:      q.Text,
			CorrectAnswer: q.CorrectAnswer,
			UserAnswer:    q.UserAnswer,
			IsCorrect:     isCorrect,
		})
	}

	score := 0.0
	if len(req.Questions) > 0 {
		score = float64(correct) / float64(len(req.Questions)) * 100.0
	}

	advice := "AI service unavailable"
	if h.services.Gemini != nil {
		var err error
		advice, err = h.services.Gemini.AnalyzeTest(
			req.Subject,
			req.Topic,
			req.Language,
			answerSummaries,
			score,
		)
		if err != nil {
			log.Printf("AI analysis error: %v", sanitizeForLog(err.Error()))
			advice = "Failed to get AI analysis"
		}
	}

	attempt := &models.TestAttempt{
		UserID:     req.UserID,
		Subject:    req.Subject,
		Topic:      req.Topic,
		Score:      score,
		Language:   req.Language,
		AIFeedback: advice,
	}

	if err := h.services.TestAttempts.CreateAttempt(c.Request.Context(), attempt); err != nil {
		log.Printf("Failed to save attempt: %v", sanitizeForLog(err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save result"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"attempt":       attempt,
		"topic_mastery": advice,
	})
}

func (h *Handler) createAttempt(c *gin.Context) {
	var req models.TestAttempt
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": sanitizeError(err.Error())})
		return
	}
	if err := h.services.TestAttempts.CreateAttempt(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create attempt"})
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *Handler) listAttemptsByUser(c *gin.Context) {
	userIDParam := c.Param("userID")
	userID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}
	attempts, err := h.services.TestAttempts.ListAttemptsByUser(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch attempts"})
		return
	}
	c.JSON(http.StatusOK, attempts)
}

type GroupAnalyticsResponse struct {
	AverageScore float64  `json:"average_score"`
	WeakTopics   []string `json:"weak_topics"`
}

func (h *Handler) GetGroupAnalytics(c *gin.Context) {
	avg, weakTopics, err := h.services.TestAttempts.GetGroupAnalytics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch analytics"})
		return
	}

	resp := GroupAnalyticsResponse{
		AverageScore: avg,
		WeakTopics:   weakTopics,
	}
	c.JSON(http.StatusOK, resp)
}

type signupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Username string `json:"username" binding:"required"`
	Role     string `json:"role"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signup(c *gin.Context) {
	var req signupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": sanitizeError(err.Error())})
		return
	}

	user, err := h.services.Auth.Signup(c.Request.Context(), req.Email, req.Password, req.Username, req.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": sanitizeError(err.Error())})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
		"role":     user.Role,
	})
}

func (h *Handler) login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": sanitizeError(err.Error())})
		return
	}

	token, user, err := h.services.Auth.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":    token,
		"username": user.Username,
		"role":     user.Role,
	})
}

func (h *Handler) logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully. Remove token on client.",
	})
}

func sanitizeError(err string) string {
	return html.EscapeString(err)
}

func sanitizeForLog(input string) string {
	input = strings.ReplaceAll(input, "\n", " ")
	input = strings.ReplaceAll(input, "\r", " ")
	input = strings.ReplaceAll(input, "\t", " ")
	return strings.TrimSpace(input)
}
