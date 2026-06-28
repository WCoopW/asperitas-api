package post

import (
	"context"
	"slices"

	domain "reddit/internal/domain/post"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type MemRepository struct {
	logger *zap.SugaredLogger
	items  []*domain.Post
}

func NewMem(logger *zap.SugaredLogger) domain.PostRepository {
	return &MemRepository{
		logger: logger,
		items:  []*domain.Post{},
	}
}

func (r *MemRepository) GetByID(ctx context.Context, id string) (domain.Post, error) {
	p, err := r.findPostByID(ctx, id)
	if err != nil {
		return domain.Post{}, err
	}
	return *p, nil
}

func (r *MemRepository) Create(ctx context.Context, post *domain.Post) (domain.Post, error) {
	post.ID = uuid.New().String()
	r.items = append(r.items, post)
	return *post, nil
}

func (r *MemRepository) Update(ctx context.Context, post *domain.Post) (domain.Post, error) {
	return domain.Post{}, nil
}

func (r *MemRepository) Delete(ctx context.Context, id string) error {
	r.items = slices.DeleteFunc(r.items, func(p *domain.Post) bool {
		return p.ID == id
	})
	return nil
}

func (r *MemRepository) DeleteComment(ctx context.Context, postID string, commentID string) error {
	p, err := r.findPostByID(ctx, postID)
	if err != nil {
		return err
	}
	p.Comments = slices.DeleteFunc(p.Comments, func(c domain.Comment) bool {
		return c.ID == commentID
	})
	return nil
}

func (r *MemRepository) AddComment(ctx context.Context, postID string, comment domain.Comment) (domain.Comment, error) {
	p, err := r.findPostByID(ctx, postID)
	if err != nil {
		return domain.Comment{}, err
	}
	comment.ID = uuid.New().String()
	p.Comments = append(p.Comments, comment)
	return comment, nil
}

func (r *MemRepository) List(ctx context.Context, filter domain.PostFilter, limit int, offset int) ([]domain.Post, error) {
	out := make([]domain.Post, 0)
	for _, p := range r.items {
		if filter.ID != "" && p.ID != filter.ID {
			continue
		}
		if filter.Category != "" && p.Category != filter.Category {
			continue
		}
		if filter.AuthorID != "" && p.Author.ID != filter.AuthorID {
			continue
		}
		out = append(out, *p)
	}
	if offset >= len(out) {
		return []domain.Post{}, nil
	}
	out = out[offset:]
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

func (r *MemRepository) UpdateVote(ctx context.Context, id string, userID string, value int) (domain.Post, error) {
	p, err := r.findPostByID(ctx, id)
	if err != nil {
		return domain.Post{}, err
	}

	voteIndex := -1
	for i := range p.Votes {
		if p.Votes[i].UserID == userID {
			voteIndex = i
			break
		}
	}

	if voteIndex == -1 {
		if value != 0 {
			p.Votes = append(p.Votes, domain.Vote{
				UserID: userID,
				Vote:   value,
			})
			p.Score += value
		}
	} else {
		oldValue := p.Votes[voteIndex].Vote

		switch value {
		case 1, -1:
			p.Votes[voteIndex].Vote = value
			p.Score += value - oldValue
		case 0:
			p.Votes = slices.Delete(p.Votes, voteIndex, voteIndex+1)
			p.Score -= oldValue
		default:
			return domain.Post{}, domain.ErrInvalidData
		}
	}

	return *p, nil
}

func (r *MemRepository) findPostByID(ctx context.Context, id string) (*domain.Post, error) {
	for _, item := range r.items {
		if item.ID == id {
			return item, nil
		}
	}
	return nil, domain.ErrNotFound
}
