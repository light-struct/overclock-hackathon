package repository

import (
	"context"

	"backend/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProfileRepository struct {
	db *pgxpool.Pool
}

func NewProfileRepository(db *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (r *ProfileRepository) Create(ctx context.Context, p *domain.Profile) error {
	const q = `
		INSERT INTO profiles (full_name, role, preferred_lang)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at;
	`
	return r.db.QueryRow(ctx, q, p.FullName, p.Role, p.PreferredLang).
		Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *ProfileRepository) GetByID(ctx context.Context, id int64) (*domain.Profile, error) {
	const q = `
		SELECT id, full_name, role, preferred_lang, created_at, updated_at
		FROM profiles
		WHERE id = $1;
	`
	var p domain.Profile
	err := r.db.QueryRow(ctx, q, id).
		Scan(&p.ID, &p.FullName, &p.Role, &p.PreferredLang, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProfileRepository) List(ctx context.Context) ([]domain.Profile, error) {
	const q = `
		SELECT id, full_name, role, preferred_lang, created_at, updated_at
		FROM profiles
		ORDER BY id;
	`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.Profile
	for rows.Next() {
		var p domain.Profile
		if err := rows.Scan(&p.ID, &p.FullName, &p.Role, &p.PreferredLang, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, rows.Err()
}

func (r *ProfileRepository) Delete(ctx context.Context, id int64) error {
	const q = `DELETE FROM profiles WHERE id = $1;`
	_, err := r.db.Exec(ctx, q, id)
	return err
}

