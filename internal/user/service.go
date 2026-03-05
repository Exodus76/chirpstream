package user

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	CreateUser(ctx context.Context, name, email, password string) error
	VerifyUser(ctx context.Context, email string, password string) (*User, error)
	GetUserById(ctx context.Context, id int) (*User, error)
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
		return fmt.Errorf("failed to encrypt password %w", err)
	}

	newUser := &User{
		Email:    email,
		Password: string(hashedPassword),
		Active:   true,
	}

	err = s.repo.CreateUser(ctx, newUser)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) VerifyUser(ctx context.Context, email string, password string) (*User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("Error verifying password %w", err)
	}

	return user, nil
}

func (s *service) GetUserById(ctx context.Context, id int) (*User, error) {
	user, err := s.repo.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) DeleteUser(ctx context.Context, id int) error {
	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
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
