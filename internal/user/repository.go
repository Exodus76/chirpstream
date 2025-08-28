package user

import (
	"database/sql"
	"time"
)

type User struct {
	ID         int64
	Email      string
	Password   string
	Active     bool
	Created_at time.Time
}

type Repository interface {
	CreateUser(user *User) error
}

type dbUserRepository struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) Repository {
	return &dbUserRepository{db: db}
}

func (r *dbUserRepository) CreateUser(user *User) error {
	return nil
}
