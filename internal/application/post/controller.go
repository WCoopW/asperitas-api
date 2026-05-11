package post

import (
	"encoding/json"
	"errors"
	"net/http"

	domain "reddit/internal/domain/post"
	"reddit/pkg/helpers"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type PostController struct {
	postService domain.PostService
	logger      *zap.SugaredLogger
}

func New(postService domain.PostService, logger *zap.SugaredLogger) *PostController {
	return &PostController{postService: postService, logger: logger}
}

func (c *PostController) GetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := c.postService.GetPosts()
	if err != nil {
		writeAppError(w, err)
		return
	}

	postsDTO := make([]PostResponseDTO, len(posts))
	for i, post := range posts {
		postsDTO[i] = NewPostResponseDTO(post)
	}

	helpers.WriteJSON(w, http.StatusOK, postsDTO)
}

func (c *PostController) GetPostsByCategory(w http.ResponseWriter, r *http.Request) {
	category := mux.Vars(r)["CATEGORY_NAME"]
	posts, err := c.postService.GetPostsByCategory(category)
	if err != nil {
		writeAppError(w, err)
		return
	}
	postsDTO := make([]PostResponseDTO, len(posts))
	for i, post := range posts {
		postsDTO[i] = NewPostResponseDTO(post)
	}
	helpers.WriteJSON(w, http.StatusOK, postsDTO)
}

func (c *PostController) GetPostByID(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["POST_ID"]
	post, err := c.postService.GetPostByID(postID)
	if err != nil {
		writeAppError(w, err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, NewPostResponseDTO(post))
}

func (c *PostController) GetPostByUser(w http.ResponseWriter, r *http.Request) {
	userLogin := mux.Vars(r)["USER_LOGIN"]
	posts, err := c.postService.GetPostsByUser(userLogin)
	if err != nil {
		writeAppError(w, err)
		return
	}
	postsDTO := make([]PostResponseDTO, len(posts))
	for i, post := range posts {
		postsDTO[i] = NewPostResponseDTO(post)
	}
	helpers.WriteJSON(w, http.StatusOK, postsDTO)
}

func (c *PostController) CreatePost(w http.ResponseWriter, r *http.Request) {
	var dto CreatePostDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		writeAppError(w, err)
		return
	}
	err = dto.Validate()
	if err != nil {
		writeAppError(w, err)
		return
	}
	userID := r.Context().Value("user_id").(string)
	p := domain.Post{
		Type:     domain.PostType(dto.Type),
		Title:    dto.Title,
		Text:     dto.Text,
		URL:      dto.URL,
		Category: dto.Category,
	}
	result, err := c.postService.CreatePost(p, userID)
	if err != nil {
		writeAppError(w, err)
		return
	}
	helpers.WriteJSON(w, http.StatusCreated, NewPostResponseDTO(result))
}

func (c *PostController) DeletePost(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["POST_ID"]
	if postID == "" {
		writeAppError(w, domain.ErrInvalidData)
		return
	}
	userID := r.Context().Value("user_id").(string)
	err := c.postService.DeletePost(postID, userID)
	if err != nil {
		writeAppError(w, err)
		return
	}
	helpers.WriteJSON(w, http.StatusNoContent, nil)
}

func (c *PostController) AddComment(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["POST_ID"]
	if postID == "" {
		writeAppError(w, domain.ErrInvalidData)
		return
	}
	var dto CreateCommentDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		writeAppError(w, err)
		return
	}
	if dto.Body == "" {
		writeAppError(w, domain.ErrInvalidData)
		return
	}
	userID := r.Context().Value("user_id").(string)
	comment := domain.Comment{
		Body: dto.Body,
	}
	result, err := c.postService.AddComment(postID, userID, comment)
	if err != nil {
		writeAppError(w, err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, result)
}

func (c *PostController) DeleteComment(w http.ResponseWriter, r *http.Request) {
	commentID := mux.Vars(r)["COMMENT_ID"]
	postID := mux.Vars(r)["POST_ID"]
	userID := r.Context().Value("user_id").(string)

	if postID == "" || userID == "" || commentID == "" {
		writeAppError(w, domain.ErrInvalidData)
		return
	}

	err := c.postService.DeleteComment(postID, commentID, userID)
	if err != nil {
		writeAppError(w, err)
		return
	}
	helpers.WriteJSON(w, http.StatusNoContent, nil)
}

func (c *PostController) UpvotePost(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["POST_ID"]
	if postID == "" {
		writeAppError(w, domain.ErrInvalidData)
		return
	}
	userID := r.Context().Value("user_id").(string)
	result, err := c.postService.UpdateVote(postID, userID, 1)
	if err != nil {
		writeAppError(w, err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, NewPostResponseDTO(result))
}

func (c *PostController) DownvotePost(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["POST_ID"]
	if postID == "" {
		writeAppError(w, domain.ErrInvalidData)
		return
	}
	userID := r.Context().Value("user_id").(string)
	result, err := c.postService.UpdateVote(postID, userID, -1)
	if err != nil {
		writeAppError(w, err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, NewPostResponseDTO(result))
}

func (c *PostController) UnvotePost(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["POST_ID"]
	if postID == "" {
		writeAppError(w, domain.ErrInvalidData)
		return
	}
	userID := r.Context().Value("user_id").(string)
	result, err := c.postService.UpdateVote(postID, userID, 0)
	if err != nil {
		writeAppError(w, err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, NewPostResponseDTO(result))
}

func writeAppError(w http.ResponseWriter, err error) {
	var vErr *domain.ValidationError
	switch {
	case errors.Is(err, domain.ErrNotFound):
		helpers.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "post not found"})
	case errors.Is(err, domain.ErrForbidden):
		helpers.WriteJSON(w, http.StatusForbidden, map[string]string{"error": "forbidden"})
	case errors.Is(err, domain.ErrInvalidData):
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid data"})
	case errors.As(err, &vErr):
		helpers.WriteJSON(w, http.StatusUnprocessableEntity, map[string]string{
			"error":   "validation failed",
			"field":   vErr.Field,
			"message": vErr.Message,
		})
	default:
		helpers.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}
}
