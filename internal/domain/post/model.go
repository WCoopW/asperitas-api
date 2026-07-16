package post

import (
	"context"
	"time"

	"reddit/internal/domain/comments"
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
	Comments         []comments.Comment
	CreatedAt        time.Time
}

type PostPopulatedFields string

const (
	AuthorPopulate   PostPopulatedFields = "author"
	VotesPopulate    PostPopulatedFields = "vote"
	CommentsPopulate PostPopulatedFields = "comment"
)

type Vote struct {
	UserID string
	Vote   int
}

type PostFilter struct {
	ID       string
	Category string
	AuthorID string
}

type PostRepository interface {
	Create(ctx context.Context, post *Post) (Post, error)
	GetByID(ctx context.Context, id string) (Post, error)
	Update(ctx context.Context, post *Post) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, request *PostRequest) ([]Post, error)
}

type PostRequest struct {
	PostFilter
	Limit     *int
	Offset    *int
	Populated []PostPopulatedFields
}
