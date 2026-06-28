package app

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"

	authhandler "reddit/internal/application/auth"
	"reddit/internal/application/middleware"
	post_app "reddit/internal/application/post"
	"reddit/internal/config"
	"reddit/internal/domain/auth"
	"reddit/internal/domain/post"
	"reddit/internal/domain/user"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type DIContainer struct {
	AuthService    auth.AuthService
	TokenValidator auth.TokenValidator
	AuthController authhandler.AuthController
	PostController post_app.PostController
	UserRepository user.UserRepository
	PostRepository post.PostRepository
	UserService    user.UserService
	PostService    post.PostService
}
type App struct {
	cfg       *config.Config
	mux       *mux.Router
	di        *DIContainer
	logger    *zap.SugaredLogger
	staticDir string
}

func New(cfg *config.Config, di *DIContainer, logger *zap.SugaredLogger, staticDir string) (*App, error) {
	app := &App{
		cfg:       cfg,
		di:        di,
		mux:       mux.NewRouter(),
		logger:    logger,
		staticDir: staticDir,
	}
	app.logger.Info("App initialized")
	app.setupStatic()
	app.setuphandlers()

	return app, nil
}

func (app *App) Start(ctx context.Context) error {
	app.logger.Infof("Server is listening on port %d", app.cfg.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", app.cfg.Port), app.mux)
	if err != nil {
		app.logger.Error("error start server", zap.Error(err))
		return err
	}
	return nil
}

func (app *App) setupStatic() {
	fs := http.FileServer(http.Dir(app.staticDir))
	app.mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	app.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(app.staticDir, "html", "index.html"))
	})
}

func (app *App) setuphandlers() {
	authMiddleware := middleware.NewAuthMiddleware(app.di.TokenValidator)
	loggerMiddleware := middleware.NewLoggerMiddleware(app.logger)
	app.mux.Use(loggerMiddleware.Log)

	publicMux := app.mux.PathPrefix("/api").Subrouter()
	publicMux.HandleFunc("/register", app.di.AuthController.Register).Methods("POST")
	publicMux.HandleFunc("/login", app.di.AuthController.Login).Methods("POST")
	publicMux.HandleFunc("/logout", app.di.AuthController.Logout).Methods("POST")
	publicMux.HandleFunc("/posts/", app.di.PostController.GetPosts).Methods("GET")
	publicMux.HandleFunc("/posts/{CATEGORY_NAME}", app.di.PostController.GetPostsByCategory).Methods("GET")
	publicMux.HandleFunc("/post/{POST_ID}", app.di.PostController.GetPostByID).Methods("GET")
	publicMux.HandleFunc("/user/{USER_LOGIN}", app.di.PostController.GetPostByUser).Methods("GET")

	protectedMux := app.mux.PathPrefix("/api").Subrouter()
	protectedMux.Use(authMiddleware.WithAuth)
	protectedMux.HandleFunc("/posts", app.di.PostController.CreatePost).Methods("POST")
	protectedMux.HandleFunc("/post/{POST_ID}", app.di.PostController.AddComment).Methods("POST")
	protectedMux.HandleFunc("/post/{POST_ID}/{COMMENT_ID}", app.di.PostController.DeleteComment).Methods("DELETE")
	protectedMux.HandleFunc("/post/{POST_ID}/upvote", app.di.PostController.UpvotePost).Methods("GET")
	protectedMux.HandleFunc("/post/{POST_ID}/downvote", app.di.PostController.DownvotePost).Methods("GET")
	protectedMux.HandleFunc("/post/{POST_ID}/unvote", app.di.PostController.UnvotePost).Methods("GET")
	protectedMux.HandleFunc("/post/{POST_ID}", app.di.PostController.DeletePost).Methods("DELETE")
}
