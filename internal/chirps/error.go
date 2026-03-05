package chirps

import "errors"

var (
	ErrChirpNotFound = errors.New("chirp not found")
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
	ErrUserNotFound  = errors.New("user not found")
)
