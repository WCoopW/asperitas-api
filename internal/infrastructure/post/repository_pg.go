package post

import (
	"context"
	"fmt"

	"reddit/internal/apperrors"
	domain "reddit/internal/domain/post"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type PGRepository struct {
	db           *sqlx.DB
	mapper       *PostMapper
	queryBuilder *PostQueryBuilder
}

func NewPG(db *sqlx.DB, logger *zap.SugaredLogger) domain.PostRepository {
	_ = logger
	return &PGRepository{
		db:           db,
		mapper:       &PostMapper{},
		queryBuilder: &PostQueryBuilder{},
	}
}

func (r *PGRepository) GetByID(ctx context.Context, id string) (domain.Post, error) {
	query := "SELECT * FROM posts WHERE id = $1"
	var post domain.Post
	err := r.db.Select(&post, query, id)
	if err != nil {
		return domain.Post{}, apperrors.ErrNotFound
	}
	return post, nil
}

func (r *PGRepository) Create(ctx context.Context, post *domain.Post) (domain.Post, error) {
	schema := r.mapper.EntityToSchema(*post)
	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO posts (author_id, content_type, title, text_content, url_content, category)
         VALUES ($1, $2, $3, $4, $5, $6)
         RETURNING id, created_at`,
		schema.AuthorID,
		schema.Type,
		schema.Title,
		schema.Text,
		schema.URL,
		schema.Category,
	).Scan(&schema.ID, &schema.CreatedAt)
	if err != nil {
		return domain.Post{}, err
	}
	post.ID = schema.ID
	post.CreatedAt = schema.CreatedAt
	return *post, nil
}

func (r *PGRepository) Update(ctx context.Context, post *domain.Post) error {
	schema := r.mapper.EntityToSchema(*post)
	query := `UPDATE posts SET 
	author_id: author_id, type = :type, title = :title, url = :url,
	text = :text, category = :category,	score = :score,	
	views = :views,	upvote_percentage = :upvote_percentage,
	created_at = :created_at, updated_at = :updated_at
	WHERE id = :id`
	_, err := r.db.NamedExecContext(ctx, query, schema)
	return err
}

func (r *PGRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM posts WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PGRepository) List(ctx context.Context, request *domain.PostRequest) ([]domain.Post, error) {
	query, args := r.queryBuilder.BuildQuery(request)
	var posts []PostWithAuthorSchema
	err := r.db.Select(&posts, query, args...)
	if err != nil {
		return nil, fmt.Errorf("fetch posts failed: %w", err)
	}
	results := make([]domain.Post, 0, len(posts))
	for _, v := range posts {
		results = append(results, r.mapper.SchemaWithAuthorToEntity(v))
	}

	return results, nil
}
