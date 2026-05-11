package post

import "errors"

type CreatePostDTO struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Text     string `json:"text"`
	URL      string `json:"url"`
	Category string `json:"category"`
}

func (dto *CreatePostDTO) Validate() error {
	if dto.Type != "link" && dto.Type != "text" {
		return errors.New("invalid post type")
	}

	if dto.Type == "link" && dto.URL == "" {
		return errors.New("url is required for link posts")
	}

	if dto.Type == "text" && dto.Text == "" {
		return errors.New("text is required for text posts")
	}

	if dto.Title == "" {
		return errors.New("title is required")
	}
	if dto.Category == "" {
		return errors.New("category is required")
	}
	return nil
}

type CreateCommentDTO struct {
	Body string `json:"body"`
}
