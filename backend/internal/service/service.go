package service

import (
	"context"
	"exam-system/internal/models"
	"exam-system/internal/repository"
	"fmt"
)

type Services struct {
	Profiles     ProfileService
	TestAttempts TestAttemptService
	Gemini       GeminiService
}

// Интерфейсы
type ProfileService interface {
	CreateProfile(ctx context.Context, p *models.Profile) error
	GetProfile(ctx context.Context, id uint) (*models.Profile, error)
}

type TestAttemptService interface {
	CreateAttempt(ctx context.Context, t *models.TestAttempt) error
	ListAttemptsByUser(ctx context.Context, userID uint) ([]models.TestAttempt, error)
	GetGroupAnalytics(ctx context.Context) (float64, []string, error)
}

// Реализация
type profileService struct {
	repo repository.ProfileRepository
}

type testAttemptService struct {
	repo repository.TestAttemptRepository
}

// КОНСТРУКТОР (Исправленная версия)
func NewServices(r *repository.Repositories) *Services {
	// Инициализируем Gemini прямо здесь
	geminiService, err := NewGeminiServiceFromEnv(context.Background())
	if err != nil {
		fmt.Printf("⚠️ Ошибка инициализации Gemini: %v\n", err)
		// В продакшене тут надо падать, но для хакатона продолжим, чтобы не крашить сервер
	}

	return &Services{
		Profiles: &profileService{
			repo: r.Profiles,
		},
		TestAttempts: &testAttemptService{
			repo: r.TestAttempts,
		},
		Gemini: geminiService,
	}
}

// Методы сервисов
func (s *profileService) CreateProfile(ctx context.Context, p *models.Profile) error {
	return s.repo.Create(ctx, p)
}

func (s *profileService) GetProfile(ctx context.Context, id uint) (*models.Profile, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *testAttemptService) CreateAttempt(ctx context.Context, t *models.TestAttempt) error {
	return s.repo.Create(ctx, t)
}

func (s *testAttemptService) ListAttemptsByUser(ctx context.Context, userID uint) ([]models.TestAttempt, error) {
	return s.repo.ListByUser(ctx, userID)
}

// GetGroupAnalytics считает средний балл и слабые темы для группы.
func (s *testAttemptService) GetGroupAnalytics(ctx context.Context) (float64, []string, error) {
	avg, err := s.repo.GetOverallAverageScore(ctx)
	if err != nil {
		return 0, nil, err
	}

	stats, err := s.repo.GetWeakTopicStats(ctx, 50.0)
	if err != nil {
		return 0, nil, err
	}

	weakTopics := make([]string, 0, len(stats))
	for _, st := range stats {
		weakTopics = append(weakTopics, fmt.Sprintf("%s/%s", st.Subject, st.Topic))
	}

	return avg, weakTopics, nil
}