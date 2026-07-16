package comments

import (
	domain "reddit/internal/domain/comments"
	domainuser "reddit/internal/domain/user"
)

type CommentMapper struct{}

func (m *CommentMapper) EntityToSchema(entity domain.Comment, postID string) CommentSchema {
	return CommentSchema{
		ID:        entity.ID,
		PostID:    postID,
		AuthorID:  entity.Author.ID,
		Body:      entity.Body,
		CreatedAt: entity.CreatedAt,
	}
}

// SchemaToEntity maps a row from comments only (no join).
func (m *CommentMapper) SchemaToEntity(schema *CommentSchema) domain.Comment {
	return domain.Comment{
		ID:        schema.ID,
		Body:      schema.Body,
		CreatedAt: schema.CreatedAt,
		Author: domainuser.User{
			ID: schema.AuthorID,
		},
	}
}

// SchemaWithAuthorToEntity maps flat JOIN columns into nested domain.Comment.Author.
func (m *CommentMapper) SchemaWithAuthorToEntity(schema *CommentWithAuthorSchema) domain.Comment {
	comment := m.SchemaToEntity(&schema.CommentSchema)
	comment.Author = domainuser.User{
		ID:       schema.AuthorID,
		Username: schema.AuthorUsername,
	}
	return comment
}
