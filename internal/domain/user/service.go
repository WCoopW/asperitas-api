package user

import "go.uber.org/zap"

type UserRepository interface {
	GetUserByUsername(username string) (User, error)
	GetUserByID(id string) (User, error)
	CreateUser(u User) (User, error)
}

type UserService interface {
	GetUserByUsername(username string) (User, error)
	GetUserByID(id string) (User, error)
	CreateUser(u User) (User, error)
}

type service struct {
	repo   UserRepository
	logger *zap.SugaredLogger
}

func New(repo UserRepository, logger *zap.SugaredLogger) UserService {
	return &service{repo: repo, logger: logger}
}

func (s *service) GetUserByUsername(username string) (User, error) {
	return s.repo.GetUserByUsername(username)
}

func (s *service) GetUserByID(id string) (User, error) {
	return s.repo.GetUserByID(id)
}

func (s *service) CreateUser(u User) (User, error) {
	return s.repo.CreateUser(u)
}
