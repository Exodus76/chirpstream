package user

import (
	"chirpstream/pkg/apperror"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	Password   string    `json:"-"`
	Active     bool      `json:"active"`
	Created_at time.Time `json:"created_at"`
}

type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserById(ctx context.Context, id int) (*User, error)
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

	_, err := r.db.Exec(ctx, query, user.Name, user.Email, user.Password)
	if err != nil {
		return fmt.Errorf("User.CreateUser: could not insert user: %w", err)
	}

	return nil
}

func (r *dbUserRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	query := "SELECT id, name, email, password, created_at FROM Users WHERE email=$1 AND active=TRUE"

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Created_at,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, &apperror.RepoError{
				Op:  "User.GetUserByEmail",
				Err: ErrNotFound,
				ID:  email,
			}
		}
		return nil, fmt.Errorf("User.GetUserbyEmail: failed to execute query %w", err)
	}

	return &user, nil
}

func (r *dbUserRepository) GetUserById(ctx context.Context, id int) (*User, error) {
	var user User
	query := "SELECT id, name, username, email FROM Users WHERE id=$1 AND active=TRUE"

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			// return nil, fmt.Errorf("GetUserById: no user with this id: %d %w", id, err)
			return nil, &apperror.RepoError{
				Op:  "user.GetUserById",
				Err: ErrNotFound,
				ID:  strconv.Itoa(id),
			}
		}
		return nil, fmt.Errorf("User.UserGetUserById: failed to execute query %w", err)
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
		// return fmt.Errorf("DeleteUser: could not delete user with id:%d %w", id, err)
		return &apperror.RepoError{
			Op:  "User.DeleteUser",
			Err: err,
			ID:  strconv.Itoa(id),
		}
	}

	return nil
}
