package comments

import (
	"context"

	"reddit/internal/apperrors"
	user "reddit/internal/domain/user"
)

type CommentService interface {
	Create(ctx context.Context, body, postID, userID string) (Comment, error)
	Delete(ctx context.Context, id, userID string) error
}

type service struct {
	repo        CommentRepository
	userService user.UserService
}

func New(repo CommentRepository, userService user.UserService) CommentService {
	return &service{
		repo:        repo,
		userService: userService,
	}
}

func (s *service) Create(ctx context.Context, body, postID, userID string) (Comment, error) {
	author, err := s.userService.GetUserByID(userID)
	if err != nil {
		return Comment{}, err
	}
	comment := Comment{
		Author: author,
		Body:   body,
	}
	return s.repo.Create(ctx, &comment, postID)
}

func (s *service) Delete(ctx context.Context, id, userID string) error {
	author, err := s.userService.GetUserByID(userID)
	if err != nil {
		return err
	}
	comment, err := s.getByID(ctx, id)
	if err != nil {
		return err
	}
	if author.ID != comment.Author.ID {
		return apperrors.ErrForbidden
	}

	return s.repo.Delete(ctx, id)
}

func (s *service) getByID(ctx context.Context, id string) (Comment, error) {
	return s.repo.GetByID(ctx, id)
}
