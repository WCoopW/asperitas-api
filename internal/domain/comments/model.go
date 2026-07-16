package comments

import (
	"context"
	"time"

	user "reddit/internal/domain/user"
)

type Comment struct {
	ID        string
	Body      string
	Author    user.User
	CreatedAt time.Time
}

type CommentRepository interface {
	Create(ctx context.Context, comment *Comment, postId string) (Comment, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (Comment, error)
}
