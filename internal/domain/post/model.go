package post

import (
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
