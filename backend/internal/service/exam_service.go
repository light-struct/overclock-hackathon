package service

import (
	"context"
	"encoding/json"
	"fmt"

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
	AIProvider   string `json:"aiProvider"` // "gemini" or "groq"
}

// GenerateQuiz вызывает Gemini и возвращает уже распаршенный массив вопросов.
func (s *ExamService) GenerateQuiz(ctx context.Context, in GenerateQuizInput) (any, error) {
	difficultyGuide := map[string]string{
		"easy":   "Simple recall questions about basic facts and definitions.",
		"medium": "Applied knowledge questions requiring reasoning and connections between concepts.",
		"hard":   "Advanced analytical questions involving edge cases, synthesis, and deep understanding.",
	}[in.Difficulty]

	// Get knowledge base context for topic
	knowledgeContext := s.ai.GetKnowledgeForTopic(in.Topic)
	
	// Default to gemini if not specified
	aiProvider := in.AIProvider
	if aiProvider == "" {
		aiProvider = "gemini"
	}

	prompt := fmt.Sprintf(`Generate a quiz about "%s" with exactly %d questions at %s difficulty level.

%s

%s

You MUST respond with a valid JSON array only. No markdown, no explanation, no code blocks. Just a raw JSON array.

Each element must be an object with these exact fields:
- "questionNumber": integer starting from 1
- "question": the question text (string)
- "options": array of 4 answer options as strings (A, B, C, D) — even for true/false, provide 4 options
- "correctAnswer": the exact text of the correct option (must match one of the options exactly)

Return ONLY the JSON array.`, in.Topic, in.NumQuestions, in.Difficulty, difficultyGuide, knowledgeContext)

	text, err := s.ai.GenerateQuizQuestions(ctx, prompt, aiProvider)
	if err != nil {
		return nil, err
	}

	var questions any
	if err := json.Unmarshal([]byte(text), &questions); err != nil {
		return nil, fmt.Errorf("parse gemini response: %w", err)
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

