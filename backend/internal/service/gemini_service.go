package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// Интерфейс
type GeminiService interface {
	GenerateTest(subject, topic, lang, difficulty string) (*GeneratedTest, error)
	AnalyzeTest(subject, topic, lang string, answers []AnswerSummary, score float64) (string, error)
}

// Реализация
type geminiService struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

// Конструктор (Внимание: имя функции должно совпадать с вызовом в service.go)
func NewGeminiServiceFromEnv(ctx context.Context) (GeminiService, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		// Если ключа нет, вернем ошибку или заглушку (для хакатона лучше ошибку в лог, но продолжить работу)
		fmt.Println("⚠️ WARNING: GEMINI_API_KEY not found in env")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	// Используем Flash модель
	model := client.GenerativeModel("gemini-1.5-flash")
	// Важно: настраиваем модель на ответ JSON'ом
	model.ResponseMIMEType = "application/json"

	return &geminiService{client: client, model: model}, nil
}

// Структуры данных
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

// --- ЛОГИКА 1: ГЕНЕРАЦИЯ ТЕСТА ---
func (s *geminiService) GenerateTest(subject, topic, lang, difficulty string) (*GeneratedTest, error) {
	ctx := context.Background()

	// Формируем промпт через конкатенацию, чтобы избежать ошибок синтаксиса с кавычками
	prompt := "Ты строгий экзаменатор.\n"
	prompt += fmt.Sprintf("Создай тест по предмету: '%s', тема: '%s'.\n", subject, topic)
	prompt += fmt.Sprintf("Сложность: %s.\n", difficulty)
	prompt += fmt.Sprintf("Язык: %s (ВАЖНО: Если 'kk' - пиши на Казахском, если 'ru' - на Русском).\n", lang)
	prompt += "Требования:\n"
	prompt += "1. Ровно 5 вопросов.\n"
	prompt += "2. 4 варианта ответа на каждый вопрос.\n"
	prompt += "3. Формат вывода: ТОЛЬКО валидный JSON.\n"
	prompt += "JSON Схема: { \"test_title\": \"Название\", \"questions\": [ { \"id\": 1, \"text\": \"Вопрос?\", \"options\": [\"A\", \"B\", \"C\", \"D\"], \"answer\": \"Текст правильного ответа\" } ] }"

	resp, err := s.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к AI: %v", err)
	}

	return parseGeminiResponse[GeneratedTest](resp)
}

// --- ЛОГИКА 2: АНАЛИЗ РЕЗУЛЬТАТОВ ---
func (s *geminiService) AnalyzeTest(subject, topic, lang string, answers []AnswerSummary, score float64) (string, error) {
	ctx := context.Background()
	answersJson, _ := json.Marshal(answers)

	prompt := "Роль: Ты мудрый наставник.\n"
	prompt += fmt.Sprintf("Предмет: %s, Тема: %s.\n", subject, topic)
	prompt += fmt.Sprintf("Язык ответа: %s (kk - казахский, ru - русский).\n", lang)
	prompt += fmt.Sprintf("Студент набрал: %.0f баллов из 100.\n", score)
	prompt += fmt.Sprintf("Ответы студента: %s\n", string(answersJson))
	prompt += "Твоя задача:\n"
	prompt += "1. Кратко похвали или поддержи.\n"
	prompt += "2. Объясни главную ошибку (если есть).\n"
	prompt += "3. Дай совет, что почитать.\n"
	prompt += "Верни JSON: { \"feedback\": \"Твой текст здесь...\" }"

	resp, err := s.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	// Используем анонимную структуру для парсинга
	res, err := parseGeminiResponse[struct{ Feedback string `json:"feedback"` }](resp)
	if err != nil {
		// Если AI вернул просто текст, а не JSON, попробуем вернуть raw текст
		if len(resp.Candidates) > 0 {
			part := resp.Candidates[0].Content.Parts[0]
			return fmt.Sprintf("%v", part), nil
		}
		return "Ошибка анализа ответов", nil
	}
	return res.Feedback, nil
}

// Вспомогательная функция (Generics)
func parseGeminiResponse[T any](resp *genai.GenerateContentResponse) (*T, error) {
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("пустой ответ от AI")
	}
	part := resp.Candidates[0].Content.Parts[0]
	jsonText := fmt.Sprintf("%v", part)
	
	// Чистим мусор (markdown)
	jsonText = strings.TrimPrefix(jsonText, "```json")
	jsonText = strings.TrimSuffix(jsonText, "```")
	jsonText = strings.TrimSpace(jsonText)

	var result T
	if err := json.Unmarshal([]byte(jsonText), &result); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %v. Текст от AI: %s", err, jsonText)
	}
	return &result, nil
}