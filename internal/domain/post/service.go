package post

import (
	"context"

	"reddit/internal/domain/user"

	"go.uber.org/zap"
)

type PostService interface {
	GetPostByID(ctx context.Context, id string) (Post, error)
	GetPosts(ctx context.Context, filter PostFilter, limit int, offset int) ([]Post, error)
	GetPostsByCategory(ctx context.Context, category string) ([]Post, error)
	GetPostsByUser(ctx context.Context, userLogin string) ([]Post, error)
	AddComment(ctx context.Context, postID string, userID string, comment Comment) (Comment, error)
	CreatePost(ctx context.Context, post Post, userID string) (Post, error)
	DeletePost(ctx context.Context, id string, userID string) error
	DeleteComment(ctx context.Context, postID string, commentID string, userID string) error
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
	return s.repo.List(ctx, filter, limit, offset)
}

func (s *service) GetPostsByCategory(ctx context.Context, category string) ([]Post, error) {
	return s.repo.List(ctx, PostFilter{Category: category}, 0, 0)
}

func (s *service) GetPostsByUser(ctx context.Context, userLogin string) ([]Post, error) {
	author, err := s.userService.GetUserByUsername(userLogin)
	if err != nil {
		return []Post{}, err
	}
	return s.repo.List(ctx, PostFilter{AuthorID: author.ID}, 0, 0)
}

func (s *service) CreatePost(ctx context.Context, post Post, userID string) (Post, error) {
	author, err := s.userService.GetUserByID(userID)
	if err != nil {
		return Post{}, err
	}
	post.Author = author

	return s.repo.Create(ctx, &post)
}

func (s *service) AddComment(ctx context.Context, postID string, userID string, comment Comment) (Comment, error) {
	author, err := s.userService.GetUserByID(userID)
	if err != nil {
		return Comment{}, err
	}
	comment.Author = author
	return s.repo.AddComment(ctx, postID, comment)
}

func (s *service) DeletePost(ctx context.Context, id string, userID string) error {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if p.Author.ID != userID {
		return ErrForbidden
	}
	return s.repo.Delete(ctx, id)
}

func (s *service) DeleteComment(ctx context.Context, id, commentID, userID string) error {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	commentIndex := -1
	for i := range p.Comments {
		if p.Comments[i].ID == commentID {
			if p.Comments[i].Author.ID != userID {
				return ErrForbidden
			}
			commentIndex = i
			break
		}
	}
	if commentIndex == -1 {
		return ErrNotFound
	}

	return s.repo.DeleteComment(ctx, id, commentID)
}

func (s *service) UpdateVote(ctx context.Context, id string, userID string, value int) (Post, error) {
	return s.repo.UpdateVote(ctx, id, userID, value)
}
