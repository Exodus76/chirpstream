package user

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	CreateUser(ctx context.Context, name, email, password string) error
	VerifyUser(ctx context.Context, email string, password string) (*User, error)
	DeleteUser(ctx context.Context, id int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateUser(ctx context.Context, name, email, password string) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return fmt.Errorf("Error encrypting password %w", err)
	}

	newUser := &User{
		Email:    email,
		Password: string(hashedPassword),
	}

	err = s.repo.CreateUser(ctx, newUser)
	if err != nil {
		return fmt.Errorf("Error creating user %w", err)
	}

	return nil
}

func (s *service) VerifyUser(ctx context.Context, email string, password string) (*User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("Error fetching user %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("No user with this email %w", user.Email)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("Error verifying password %w", err)
	}

	return user, nil
}

func (s *service) DeleteUser(ctx context.Context, id int) error {
	panic("not implemented") // TODO: Implement
}

//backup
// func (s *service) CreateUser(ctx context.Context, name, email, password string) error {
//
// 	tx, err := s.pool.Begin(ctx)
// 	if err != nil {
// 		return fmt.Errorf("could not begin transaction: %w", err)
// 	}
//
// 	defer tx.Rollback(ctx)
//
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
// 	if err != nil {
// 		return fmt.Errorf("Error encrypting password %w", err)
// 	}
//
// 	newUser := &User{
// 		Email:    email,
// 		Password: string(hashedPassword),
// 	}
//
// 	err = s.repo.CreateUser(ctx, newUser)
// 	if err != nil {
// 		return fmt.Errorf("Error creating user %w", err)
// 	}
//
// 	if err := tx.Commit(ctx); err != nil {
// 		return fmt.Errorf("could not commit transaction: %w", err)
// 	}
//
// 	return nil
// }
