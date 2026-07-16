package post

import "time"

// PostSchema maps 1:1 to the posts table.
// Use for INSERT/UPDATE and simple SELECT without joins.
type PostSchema struct {
	ID               string    `db:"id"`
	AuthorID         string    `db:"author_id"`
	Type             string    `db:"content_type"`
	Title            string    `db:"title"`
	URL              string    `db:"url_content"`
	Text             string    `db:"text_content"`
	Category         string    `db:"category"`
	Score            int       `db:"score"`
	Views            int       `db:"views"`
	UpvotePercentage int       `db:"upvote_percentage"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

// PostWithAuthorSchema is for reads with JOIN users.
// Flat columns from the join are mapped into domain.Post.Author in the mapper.
//
// Example query:
//
//	SELECT p.id, p.author_id, p.type, p.title, p.url, p.text, p.category,
//	       p.score, p.views, p.upvote_percentage, p.created_at, p.updated_at,
//	       u.username AS author_username
//	FROM posts p
//	JOIN users u ON u.id = p.author_id
//	WHERE p.id = $1
type PostWithAuthorSchema struct {
	PostSchema
	AuthorUsername string `db:"username"`
}
