package auth

import "reddit/internal/domain/user"

type TokenIssuer interface {
	GenerateToken(user user.User) (string, error)
}

type TokenValidator interface {
	ValidateToken(raw string) (userID string, err error)
}
