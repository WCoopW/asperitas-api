package post

import (
	domain "reddit/internal/domain/post"
)

type CreatePostDTO struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Text     string `json:"text"`
	URL      string `json:"url"`
	Category string `json:"category"`
}

func (dto *CreatePostDTO) Validate() error {
	if dto.Type != "link" && dto.Type != "text" {
		return &domain.ValidationError{
			Field:   "type",
			Message: "invalid post type",
		}
	}

	if dto.Type == "link" && dto.URL == "" {
		return &domain.ValidationError{
			Field:   "url",
			Message: "url is required for link posts",
		}
	}

	if dto.Type == "text" && dto.Text == "" {
		return &domain.ValidationError{
			Field:   "text",
			Message: "text is required for text posts",
		}
	}

	if dto.Title == "" {
		return &domain.ValidationError{
			Field:   "title",
			Message: "title is required",
		}
	}
	if dto.Category == "" {
		return &domain.ValidationError{
			Field:   "category",
			Message: "category is required",
		}
	}
	return nil
}

type CreateCommentDTO struct {
	Body string `json:"body"`
}
