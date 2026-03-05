package apperror

import "fmt"

type RepoError struct {
	Op  string // operation
	Err error  // underlying error
	ID  string // optional context
}

func (e *RepoError) Error() string {
	return fmt.Sprintf("op=%s id=%s: %v", e.Op, e.ID, e.Err)
}

func (e *RepoError) Unwrap() error {
	return e.Err
}
