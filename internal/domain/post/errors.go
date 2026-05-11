package post

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidData = errors.New("invalid data")
	ErrNotFound    = errors.New("not found")
	ErrForbidden   = errors.New("forbidden")
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s: %s", e.Field, e.Message)
}
