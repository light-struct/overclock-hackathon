package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiService interface {
	GenerateTest(subject, topic, lang, difficulty string) (*GeneratedTest, error)
	AnalyzeTest(subject, topic, lang string, answers []AnswerSummary, score float64) (string, error)
}

type geminiService struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewGeminiService(ctx context.Context, apiKey string) (GeminiService, error) {
	if apiKey == "" {
		return nil, errors.New("GEMINI_API_KEY is required")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	model := client.GenerativeModel("gemini-1.5-flash")
	model.ResponseMIMEType = "application/json"

	return &geminiService{client: client, model: model}, nil
}

type GeneratedTest struct {
	TestTitle string `json:"test_title"`
	Questions []struct {
		ID      int      `json:"id"`
		Text    string   `json:"text"`
		Options []string `json:"options"`
		Answer  string   `json:"answer"`
	} `json:"questions"`
}

type AnswerSummary struct {
	Question      string `json:"question"`
	CorrectAnswer string `json:"correct_answer"`
	UserAnswer    string `json:"user_answer"`
	IsCorrect     bool   `json:"is_correct"`
}

func (s *geminiService) GenerateTest(subject, topic, lang, difficulty string) (*GeneratedTest, error) {
	ctx := context.Background()

	prompt := fmt.Sprintf(`You are a strict examiner.
Create a test for subject: '%s', topic: '%s'.
Difficulty: %s.
Language: %s (IMPORTANT: If 'kk' - write in Kazakh, if 'ru' - in Russian, if 'en' - in English).
Requirements:
1. Exactly 5 questions.
2. 4 answer options for each question.
3. Output format: ONLY valid JSON.
JSON Schema: {"test_title": "Title", "questions": [{"id": 1, "text": "Question?", "options": ["A", "B", "C", "D"], "answer": "Correct answer text"}]}`,
		sanitizeInput(subject), sanitizeInput(topic), sanitizeInput(difficulty), sanitizeInput(lang))

	resp, err := s.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Printf("AI generation error (sanitized): %v", sanitizeForLog(err.Error()))
		return nil, fmt.Errorf("AI request failed: %w", err)
	}

	return parseGeminiResponse[GeneratedTest](resp)
}

func (s *geminiService) AnalyzeTest(subject, topic, lang string, answers []AnswerSummary, score float64) (string, error) {
	ctx := context.Background()
	answersJson, _ := json.Marshal(answers)

	prompt := fmt.Sprintf(`Role: You are a wise mentor.
Subject: %s, Topic: %s.
Response language: %s (kk - Kazakh, ru - Russian, en - English).
Student scored: %.0f points out of 100.
Student answers: %s
Your task:
1. Briefly praise or support.
2. Explain the main mistake (if any).
3. Give advice on what to read.
Return JSON: {"feedback": "Your text here..."}`,
		sanitizeInput(subject), sanitizeInput(topic), sanitizeInput(lang), score, string(answersJson))

	resp, err := s.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Printf("AI analysis error (sanitized): %v", sanitizeForLog(err.Error()))
		return "Failed to get AI analysis", err
	}

	res, err := parseGeminiResponse[struct{ Feedback string `json:"feedback"` }](resp)
	if err != nil {
		if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
			return fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0]), nil
		}
		return "Analysis error", nil
	}
	return res.Feedback, nil
}

func parseGeminiResponse[T any](resp *genai.GenerateContentResponse) (*T, error) {
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, errors.New("empty AI response")
	}
	part := resp.Candidates[0].Content.Parts[0]
	jsonText := fmt.Sprintf("%v", part)

	jsonText = strings.TrimPrefix(jsonText, "```json")
	jsonText = strings.TrimSuffix(jsonText, "```")
	jsonText = strings.TrimSpace(jsonText)

	var result T
	if err := json.Unmarshal([]byte(jsonText), &result); err != nil {
		return nil, fmt.Errorf("JSON parse error: %w", err)
	}
	return &result, nil
}

func sanitizeInput(input string) string {
	input = strings.ReplaceAll(input, "\n", " ")
	input = strings.ReplaceAll(input, "\r", " ")
	input = strings.ReplaceAll(input, "\t", " ")
	return strings.TrimSpace(input)
}

func sanitizeForLog(input string) string {
	input = strings.ReplaceAll(input, "\n", " ")
	input = strings.ReplaceAll(input, "\r", " ")
	input = strings.ReplaceAll(input, "\t", " ")
	return strings.TrimSpace(input)
}