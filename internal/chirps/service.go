package chirps

import (
	"context"
	"fmt"

	"github.com/gocql/gocql"
)

type Service interface {
	CreateChirp(ctx context.Context, content string, userId int) error
	GetChirpById(ctx context.Context, userId int, chirpId gocql.UUID) (*Chirp, error)
	GetChirpsByUserId(ctx context.Context, userId int, pageState []byte, limit int) ([]Chirp, []byte, error)
	UpdateChirp(ctx context.Context, userId int, chirpId gocql.UUID, content string) error
	DeleteChirp(ctx context.Context, userId int, chirpId gocql.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) CreateChirp(ctx context.Context, content string, userId int) error {

	//TODO: check if user exist before creating new chirp
	err := s.repo.CreateChirp(ctx, content, userId)
	if err != nil {
		return fmt.Errorf("failed to create chirp %w", err)
	}

	return nil
}

func (s *service) GetChirpById(ctx context.Context, userId int, chirpId gocql.UUID) (*Chirp, error) {
	var chirp *Chirp

	chirp, err := s.repo.GetChirpById(ctx, userId, chirpId)
	if err != nil {
		return nil, fmt.Errorf("failed to get chirp with id: %s %w", chirpId, err)
	}

	return chirp, nil
}

func (s *service) GetChirpsByUserId(ctx context.Context, userId int, pageState []byte, limit int) ([]Chirp, []byte, error) {
	var chirp []Chirp
	var nextPageState []byte

	chirp, nextPageState, err := s.repo.GetChirpsByUserId(ctx, userId, pageState, limit)
	if err != nil {
		return nil, nil, fmt.Errorf("failed getting chirps for user id: %d %w", userId, err)
	}

	return chirp, nextPageState, nil
}

func (s *service) UpdateChirp(ctx context.Context, userId int, chirpId gocql.UUID, content string) error {
	//check if the chirp exists
	chirp, err := s.repo.GetChirpById(ctx, userId, chirpId)
	if err != nil {
		return fmt.Errorf("failed to get chirp with id: %s %w", chirpId, err)
	}

	err = s.repo.UpdateChirp(ctx, chirp.UserId, chirp.ChirpId, content)
	if err != nil {
		return fmt.Errorf("failed to update chirp with id: %s %w", chirpId, err)
	}

	return nil
}

func (s *service) DeleteChirp(ctx context.Context, userId int, chirpId gocql.UUID) error {
	chirp, err := s.repo.GetChirpById(ctx, userId, chirpId)
	if err != nil {
		return fmt.Errorf("failed to get chirp with id: %s %w", chirpId, err)
	}

	err = s.repo.DeleteChirp(ctx, chirp.UserId, chirp.ChirpId)
	if err != nil {
		return fmt.Errorf("Could not delete chirp with id: %s %w", chirpId, err)
	}

	return nil
}
