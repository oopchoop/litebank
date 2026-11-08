package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, user NewUser) (int64, error) {
	var id int64

	err := r.pool.QueryRow(ctx,
		`INSERT INTO users(phone_number, password_hash, first_name, last_name, patronymic)
		 VALUES ($1, $2, $3, $4, $5) 
		 RETURNING id`,
		user.PhoneNumber, user.PasswordHash, user.FirstName, user.LastName, user.Patronymic,
	).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "users_phone_number_key" {
				return 0, ErrPhoneAlreadyUsed
			}
		}

		return 0, fmt.Errorf("Создание пользователя: %w", err)
	}

	return id, nil
}

func (r *Repository) GetByPhone(ctx context.Context, phone string) (*User, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, phone_number, password_hash, first_name, last_name, patronymic, created_at 
		 FROM users WHERE phone_number = $1 
		 LIMIT 1`,
		phone)
	if err != nil {
		return nil, fmt.Errorf("query user: %w", err)
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotRegister
		}

		return nil, fmt.Errorf("get user by phone: %w", err)
	}

	return &user, nil
}
