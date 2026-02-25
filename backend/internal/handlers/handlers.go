package handlers

import (
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

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// Профили
		api.POST("/profiles", h.createProfile)
		api.GET("/profiles/:id", h.getProfile)

		// Тесты
		api.POST("/test/generate", h.GenerateTestHandler) // Генерация
		api.POST("/test/submit", h.submitTest)            // Сдача и анализ

		// Аналитика для учителя
		api.GET("/analytics/group", h.GetGroupAnalytics)

		// Попытки (история)
		api.POST("/attempts", h.createAttempt)
		api.GET("/attempts/user/:userID", h.listAttemptsByUser)
	}
}

// --- ХЕНДЛЕРЫ ПРОФИЛЕЙ ---
func (h *Handler) createProfile(c *gin.Context) {
	var req models.Profile
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.services.Profiles.CreateProfile(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}

// --- ХЕНДЛЕРЫ ГЕНЕРАЦИИ ТЕСТА ---
type GenerateRequest struct {
	Subject    string `json:"subject"`
	Topic      string `json:"topic"`
	Difficulty string `json:"difficulty"`
	Lang       string `json:"lang"`
}

func (h *Handler) GenerateTestHandler(c *gin.Context) {
	var req GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.services.Gemini == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI service is not initialized"})
		return
	}

	test, err := h.services.Gemini.GenerateTest(req.Subject, req.Topic, req.Lang, req.Difficulty)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI error: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, test)
}

// --- ХЕНДЛЕРЫ СДАЧИ ТЕСТА ---
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var correct int
	answerSummaries := make([]service.AnswerSummary, 0, len(req.Questions))

	for _, q := range req.Questions {
		// Сравнение без учета регистра и пробелов
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

	// Анализ через AI
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
			// Логируем ошибку, но не валим запрос пользователю
			println("AI Analysis failed:", err.Error())
			advice = "Не удалось получить анализ от ИИ."
		}
	}

	// Сохранение в БД
	attempt := &models.TestAttempt{
		UserID:     req.UserID,
		Subject:    req.Subject,
		Topic:      req.Topic,
		Score:      score,
		Language:   req.Language,
		AIFeedback: advice,
	}

	if err := h.services.TestAttempts.CreateAttempt(c.Request.Context(), attempt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save result", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"attempt":       attempt,
		"topic_mastery": advice,
	})
}

// Прочие методы
func (h *Handler) createAttempt(c *gin.Context) {
	var req models.TestAttempt
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.services.TestAttempts.CreateAttempt(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, attempts)
}

// --- Аналитика для учителя ---
type GroupAnalyticsResponse struct {
	AverageScore float64  `json:"average_score"`
	WeakTopics   []string `json:"weak_topics"`
}

// GetGroupAnalytics возвращает средний балл и слабые темы по всем попыткам.
func (h *Handler) GetGroupAnalytics(c *gin.Context) {
	avg, weakTopics, err := h.services.TestAttempts.GetGroupAnalytics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := GroupAnalyticsResponse{
		AverageScore: avg,
		WeakTopics:   weakTopics,
	}
	c.JSON(http.StatusOK, resp)
}