package post

import (
	"context"
	"time"

	"reddit/internal/domain/user"
)

type PostType string

const (
	Link PostType = "link"
	Text PostType = "text"
)

type Post struct {
	ID               string
	Score            int
	Views            int
	Type             PostType
	Author           user.User
	Title            string
	Text             string
	URL              string
	Votes            []Vote
	UpvotePercentage int
	Category         string
	Comments         []Comment
	CreatedAt        time.Time
}

type Vote struct {
	UserID string
	Vote   int
}

type Comment struct {
	ID        string
	CreatedAt time.Time
	Author    user.User
	Body      string
}

type PostFilter struct {
	ID       string
	Category string
	AuthorID string
}

type PostRepository interface {
	Create(ctx context.Context, post *Post) (Post, error)
	GetByID(ctx context.Context, id string) (Post, error)
	Update(ctx context.Context, post *Post) (Post, error)
	Delete(ctx context.Context, id string) error
	DeleteComment(ctx context.Context, postID string, commentID string) error
	AddComment(ctx context.Context, postID string, comment Comment) (Comment, error)
	List(ctx context.Context, filter PostFilter, limit int, offset int) ([]Post, error)
	UpdateVote(ctx context.Context, id string, userID string, value int) (Post, error)
}

type PostRequest struct {
	Filter PostFilter
	Limit  int
	Offset int
}
