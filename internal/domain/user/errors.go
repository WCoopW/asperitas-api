package user

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrWrongCredentials = errors.New("wrong credentials")
	ErrUsernameTaken    = errors.New("username already taken")
)
