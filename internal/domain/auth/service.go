package auth

import (
	"reddit/internal/domain/user"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(username, password string) (token string, err error)
	Register(username, password string) (token string, err error)
	Logout(token string) error
}

type service struct {
	users  user.UserService
	tokens TokenIssuer
	logger *zap.SugaredLogger
}

func New(users user.UserService, tokens TokenIssuer, logger *zap.SugaredLogger) AuthService {
	return &service{users: users, tokens: tokens, logger: logger}
}

func (s *service) Login(username, password string) (string, error) {
	u, err := s.users.GetUserByUsername(username)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", user.ErrWrongCredentials
	}

	return s.tokens.GenerateToken(u)
}

func (s *service) Register(username, password string) (string, error) {
	if _, err := s.users.GetUserByUsername(username); err == nil {
		return "", user.ErrUsernameTaken
	} else if err != nil && err != user.ErrNotFound {
		return "", err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	u, err := s.users.CreateUser(user.User{
		Username:     username,
		PasswordHash: string(hash),
	})
	if err != nil {
		return "", err
	}
	return s.tokens.GenerateToken(u)
}

func (s *service) Logout(token string) error {
	_ = token
	return nil
}
