package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config содержит все основные переменные окружения, нужные бэкенду.
type Config struct {
	Port         string
	DatabaseURL  string
	JWTSecret    string
	GeminiAPIKey string
}

// Load загружает .env и формирует структуру конфигурации.
func Load() *Config {
	// Пытаемся загрузить .env как из корня backend, так и из корня репо.
	_ = godotenv.Load(".env")
	_ = godotenv.Load("backend/.env")

	cfg := &Config{
		Port:         getEnv("PORT", "8080"),
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	// Для клиента google.golang.org/genai ожидается GOOGLE_API_KEY.
	// Пробрасываем туда GEMINI_API_KEY, если он задан.
	if cfg.GeminiAPIKey != "" && os.Getenv("GOOGLE_API_KEY") == "" {
		if err := os.Setenv("GOOGLE_API_KEY", cfg.GeminiAPIKey); err != nil {
			log.Printf("warning: failed to set GOOGLE_API_KEY from GEMINI_API_KEY: %v", err)
		}
	}

	return cfg
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

