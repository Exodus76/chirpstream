package userrelationship

import "context"

type Service interface {
	FollowUser(ctx context.Context, userId int, followingId int) error
	UnfollowUser(ctx context.Context, userId int, followingId int) error
	GetFollowing(ctx context.Context, userId int) ([]int, error)
	GetFollowers(ctx context.Context, userId int) ([]int, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// FollowUser implements [Service].
func (s *service) FollowUser(ctx context.Context, userId int, followingId int) error {
	return s.repo.FollowUser(ctx, userId, followingId)
}

// GetFollowers implements [Service].
func (s *service) GetFollowers(ctx context.Context, userId int) ([]int, error) {
	return s.repo.GetFollowers(ctx, userId)
}

// GetFollowing implements [Service].
func (s *service) GetFollowing(ctx context.Context, userId int) ([]int, error) {
	return s.repo.GetFollowing(ctx, userId)
}

// UnfollowUser implements [Service].
func (s *service) UnfollowUser(ctx context.Context, userId int, followingId int) error {
	return s.repo.UnfollowUser(ctx, userId, followingId)
}
