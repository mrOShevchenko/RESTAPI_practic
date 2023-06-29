package container

import (
	"Nix_trainee_practic/config"
	"Nix_trainee_practic/internal/http/handlers"
	"Nix_trainee_practic/internal/repository"
	"Nix_trainee_practic/internal/service"
	"Nix_trainee_practic/middleware"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type Container struct {
	Services
	Handlers
	Middleware
}

type Services struct {
	service.CommentService
	service.PostService
	service.UserService
	service.AuthService
}

type Handlers struct {
	handlers.Comment
	handlers.Post
	handlers.Register
	handlers.Oauth
}

type Middleware struct {
	middleware.AuthMiddleware
}

func New(conf config.Configuration) Container {
	sess := getDbSess(conf)
	newRedis := getRedis(conf)

	userRepository := repository.NewUserRepo(sess)
	passwordGenerator := service.NewGeneratePasswordHash(bcrypt.DefaultCost)
	userService := service.NewUser(userRepository, passwordGenerator)
	authService := service.NewAuth(userService, conf, newRedis)
	registerController := handlers.NewRegister(authService)
	oauthController := handlers.NewOauth(userService, authService)

	postRepository := repository.NewPostRepo(sess)
	postService := service.NewPost(postRepository)

	commentRepository := repository.NewCommentRepo(sess)
	commentService := service.NewComment(commentRepository, userService, postService)
	commentHandler := handlers.NewComment(commentService)

	postHandler := handlers.NewPost(postService, commentService)

	authMiddleware := middleware.NewMiddleware(authService, newRedis)

	return Container{
		Services: Services{
			commentService,
			postService,
			userService,
			authService,
		},
		Handlers: Handlers{
			commentHandler,
			postHandler,
			registerController,
			oauthController,
		},
		Middleware: Middleware{
			authMiddleware,
		},
	}
}

func getDbSess(conf config.Configuration) db.Session {
	sess, err := postgresql.Open(
		postgresql.ConnectionURL{
			User:     conf.DatabaseUser,
			Host:     conf.DatabaseHost,
			Password: conf.DatabasePassword,
			Database: conf.DatabaseName,
		})
	if err != nil {
		log.Fatalf("Unable to create new DB session: %q\n", err)
	}
	return sess
}

func getRedis(conf config.Configuration) *redis.Client {
	addr := fmt.Sprintf("%s:%s", conf.RedisHost, conf.RedisPort)
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}
