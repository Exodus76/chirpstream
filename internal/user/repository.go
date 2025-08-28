package user

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID        int64
	Email     string
	UserName  string
	Password  string
	Active    bool
	CreatedAt time.Time
}

type Repository interface {
	CreateUser(ctx context.Context, user *User) error
}

type dbUserRepository struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) Repository {
	return &dbUserRepository{db: db}
}

func (r *dbUserRepository) CreateUser(ctx context.Context, user *User) error {

	query := "INSERT INTO Users (email, password, active) VALUES($1, $2)"

	err := r.db.QueryRow(ctx, query, user.Email, user.Password, 1).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("could not create user: %w", err)
	}

	return nil
}

func (r *dbUserRepository) GetUserByEmail(ctx context.Context, user *User) (*User, error) {

	query := "SELECT id, email, password, active, user_name, created_at FROM users WHERE email=$1"

	err := r.db.QueryRow(ctx, query, user.Email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Active,
		&user.UserName,
		&user.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("could not find user: %w", err)
	}

	return nil, nil
}
