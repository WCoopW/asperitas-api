package post

import (
	"context"

	"reddit/internal/apperrors"
	"reddit/internal/domain/user"

	"go.uber.org/zap"
)

type PostService interface {
	GetPostByID(ctx context.Context, id string) (Post, error)
	GetPosts(ctx context.Context, filter PostFilter, limit int, offset int) ([]Post, error)
	GetPostsByCategory(ctx context.Context, category string) ([]Post, error)
	GetPostsByUser(ctx context.Context, userLogin string) ([]Post, error)
	CreatePost(ctx context.Context, post Post, userID string) (Post, error)
	DeletePost(ctx context.Context, id string, userID string) error
	UpdateVote(ctx context.Context, id string, userID string, value int) (Post, error)
}

type service struct {
	repo        PostRepository
	userService user.UserService
	logger      *zap.SugaredLogger
}

func New(repo PostRepository, userService user.UserService, logger *zap.SugaredLogger) PostService {
	return &service{
		repo:        repo,
		userService: userService,
		logger:      logger,
	}
}

func (s *service) GetPostByID(ctx context.Context, id string) (Post, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetPosts(ctx context.Context, filter PostFilter, limit int, offset int) ([]Post, error) {
	request := PostRequest{
		Limit:      &limit,
		Offset:     &offset,
		PostFilter: filter,
		Populated:  []PostPopulatedFields{AuthorPopulate},
	}

	return s.repo.List(ctx, &request)
}

func (s *service) GetPostsByCategory(ctx context.Context, category string) ([]Post, error) {
	request := PostRequest{
		PostFilter: PostFilter{Category: category},
	}

	return s.repo.List(ctx, &request)
}

func (s *service) GetPostsByUser(ctx context.Context, userLogin string) ([]Post, error) {
	author, err := s.userService.GetUserByUsername(userLogin)
	if err != nil {
		return []Post{}, err
	}
	request := PostRequest{
		Limit:      nil,
		Offset:     nil,
		PostFilter: PostFilter{AuthorID: author.ID},
	}
	return s.repo.List(ctx, &request)
}

func (s *service) CreatePost(ctx context.Context, post Post, userID string) (Post, error) {
	author, err := s.userService.GetUserByID(userID)
	if err != nil {
		return Post{}, err
	}
	post.Author = author

	return s.repo.Create(ctx, &post)
}

func (s *service) DeletePost(ctx context.Context, id string, userID string) error {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if p.Author.ID != userID {
		return apperrors.ErrForbidden
	}
	return s.repo.Delete(ctx, id)
}

func (s *service) UpdateVote(ctx context.Context, id string, userID string, value int) (Post, error) {
	// return s.repo.UpdateVote(ctx, id, userID, value)
	return Post{}, nil
}
