package repository

import (
	"context"

	"exam-system/internal/models"

	"gorm.io/gorm"
)

type ProfileRepository interface {
	Create(ctx context.Context, p *models.Profile) error
	FindByID(ctx context.Context, id uint) (*models.Profile, error)
}

type TestAttemptRepository interface {
	Create(ctx context.Context, t *models.TestAttempt) error
	ListByUser(ctx context.Context, userID uint) ([]models.TestAttempt, error)
	GetOverallAverageScore(ctx context.Context) (float64, error)
	GetWeakTopicStats(ctx context.Context, threshold float64) ([]TopicStat, error)
}

// Repositories is the aggregate of all domain repositories.
// Exported so it can be used by the service layer.
type Repositories struct {
	Profiles     ProfileRepository
	TestAttempts TestAttemptRepository
}

type profileRepo struct {
	db *gorm.DB
}

type testAttemptRepo struct {
	db *gorm.DB
}

// TopicStat содержит средний балл по конкретной теме.
type TopicStat struct {
	Subject  string  `gorm:"column:subject"`
	Topic    string  `gorm:"column:topic"`
	AvgScore float64 `gorm:"column:avg_score"`
}

// NewRepositories wires all concrete repository implementations.
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Profiles:     &profileRepo{db: db},
		TestAttempts: &testAttemptRepo{db: db},
	}
}

func (r *profileRepo) Create(ctx context.Context, p *models.Profile) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *profileRepo) FindByID(ctx context.Context, id uint) (*models.Profile, error) {
	var profile models.Profile
	if err := r.db.WithContext(ctx).First(&profile, id).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *testAttemptRepo) Create(ctx context.Context, t *models.TestAttempt) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *testAttemptRepo) ListByUser(ctx context.Context, userID uint) ([]models.TestAttempt, error) {
	var attempts []models.TestAttempt
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&attempts).Error; err != nil {
		return nil, err
	}
	return attempts, nil
}

func (r *testAttemptRepo) GetOverallAverageScore(ctx context.Context) (float64, error) {
	var avg float64
	row := r.db.WithContext(ctx).
		Table("test_attempts").
		Select("AVG(score)").
		Row()
	if err := row.Scan(&avg); err != nil {
		return 0, err
	}
	return avg, nil
}

func (r *testAttemptRepo) GetWeakTopicStats(ctx context.Context, threshold float64) ([]TopicStat, error) {
	var stats []TopicStat
	if err := r.db.WithContext(ctx).
		Table("test_attempts").
		Select("subject, topic, AVG(score) AS avg_score").
		Group("subject, topic").
		Having("AVG(score) < ?", threshold).
		Scan(&stats).Error; err != nil {
		return nil, err
	}
	return stats, nil
}

