package post

import (
	"time"

	comment "reddit/internal/domain/comments"
	domain "reddit/internal/domain/post"
	"reddit/internal/domain/user"
)

type AuthorResponseDTO struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func newAuthorResponseDTO(u user.User) AuthorResponseDTO {
	return AuthorResponseDTO{
		ID:       u.ID,
		Username: u.Username,
	}
}

type VoteResponseDTO struct {
	UserID string `json:"user"`
	Vote   int    `json:"vote"`
}

type CommentResponseDTO struct {
	ID        string            `json:"id"`
	CreatedAt time.Time         `json:"created_at"`
	Author    AuthorResponseDTO `json:"author"`
	Body      string            `json:"body"`
}

func newCommentResponseDTO(comment comment.Comment) CommentResponseDTO {
	return CommentResponseDTO{
		ID:        comment.ID,
		CreatedAt: comment.CreatedAt,
		Body:      comment.Body,
		Author:    newAuthorResponseDTO(comment.Author),
	}
}

type PostResponseDTO struct {
	ID               string               `json:"id"`
	Score            int                  `json:"score"`
	Views            int                  `json:"views"`
	Type             string               `json:"type"`
	Author           AuthorResponseDTO    `json:"author"`
	Title            string               `json:"title"`
	Text             string               `json:"text,omitempty"`
	URL              string               `json:"url,omitempty"`
	Votes            []VoteResponseDTO    `json:"votes"`
	UpvotePercentage int                  `json:"upvote_percentage"`
	Category         string               `json:"category"`
	Comments         []CommentResponseDTO `json:"comments"`
	CreatedAt        time.Time            `json:"created_at"`
}

func NewPostResponseDTO(post domain.Post) PostResponseDTO {
	votes := make([]VoteResponseDTO, len(post.Votes))
	for i, v := range post.Votes {
		votes[i] = VoteResponseDTO{UserID: v.UserID, Vote: v.Vote}
	}

	comments := make([]CommentResponseDTO, len(post.Comments))
	for i, c := range post.Comments {
		comments[i] = newCommentResponseDTO(c)
	}

	return PostResponseDTO{
		ID:               post.ID,
		Score:            post.Score,
		Views:            post.Views,
		Type:             string(post.Type),
		Author:           newAuthorResponseDTO(post.Author),
		Title:            post.Title,
		Text:             post.Text,
		URL:              post.URL,
		Votes:            votes,
		UpvotePercentage: post.UpvotePercentage,
		Category:         post.Category,
		Comments:         comments,
		CreatedAt:        post.CreatedAt,
	}
}
