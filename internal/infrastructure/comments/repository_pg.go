package comments

import (
	"context"

	domain "reddit/internal/domain/comments"
	"reddit/internal/infrastructure"

	"github.com/jmoiron/sqlx"
)

type PGRepository struct {
	db     *sqlx.DB
	mapper *CommentMapper
	// queryBuilder *PostQueryBuilder
}

func NewRepo(db *sqlx.DB) domain.CommentRepository {
	return &PGRepository{
		db:     db,
		mapper: &CommentMapper{},
	}
}

func (r *PGRepository) Create(ctx context.Context, comment *domain.Comment, postId string) (domain.Comment, error) {
	schema := r.mapper.EntityToSchema(*comment, postId)
	query := `INSERT INTO comments (post_id, author_id, body) 
	VALUES ($1, $2, $3) 
	RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		schema.PostID,
		schema.AuthorID,
		schema.Body).Scan(&schema.ID, &schema.CreatedAt)
	if err != nil {
		return domain.Comment{}, err
	}
	comment.ID = schema.ID
	comment.CreatedAt = schema.CreatedAt
	return *comment, nil
}

func (r *PGRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM comments WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PGRepository) GetByID(ctx context.Context, id string) (domain.Comment, error) {
	query := "SELECT * FROM comments WHERE id = $1"
	var comment domain.Comment
	err := r.db.Select(&comment, query, id)
	if err != nil {
		return domain.Comment{}, infrastructure.ErrNotFound
	}
	return comment, nil
}
