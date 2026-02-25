package repository

import (
	"context"
	"errors"

	"backend/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *domain.User) error {
	const q = `
		INSERT INTO users (email, username, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at;
	`
	return r.db.QueryRow(ctx, q, u.Email, u.Username, u.PasswordHash, u.Role).
		Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const q = `
		SELECT id, email, username, password_hash, role, created_at, updated_at
		FROM users
		WHERE email = $1;
	`
	var u domain.User
	err := r.db.QueryRow(ctx, q, email).
		Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	const q = `
		SELECT id, email, username, password_hash, role, created_at, updated_at
		FROM users
		WHERE id = $1;
	`
	var u domain.User
	err := r.db.QueryRow(ctx, q, id).
		Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) List(ctx context.Context) ([]domain.User, error) {
	const q = `
		SELECT id, email, username, role, created_at, updated_at
		FROM users
		ORDER BY id;
	`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Username, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		res = append(res, u)
	}
	return res, rows.Err()
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	const q = `DELETE FROM users WHERE id = $1;`
	_, err := r.db.Exec(ctx, q, id)
	return err
}

