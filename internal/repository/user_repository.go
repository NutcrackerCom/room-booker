package repository

import (
	"context"

	"room-booking/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, email, passwordHash string) (*domain.User, error) {
	query := `
		insert into users (id, email, password_hash, role)
		values (gen_random_uuid(), $1, $2, 'user')
		returning id, email, role, created_at
	`

	var user domain.User
	err := r.db.QueryRow(ctx, query, email, passwordHash).Scan(
		&user.ID,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, string, error) {
	query := `
		select id, email, role, created_at, password_hash
		from users
		where email = $1
	`

	var user domain.User
	var passwordHash string

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&passwordHash,
	)
	if err != nil {
		return nil, "", err
	}

	return &user, passwordHash, nil
}
