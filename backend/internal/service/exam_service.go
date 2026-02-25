package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"backend/internal/domain"
	"backend/internal/repository"
)

type ExamService struct {
	ai       *AIService
	attempts *repository.TestAttemptRepository
}

func NewExamService(ai *AIService, attempts *repository.TestAttemptRepository) *ExamService {
	return &ExamService{
		ai:       ai,
		attempts: attempts,
	}
}

type GenerateQuizInput struct {
	Topic        string `json:"topic"`
	NumQuestions int    `json:"numQuestions"`
	Difficulty   string `json:"difficulty"`
}

// GenerateQuiz вызывает Gemini и возвращает уже распаршенный массив вопросов.
func (s *ExamService) GenerateQuiz(ctx context.Context, in GenerateQuizInput) (any, error) {
	difficultyGuide := map[string]string{
		"easy":   "Simple recall questions about basic facts and definitions.",
		"medium": "Applied knowledge questions requiring reasoning and connections between concepts.",
		"hard":   "Advanced analytical questions involving edge cases, synthesis, and deep understanding.",
	}[in.Difficulty]

	prompt := fmt.Sprintf(`Сгенерируй тест по теме "%s" из %d вопросов уровня сложности %s.

%s

ИСПОЛЬЗУЙ ТОЛЬКО ВОПРОСЫ ИЗ БАЗЫ ДАННЫХ (ai-knowledge-base/questions/).
Выбери вопросы случайно из соответствующего файла и уровня сложности.

ВЕРНИ ТОЛЬКО ВАЛИДНЫЙ JSON МАССИВ БЕЗ MARKDOWN, БЕЗ ОБЪЯСНЕНИЙ.

Формат каждого элемента:
{
  "questionNumber": 1,
  "question": "текст вопроса",
  "options": ["A) вариант", "B) вариант", "C) вариант", "D) вариант"],
  "correctAnswer": "точный текст правильного варианта"
}

ВЕРНИ ТОЛЬКО JSON МАССИВ.`, in.Topic, in.NumQuestions, in.Difficulty, difficultyGuide)

	text, err := s.ai.GenerateQuizQuestions(ctx, in.Topic, in.Difficulty, prompt)
	if err != nil {
		return nil, err
	}

	fmt.Printf("[DEBUG] Gemini response (first 500 chars): %s\n", text[:min(len(text), 500)])

	// Очистка от markdown блоков
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "```") {
		lines := strings.Split(text, "\n")
		if len(lines) > 2 {
			text = strings.Join(lines[1:len(lines)-1], "\n")
		}
	}
	text = strings.TrimSpace(text)

	var questions any
	if err := json.Unmarshal([]byte(text), &questions); err != nil {
		return nil, fmt.Errorf("parse gemini response: %w (response: %s)", err, text[:min(len(text), 200)])
	}
	return questions, nil
}

type SaveAttemptInput struct {
	UserID   int64   `json:"user_id"`
	Subject  string  `json:"subject"`
	Topic    string  `json:"topic"`
	Score    float64 `json:"score"`
	Language string  `json:"language"`
	Feedback string  `json:"ai_feedback"`
}

func (s *ExamService) SaveAttempt(ctx context.Context, in SaveAttemptInput) (*domain.TestAttempt, error) {
	att := &domain.TestAttempt{
		UserID:     in.UserID,
		Subject:    in.Subject,
		Topic:      in.Topic,
		Score:      in.Score,
		Language:   in.Language,
		AIFeedback: in.Feedback,
	}
	if err := s.attempts.Create(ctx, att); err != nil {
		return nil, err
	}
	return att, nil
}

func (s *ExamService) ListAttemptsForUser(ctx context.Context, userID int64) ([]domain.TestAttempt, error) {
	return s.attempts.ListByUser(ctx, userID)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
