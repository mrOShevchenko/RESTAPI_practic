package http

import (
	"Nix_trainee_practic/config"
	"Nix_trainee_practic/config/container"
	"Nix_trainee_practic/internal/http/validators"
	"Nix_trainee_practic/internal/repository"
	MW "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func EchoRouter(s *Server, cont container.Container) {

	e := s.Echo
	e.GET("/auth/google/login", cont.Oauth.GetInfo)
	e.GET("/auth/google/callback", cont.Oauth.CallBackRegister)
	e.Use(MW.Logger())
	e.Validator = validators.NewValidator()

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.POST("/register", cont.Register.Register)
	e.POST("/login", cont.Register.Login)

	v1 := e.Group("/api/v1")
	v1.GET("", repository.PingHandler)

	authMW := cont.AuthMiddleware.JWT(config.GetConfiguration().AccessSecret)
	validToken := cont.AuthMiddleware.ValidateJWT()
	commRouter := v1.Group("/comments/")
	postRouter := v1.Group("/posts/")

	commRouter.Use(authMW, validToken)
	postRouter.Use(authMW, validToken)

	commRouter.POST("save/:post_id", cont.Comment.SaveComment)
	commRouter.GET("comment/:id", cont.Comment.GetComment)
	commRouter.PUT("update/:id", cont.Comment.UpdateComment)
	commRouter.DELETE("delete/:id", cont.Comment.DeleteComment)

	postRouter.POST("save", cont.Post.SavePost)
	postRouter.GET("post/:id", cont.Post.GetPost)
	postRouter.PUT("update/:id", cont.Post.UpdatePost)
	postRouter.DELETE("delete/:id", cont.Post.DeletePost)
}
