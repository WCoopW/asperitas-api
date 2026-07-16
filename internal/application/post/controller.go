package post

import (
	"encoding/json"
	"errors"
	"net/http"

	"reddit/internal/apperrors"
	d_comment "reddit/internal/domain/comments"
	d_post "reddit/internal/domain/post"

	"reddit/pkg/helpers"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type PostController struct {
	postService    d_post.PostService
	commentService d_comment.CommentService
	logger         *zap.SugaredLogger
}

func New(postService d_post.PostService, commentService d_comment.CommentService, logger *zap.SugaredLogger) *PostController {
	return &PostController{
		postService:    postService,
		commentService: commentService,
		logger:         logger,
	}
}

func (c *PostController) GetPosts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := helpers.CtxDefaultTimeout(r.Context(), nil)
	defer cancel()
	posts, err := c.postService.GetPosts(ctx, d_post.PostFilter{}, 0, 0)
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
	ctx, cancel := helpers.CtxDefaultTimeout(r.Context(), nil)
	defer cancel()
	posts, err := c.postService.GetPostsByCategory(ctx, category)
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
	ctx, cancel := helpers.CtxDefaultTimeout(r.Context(), nil)
	defer cancel()
	post, err := c.postService.GetPostByID(ctx, postID)
	if err != nil {
		writeAppError(w, err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, NewPostResponseDTO(post))
}

func (c *PostController) GetPostByUser(w http.ResponseWriter, r *http.Request) {
	userLogin := mux.Vars(r)["USER_LOGIN"]
	ctx, cancel := helpers.CtxDefaultTimeout(r.Context(), nil)
	defer cancel()
	posts, err := c.postService.GetPostsByUser(ctx, userLogin)
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
	ctx, cancel := helpers.CtxDefaultTimeout(r.Context(), nil)
	defer cancel()
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
	p := d_post.Post{
		Type:     d_post.PostType(dto.Type),
		Title:    dto.Title,
		Text:     dto.Text,
		URL:      dto.URL,
		Category: dto.Category,
	}
	result, err := c.postService.CreatePost(ctx, p, userID)
	if err != nil {
		writeAppError(w, err)
		return
	}
	helpers.WriteJSON(w, http.StatusCreated, NewPostResponseDTO(result))
}

func (c *PostController) DeletePost(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["POST_ID"]
	ctx, cancel := helpers.CtxDefaultTimeout(r.Context(), nil)
	defer cancel()
	if postID == "" {
		writeAppError(w, apperrors.ErrInvalidData)
		return
	}
	userID := r.Context().Value("user_id").(string)
	err := c.postService.DeletePost(ctx, postID, userID)
	if err != nil {
		writeAppError(w, err)
		return
	}
	helpers.WriteNoContent(w)
}

func (c *PostController) AddComment(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["POST_ID"]
	ctx, cancel := helpers.CtxDefaultTimeout(r.Context(), nil)
	defer cancel()
	if postID == "" {
		writeAppError(w, apperrors.ErrInvalidData)
		return
	}
	var dto CreateCommentDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		writeAppError(w, err)
		return
	}
	if dto.Body == "" {
		writeAppError(w, apperrors.ErrInvalidData)
		return
	}
	userID := r.Context().Value("user_id").(string)
	result, err := c.commentService.Create(ctx, postID, userID, userID)
	if err != nil {
		writeAppError(w, err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, result)
}

func (c *PostController) DeleteComment(w http.ResponseWriter, r *http.Request) {
	commentID := mux.Vars(r)["COMMENT_ID"]
	postID := mux.Vars(r)["POST_ID"]
	ctx, cancel := helpers.CtxDefaultTimeout(r.Context(), nil)
	defer cancel()
	userID := r.Context().Value("user_id").(string)

	if postID == "" || userID == "" || commentID == "" {
		writeAppError(w, apperrors.ErrInvalidData)
		return
	}

	err := c.commentService.Delete(ctx, commentID, userID)
	if err != nil {
		writeAppError(w, err)
		return
	}
	helpers.WriteNoContent(w)
}

func (c *PostController) UpvotePost(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["POST_ID"]
	ctx, cancel := helpers.CtxDefaultTimeout(r.Context(), nil)
	defer cancel()
	if postID == "" {
		writeAppError(w, apperrors.ErrInvalidData)
		return
	}
	userID := r.Context().Value("user_id").(string)
	result, err := c.postService.UpdateVote(ctx, postID, userID, 1)
	if err != nil {
		writeAppError(w, err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, NewPostResponseDTO(result))
}

func (c *PostController) DownvotePost(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["POST_ID"]
	ctx, cancel := helpers.CtxDefaultTimeout(r.Context(), nil)
	defer cancel()
	if postID == "" {
		writeAppError(w, apperrors.ErrInvalidData)
		return
	}
	userID := r.Context().Value("user_id").(string)
	result, err := c.postService.UpdateVote(ctx, postID, userID, -1)
	if err != nil {
		writeAppError(w, err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, NewPostResponseDTO(result))
}

func (c *PostController) UnvotePost(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["POST_ID"]
	ctx, cancel := helpers.CtxDefaultTimeout(r.Context(), nil)
	defer cancel()
	if postID == "" {
		writeAppError(w, apperrors.ErrInvalidData)
		return
	}
	userID := r.Context().Value("user_id").(string)
	result, err := c.postService.UpdateVote(ctx, postID, userID, 0)
	if err != nil {
		writeAppError(w, err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, NewPostResponseDTO(result))
}

func writeAppError(w http.ResponseWriter, err error) {
	var vErr *apperrors.ValidationError
	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		helpers.WriteError(w, http.StatusNotFound, "post not found")
	case errors.Is(err, apperrors.ErrForbidden):
		helpers.WriteError(w, http.StatusForbidden, "forbidden")
	case errors.Is(err, apperrors.ErrInvalidData):
		helpers.WriteError(w, http.StatusBadRequest, "invalid data")
	case errors.As(err, &vErr):
		helpers.WriteJSON(w, http.StatusUnprocessableEntity, map[string]string{
			"message": vErr.Message,
			"field":   vErr.Field,
		})
	default:
		helpers.WriteError(w, http.StatusInternalServerError, "internal error")
	}
}
