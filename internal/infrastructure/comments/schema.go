package comments

import "time"

// CommentSchema maps 1:1 to the comments table.
type CommentSchema struct {
	ID        string    `db:"id"`
	PostID    string    `db:"post_id"`
	AuthorID  string    `db:"author_id"`
	Body      string    `db:"body"`
	CreatedAt time.Time `db:"created_at"`
}

// CommentWithAuthorSchema is for reads with JOIN users.
//
// Example query:
//
//	SELECT c.id, c.post_id, c.author_id, c.body, c.created_at,
//	       u.username AS author_username
//	FROM comments c
//	JOIN users u ON u.id = c.author_id
//	WHERE c.post_id = $1
type CommentWithAuthorSchema struct {
	CommentSchema
	AuthorUsername string `db:"author_username"`
}
