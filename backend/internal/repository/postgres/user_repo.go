package postgres

import (
	"context"
	"errors"

	"github.com/soydoradesu/product_discovery/internal/domain"
	"github.com/soydoradesu/product_discovery/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) repository.UserRepository {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var u domain.User
	err := r.pool.QueryRow(ctx, `
		SELECT id, email, password_hash, google_id, created_at
		FROM users
		WHERE email = $1
	`, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.GoogleID, &u.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, repository.ErrNotFound
	}
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (domain.User, error) {
	var u domain.User
	err := r.pool.QueryRow(ctx, `
		SELECT id, email, password_hash, google_id, created_at
		FROM users
		WHERE id = $1
	`, id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.GoogleID, &u.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, repository.ErrNotFound
	}
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (r *UserRepo) GetByGoogleID(ctx context.Context, googleID string) (domain.User, error) {
	var u domain.User
	err := r.pool.QueryRow(ctx, `
		SELECT id, email, password_hash, google_id, created_at
		FROM users
		WHERE google_id = $1
	`, googleID).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.GoogleID, &u.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, repository.ErrNotFound
	}
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (r *UserRepo) SetGoogleID(ctx context.Context, userID int64, googleID string) error {
	ct, err := r.pool.Exec(ctx, `
		UPDATE users
		SET google_id = $1
		WHERE id = $2
	`, googleID, userID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *UserRepo) CreateOAuthUser(ctx context.Context, email, googleID string) (int64, error) {
	var id int64
	err := r.pool.QueryRow(ctx, `
		INSERT INTO users(email, google_id)
		VALUES ($1, $2)
		RETURNING id
	`, email, googleID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}