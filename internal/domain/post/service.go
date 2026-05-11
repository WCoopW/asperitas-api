package post

import (
	"reddit/internal/domain/user"

	"go.uber.org/zap"
)

type PostService interface {
	GetPostByID(id string) (Post, error)
	GetPosts() ([]Post, error)
	GetPostsByCategory(category string) ([]Post, error)
	GetPostsByUser(userLogin string) ([]Post, error)
	AddComment(postID string, userID string, comment Comment) (Comment, error)
	CreatePost(post Post, userID string) (Post, error)
	DeletePost(id string, userID string) error
	DeleteComment(postID string, commentID string, userID string) error
	UpdateVote(id string, userID string, value int) (Post, error)
}

type PostRepository interface {
	GetPostByID(id string) (Post, error)
	GetPosts() ([]Post, error)
	GetPostsByCategory(category string) ([]Post, error)
	GetPostsByUser(userLogin string) ([]Post, error)
	AddComment(postID string, comment Comment) (Comment, error)
	CreatePost(post Post) (Post, error)
	DeletePost(id string) error
	DeleteComment(postID string, commentID string) error
	UpdateVote(id string, userID string, value int) (Post, error)
}

type service struct {
	repo        PostRepository
	userService user.UserService
	logger      *zap.SugaredLogger
}

func New(repo PostRepository, userService user.UserService, logger *zap.SugaredLogger) PostService {
	return &service{repo: repo, userService: userService, logger: logger}
}

func (s *service) GetPostByID(id string) (Post, error) {
	return s.repo.GetPostByID(id)
}

func (s *service) GetPosts() ([]Post, error) {
	return s.repo.GetPosts()
}

func (s *service) GetPostsByCategory(category string) ([]Post, error) {
	return s.repo.GetPostsByCategory(category)
}

func (s *service) GetPostsByUser(userLogin string) ([]Post, error) {
	return s.repo.GetPostsByUser(userLogin)
}

func (s *service) AddComment(postID string, userID string, comment Comment) (Comment, error) {
	author, err := s.userService.GetUserByID(userID)
	if err != nil {
		return Comment{}, err
	}
	comment.Author = author
	return s.repo.AddComment(postID, comment)
}

func (s *service) CreatePost(post Post, userID string) (Post, error) {
	author, err := s.userService.GetUserByID(userID)
	if err != nil {
		return Post{}, err
	}
	post.Author = author
	return s.repo.CreatePost(post)
}

func (s *service) DeletePost(id string, userID string) error {
	p, err := s.repo.GetPostByID(id)
	if err != nil {
		return err
	}

	if p.Author.ID != userID {
		return ErrForbidden
	}
	return s.repo.DeletePost(id)
}

func (s *service) DeleteComment(id, commentID, userID string) error {
	p, err := s.repo.GetPostByID(id)
	if err != nil {
		return err
	}
	commentIndex := -1
	for i := range p.Comments {
		if p.Comments[i].ID == commentID {
			if p.Comments[i].Author.ID != userID {
				return ErrForbidden
			}
			commentIndex = i
			break
		}
	}
	if commentIndex == -1 {
		return ErrNotFound
	}

	return s.repo.DeleteComment(id, commentID)
}

func (s *service) UpdateVote(id string, userID string, value int) (Post, error) {
	return s.repo.UpdateVote(id, userID, value)
}
