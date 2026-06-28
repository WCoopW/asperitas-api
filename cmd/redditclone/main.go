package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	app "reddit/internal/application"
	authhandler "reddit/internal/application/auth"
	post_app "reddit/internal/application/post"
	"reddit/internal/config"
	"reddit/internal/db"
	authdomain "reddit/internal/domain/auth"
	post "reddit/internal/domain/post"
	user "reddit/internal/domain/user"
	post_repo "reddit/internal/infrastructure/post"
	user_repo "reddit/internal/infrastructure/user"
	"reddit/pkg/jwtutil"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("error load config: %s", err)
	}

	//  setup logger
	sugar := zap.Must(zap.NewProduction()).Sugar()
	defer sugar.Sync()

	db, err := db.Connect(cfg.DB)
	if err != nil {
		sugar.Error("error init database", zap.Error(err))
		return
	}
	defer db.Close()

	// init dependencies
	wd, _ := os.Getwd()
	staticDir := filepath.Join(wd, "static")
	sugar.Infof("staticDir: %s", staticDir)
	di := initDependencies(cfg, db, sugar)
	//  run app

	app, err := app.New(&cfg, di, sugar, staticDir)
	if err != nil {
		sugar.Error("error init app", zap.Error(err))
	}
	err = app.Start(context.Background())
	if err != nil {
		sugar.Error("error start app", zap.Error(err))
	}
}

func initDependencies(cfg config.Config, db *sqlx.DB, logger *zap.SugaredLogger) *app.DIContainer {
	jwtTokens, err := jwtutil.New(cfg.JWT)
	if err != nil {
		logger.Error("error init jwt", zap.Error(err))
		return nil
	}

	userRepository := user_repo.NewUserPGRepository(db)
	userService := user.New(userRepository, logger.Named("user_service"))
	authService := authdomain.New(userService, jwtTokens, logger.Named("auth_service"))
	authController := authhandler.New(authService, logger.Named("auth_controller"))

	postRepository := post_repo.NewMem(logger.Named("post_repository"))
	postService := post.New(postRepository, userService, logger.Named("post_service"))
	postController := post_app.New(postService, logger.Named("post_controller"))
	return &app.DIContainer{
		AuthController: *authController,
		PostController: *postController,
		AuthService:    authService,
		TokenValidator: jwtTokens,
		UserRepository: userRepository,
		PostRepository: postRepository,
		UserService:    userService,
		PostService:    postService,
	}
}
