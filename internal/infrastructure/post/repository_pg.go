package post

import (
	"context"
	"database/sql"

	domain "reddit/internal/domain/post"

	"go.uber.org/zap"
)

type PGRepository struct {
	db           *sql.DB
	queryBuilder *queryBuilder
}

func NewPG(db *sql.DB, logger *zap.SugaredLogger) domain.PostRepository {
	_ = logger
	return &PGRepository{
		db:           db,
		queryBuilder: newQueryBuilder(),
	}
}

func (r *PGRepository) GetByID(ctx context.Context, id string) (domain.Post, error) {
	query, args := r.queryBuilder.build(domain.PostFilter{ID: id}, 0, 0)
	row := r.db.QueryRowContext(context.Background(), query, args...)
	var post domain.Post
	err := row.Scan(
		&post.ID, &post.Title, &post.Text, &post.URL, &post.Category,
		&post.CreatedAt, &post.Score, &post.Views, &post.UpvotePercentage,
		&post.Votes, &post.Comments,
	)
	if err != nil {
		return domain.Post{}, domain.ErrNotFound
	}
	return post, nil
}

func (r *PGRepository) Create(ctx context.Context, post *domain.Post) (domain.Post, error) {
	return domain.Post{}, nil
}

func (r *PGRepository) Update(ctx context.Context, post *domain.Post) (domain.Post, error) {
	return domain.Post{}, nil
}

func (r *PGRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *PGRepository) DeleteComment(ctx context.Context, postID string, commentID string) error {
	return nil
}

func (r *PGRepository) AddComment(ctx context.Context, postID string, comment domain.Comment) (domain.Comment, error) {
	return domain.Comment{}, nil
}

func (r *PGRepository) List(ctx context.Context, filter domain.PostFilter, limit int, offset int) ([]domain.Post, error) {
	return []domain.Post{}, nil
}

func (r *PGRepository) UpdateVote(ctx context.Context, id string, userID string, value int) (domain.Post, error) {
	return domain.Post{}, nil
}
