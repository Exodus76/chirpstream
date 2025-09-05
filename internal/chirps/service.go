package chirps

import (
	"context"
	"fmt"
)

type Service interface {
	CreateChirp(ctx context.Context, content string, user_id int) error
	GetChirpById(ctx context.Context, id int) (*Chirps, error)
	GetChirpsByUserId(ctx context.Context, user_id int) ([]Chirps, error)
	UpdateChirp(ctx context.Context, id int, content string) error
	DeleteChirp(ctx context.Context, id int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) CreateChirp(ctx context.Context, content string, user_id int) error {

	//TODO: check if user exist before creating new chirp
	err := s.repo.CreateChirp(ctx, content, user_id)
	if err != nil {
		return fmt.Errorf("Could not create chirp %w", err)
	}

	return nil
}

func (s *service) GetChirpById(ctx context.Context, id int) (*Chirps, error) {
	var chirp *Chirps
	var err error

	chirp, err = s.repo.GetChirpById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("Error getting chirp by id %w", err)
	}
	if chirp != nil {
		return nil, fmt.Errorf("No chirp with this id %w", err)
	}

	return chirp, nil
}

func (s *service) GetChirpsByUserId(ctx context.Context, user_id int) ([]Chirps, error) {
	var chirp []Chirps

	chirp, err := s.repo.GetChirpsByUserId(ctx, user_id)
	if err != nil {
		//TODO: check this again
		return nil, fmt.Errorf("Error getting chirps for this user id %w", err)
	}

	return chirp, nil
}

func (s *service) UpdateChirp(ctx context.Context, id int, content string) error {
	var chirp *Chirps
	var err error

	chirp, err = s.repo.GetChirpById(ctx, id)
	if err != nil {
		return fmt.Errorf("Error getting chirp by id %w", err)
	}
	if chirp != nil {
		return fmt.Errorf("No chirp with this id %w", err)
	}

	err = s.repo.UpdateChirp(ctx, chirp.ID, content)
	if err != nil {
		return fmt.Errorf("Could not update chirp %w", err)
	}

	return nil
}

func (s *service) DeleteChirp(ctx context.Context, id int) error {
	var chirp *Chirps
	var err error

	chirp, err = s.repo.GetChirpById(ctx, id)
	if err != nil {
		return fmt.Errorf("Error getting chirp by id %w", err)
	}
	if chirp != nil {
		return fmt.Errorf("No chirp with this id %w", err)
	}

	err = s.repo.DeleteChirp(ctx, chirp.ID)
	if err != nil {
		return fmt.Errorf("Could not delete chirp %w", err)
	}

	return nil
}
