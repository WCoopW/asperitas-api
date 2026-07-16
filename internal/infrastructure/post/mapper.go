package post

import (
	domain "reddit/internal/domain/post"
	domainuser "reddit/internal/domain/user"
)

type PostMapper struct{}

func (m *PostMapper) EntityToSchema(entity domain.Post) PostSchema {
	return PostSchema{
		ID:               entity.ID,
		AuthorID:         entity.Author.ID,
		Type:             string(entity.Type),
		Title:            entity.Title,
		URL:              entity.URL,
		Text:             entity.Text,
		Category:         entity.Category,
		Score:            entity.Score,
		Views:            entity.Views,
		UpvotePercentage: entity.UpvotePercentage,
		CreatedAt:        entity.CreatedAt,
	}
}

// SchemaToEntity maps a row from posts only (no join).
func (m *PostMapper) SchemaToEntity(schema PostSchema) domain.Post {
	return domain.Post{
		ID:               schema.ID,
		Type:             domain.PostType(schema.Type),
		Title:            schema.Title,
		Text:             schema.Text,
		URL:              schema.URL,
		Category:         schema.Category,
		Score:            schema.Score,
		Views:            schema.Views,
		UpvotePercentage: schema.UpvotePercentage,
		CreatedAt:        schema.CreatedAt,
		Author: domainuser.User{
			ID: schema.AuthorID,
		},
	}
}

// SchemaWithAuthorToEntity maps a flat JOIN row into nested domain.Post.
// author_id -> Author.ID, author_username -> Author.Username
func (m *PostMapper) SchemaWithAuthorToEntity(schema PostWithAuthorSchema) domain.Post {
	post := m.SchemaToEntity(schema.PostSchema)
	post.Author = domainuser.User{
		ID:       schema.AuthorID,
		Username: schema.AuthorUsername,
	}
	return post
}
