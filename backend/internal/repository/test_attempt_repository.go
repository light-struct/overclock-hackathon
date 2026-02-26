package repository

import (
	"context"

	"backend/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TestAttemptRepository struct {
	db *pgxpool.Pool
}

func NewTestAttemptRepository(db *pgxpool.Pool) *TestAttemptRepository {
	return &TestAttemptRepository{db: db}
}

func (r *TestAttemptRepository) Create(ctx context.Context, t *domain.TestAttempt) error {
	const q = `
		INSERT INTO test_attempts (user_id, subject, topic, score, language, ai_feedback)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at;
	`
	return r.db.QueryRow(ctx, q,
		t.UserID, t.Subject, t.Topic, t.Score, t.Language, t.AIFeedback,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *TestAttemptRepository) GetByID(ctx context.Context, id int64) (*domain.TestAttempt, error) {
	const q = `
		SELECT id, user_id, subject, topic, score, language, ai_feedback, created_at, updated_at
		FROM test_attempts
		WHERE id = $1;
	`
	var t domain.TestAttempt
	err := r.db.QueryRow(ctx, q, id).
		Scan(&t.ID, &t.UserID, &t.Subject, &t.Topic, &t.Score, &t.Language, &t.AIFeedback, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TestAttemptRepository) ListByUser(ctx context.Context, userID int64) ([]domain.TestAttempt, error) {
	const q = `
		SELECT id, user_id, subject, topic, score, language, ai_feedback, created_at, updated_at
		FROM test_attempts
		WHERE user_id = $1
		ORDER BY created_at DESC;
	`
	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.TestAttempt
	for rows.Next() {
		var t domain.TestAttempt
		if err := rows.Scan(&t.ID, &t.UserID, &t.Subject, &t.Topic, &t.Score, &t.Language, &t.AIFeedback, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	return res, rows.Err()
}

func (r *TestAttemptRepository) Delete(ctx context.Context, id int64) error {
	const q = `DELETE FROM test_attempts WHERE id = $1;`
	_, err := r.db.Exec(ctx, q, id)
	return err
}


func (r *TestAttemptRepository) ListAll(ctx context.Context) ([]domain.TestAttempt, error) {
	const q = `
		SELECT id, user_id, subject, topic, score, language, ai_feedback, created_at, updated_at
		FROM test_attempts
		ORDER BY created_at DESC;
	`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.TestAttempt
	for rows.Next() {
		var t domain.TestAttempt
		if err := rows.Scan(&t.ID, &t.UserID, &t.Subject, &t.Topic, &t.Score, &t.Language, &t.AIFeedback, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	return res, rows.Err()
}
