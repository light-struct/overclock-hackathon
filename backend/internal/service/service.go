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
	Auth         AuthService
}

type ProfileService interface {
	CreateProfile(ctx context.Context, p *models.Profile) error
	GetProfile(ctx context.Context, id uint) (*models.Profile, error)
}

type TestAttemptService interface {
	CreateAttempt(ctx context.Context, t *models.TestAttempt) error
	ListAttemptsByUser(ctx context.Context, userID uint) ([]models.TestAttempt, error)
	GetGroupAnalytics(ctx context.Context) (float64, []string, error)
}

type profileService struct {
	repo repository.ProfileRepository
}

type testAttemptService struct {
	repo repository.TestAttemptRepository
}

func NewServices(r *repository.Repositories, geminiKey, jwtSecret string) (*Services, error) {
	geminiService, err := NewGeminiService(context.Background(), geminiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Gemini: %w", err)
	}

	auth := NewAuthService(r.Users, jwtSecret)

	return &Services{
		Profiles:     &profileService{repo: r.Profiles},
		TestAttempts: &testAttemptService{repo: r.TestAttempts},
		Gemini:       geminiService,
		Auth:         auth,
	}, nil
}

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