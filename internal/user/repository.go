package user

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	UserName  string    `json:"username.omitempty"`
	Password  string    `json:"-"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
}

type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	DeleteUser(ctx context.Context, id int) error
}

type dbUserRepository struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) Repository {
	return &dbUserRepository{db: db}
}

func (r *dbUserRepository) CreateUser(ctx context.Context, user *User) error {

	query := "INSERT INTO Users (name, email, password, active) VALUES($1, $2, $3, true)"

	commandTag, err := r.db.Exec(ctx, query, user.Name, user.Email, user.Password)
	if err != nil {
		return fmt.Errorf("could not create user: %w", err)
	}

	fmt.Printf("commandTag: %v\n", commandTag)

	return nil
}

func (r *dbUserRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	query := "SELECT id, email, password, active, userName, created_at FROM Users WHERE email=$1"

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

	return &user, nil
}

// func (r *dbUserRepository) UpdateUser(ctx context.Context, id int) error {
// 	var user User
// 	findQuery := "SELECT id, email, active, userName, FROM Users WHERE id=$1"
//
// 	err := r.db.QueryRow(ctx, findQuery, id).Scan(
// 		&user.ID,
// 		&user.Email,
// 		&user.Password,
// 		&user.UserName,
// 	)
// 	if err != nil {
// 		if err == pgx.ErrNoRows {
// 			return err
// 		}
// 		return err
// 	}
//
// 	return nil
// }

func (r *dbUserRepository) DeleteUser(ctx context.Context, id int) error {
	query := "DELETE FROM Users WHERE id=$1"

	commandTag, err := r.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return fmt.Errorf("could not find user: %w", err)
	}

	return err
}
