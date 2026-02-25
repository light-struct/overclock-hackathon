package service

import (
	"context"
	"fmt"

	"backend/internal/config"

	"google.golang.org/genai"
)

// AIService инкапсулирует работу с Gemini.
type AIService struct {
	client *genai.Client
}

func NewAIService(ctx context.Context, cfg *config.Config) (*AIService, error) {
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("create genai client: %w", err)
	}
	return &AIService{client: client}, nil
}

// GenerateQuizQuestions проксирует промпт к Gemini и возвращает текстовый ответ.
// Фронтенд уже ожидает JSON‑массив вопросов, поэтому здесь мы просто возвращаем текст,
// а парсинг делает клиент.
func (s *AIService) GenerateQuizQuestions(ctx context.Context, prompt string) (string, error) {
	res, err := s.client.Models.GenerateContent(
		ctx,
		"gemini-3-flash-preview",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", err
	}
	return res.Text(), nil
}

