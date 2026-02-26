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
	geminiAPIKey   string
	groqAPIKey     string
	systemPrompt   string
	knowledgeBase  map[string]string
}

func NewAIService(ctx context.Context, cfg *config.Config) (*AIService, error) {
	if cfg.GeminiAPIKey == "" && cfg.GroqAPIKey == "" {
		return nil, fmt.Errorf("at least one AI API key must be set")
	}
	
	svc := &AIService{
		geminiAPIKey:  cfg.GeminiAPIKey,
		groqAPIKey:    cfg.GroqAPIKey,
		knowledgeBase: make(map[string]string),
	}
	
	if err := svc.loadKnowledgeBase(); err != nil {
		return nil, fmt.Errorf("failed to load knowledge base: %w", err)
	}
	
	return svc, nil
}

func (s *AIService) loadKnowledgeBase() error {
	basePath := "../ai-knowledge-base"
	
	// Load system prompt
	configPath := filepath.Join(basePath, "ai-config.md")
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("read ai-config.md: %w", err)
	}
	s.systemPrompt = string(configData)
	
	// Load question files
	questionsPath := filepath.Join(basePath, "questions")
	files, err := os.ReadDir(questionsPath)
	if err != nil {
		return fmt.Errorf("read questions dir: %w", err)
	}
	
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".md" {
			topic := file.Name()[:len(file.Name())-3]
			data, err := os.ReadFile(filepath.Join(questionsPath, file.Name()))
			if err != nil {
				return fmt.Errorf("read %s: %w", file.Name(), err)
			}
			s.knowledgeBase[topic] = string(data)
		}
	}
	
	return nil
}

func (s *AIService) GenerateQuizQuestions(ctx context.Context, prompt string, aiProvider string) (string, error) {
	if aiProvider == "groq" {
		return s.generateWithGroq(ctx, prompt)
	}
	return s.generateWithGemini(ctx, prompt)
}

func (s *AIService) generateWithGemini(ctx context.Context, prompt string) (string, error) {
	if s.geminiAPIKey == "" {
		return "", fmt.Errorf("Gemini API key not configured")
	}
	
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=%s", s.geminiAPIKey)
	
	// Build full context with system prompt
	fullPrompt := s.systemPrompt + "\n\n" + prompt
	
	reqBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{"parts": []map[string]string{{"text": fullPrompt}}},
		},
		"systemInstruction": map[string]interface{}{
			"parts": []map[string]string{{"text": "You are QuizAgent, an AI testing platform. Follow all rules from ai-config.md strictly."}},
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
	
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("gemini API error %d: %s", resp.StatusCode, string(body))
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
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

func (s *AIService) GetKnowledgeForTopic(topic string) string {
	if content, ok := s.knowledgeBase[topic]; ok {
		return fmt.Sprintf("Knowledge base for topic '%s':\n%s", topic, content)
	}
	return fmt.Sprintf("No specific knowledge base found for topic '%s'. Generate questions based on general knowledge.", topic)
}

func (s *AIService) generateWithGroq(ctx context.Context, prompt string) (string, error) {
	if s.groqAPIKey == "" {
		return "", fmt.Errorf("Groq API key not configured")
	}
	
	url := "https://api.groq.com/openai/v1/chat/completions"
	
	fullPrompt := s.systemPrompt + "\n\n" + prompt
	
	reqBody := map[string]interface{}{
		"model": "llama-3.3-70b-versatile",
		"messages": []map[string]string{
			{"role": "system", "content": "You are QuizAgent, an AI testing platform. Follow all rules from ai-config.md strictly."},
			{"role": "user", "content": fullPrompt},
		},
		"temperature": 0.7,
	}
	
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.groqAPIKey)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("groq API error %d: %s", resp.StatusCode, string(body))
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	
	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}
	
	choice := choices[0].(map[string]interface{})
	message := choice["message"].(map[string]interface{})
	text := message["content"].(string)
	
	return text, nil
}
