package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"backend/internal/config"
)

type AIService struct {
	apiKey            string
	systemInstruction string
	questionsDir      string
}

func NewAIService(ctx context.Context, cfg *config.Config) (*AIService, error) {
	if cfg.GeminiAPIKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set")
	}

	// Загрузка базы знаний
	aiConfigPath := filepath.Join("..", "ai-knowledge-base", "ai-config.md")
	aiConfig, err := os.ReadFile(aiConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load ai-config.md: %w", err)
	}

	questionsDir := filepath.Join("..", "ai-knowledge-base", "questions")

	return &AIService{
		apiKey:            cfg.GeminiAPIKey,
		systemInstruction: string(aiConfig),
		questionsDir:      questionsDir,
	}, nil
}

func (s *AIService) GenerateQuizQuestions(ctx context.Context, topic, difficulty string, prompt string) (string, error) {
	// Маппинг тем на файлы
	topicMap := map[string]string{
		"oop":              "oop",
		"oop-java":         "oop-java",
		"java-core":        "java-core",
		"procedural-java":  "procedural-java",
		"exceptions":       "exceptions",
		"serialization":    "serialization",
		"garbage-collector": "garbage-collector",
	}

	fileName, ok := topicMap[topic]
	if !ok {
		fileName = topic // фолбэк на исходное имя
	}

	// Загружаем только нужный файл с вопросами
	topicFile := filepath.Join(s.questionsDir, fileName+".md")
	topicQuestions, err := os.ReadFile(topicFile)
	if err != nil {
		return "", fmt.Errorf("failed to load questions for topic %s: %w", topic, err)
	}

	// Добавляем контекст темы в промпт
	fullPrompt := fmt.Sprintf("=== ВОПРОСЫ ПО ТЕМЕ %s ===\n%s\n\n=== ЗАДАНИЕ ===\n%s", topic, string(topicQuestions), prompt)

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=%s", s.apiKey)
	
	reqBody := map[string]interface{}{
		"systemInstruction": map[string]interface{}{
			"parts": []map[string]string{{"text": s.systemInstruction}},
		},
		"contents": []map[string]interface{}{
			{"parts": []map[string]string{{"text": fullPrompt}}},
		},
		"generationConfig": map[string]interface{}{
			"temperature":       0.7,
			"maxOutputTokens":   4096,
			"responseMimeType": "application/json",
		},
	}
	
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	fmt.Printf("[DEBUG] Raw response status: %d, body length: %d\n", resp.StatusCode, len(body))
	
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("gemini API error %d: %s", resp.StatusCode, string(body))
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}
	
	candidates, ok := result["candidates"].([]interface{})
	if !ok || len(candidates) == 0 {
		return "", fmt.Errorf("no candidates in response")
	}
	
	candidate := candidates[0].(map[string]interface{})
	content := candidate["content"].(map[string]interface{})
	parts := content["parts"].([]interface{})
	part := parts[0].(map[string]interface{})
	text := part["text"].(string)
	
	return text, nil
}
