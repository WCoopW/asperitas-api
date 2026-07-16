package post

import (
	"reddit/internal/apperrors"
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
		return &apperrors.ValidationError{
			Field:   "type",
			Message: "invalid post type",
		}
	}

	if dto.Type == "link" && dto.URL == "" {
		return &apperrors.ValidationError{
			Field:   "url",
			Message: "url is required for link posts",
		}
	}

	if dto.Type == "text" && dto.Text == "" {
		return &apperrors.ValidationError{
			Field:   "text",
			Message: "text is required for text posts",
		}
	}

	if dto.Title == "" {
		return &apperrors.ValidationError{
			Field:   "title",
			Message: "title is required",
		}
	}
	if dto.Category == "" {
		return &apperrors.ValidationError{
			Field:   "category",
			Message: "category is required",
		}
	}
	return nil
}

type CreateCommentDTO struct {
	Body string `json:"body"`
}
