package middleware

import (
	"Nix_trainee_practic/internal/http/response"
	"Nix_trainee_practic/internal/service"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"

	MW "github.com/labstack/echo/v4/middleware"
)

type AuthMiddleware interface {
	JWT(secret string) echo.MiddlewareFunc
	ValidateJWT() echo.MiddlewareFunc
}

type authMiddleware struct {
	authService service.AuthService
	r           *redis.Client
}

func NewMiddleware(as service.AuthService, red *redis.Client) AuthMiddleware {
	return authMiddleware{
		authService: as,
		r:           red,
	}
}

func (m authMiddleware) ValidateJWT() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			claims := token.Claims.(*service.JWTClaim)

			user, err := m.authService.ValidateJWT(claims.UID, claims.ID, false)
			if err != nil {
				return response.MessageResponse(c, http.StatusUnauthorized, "Not authorized")
			}

			c.Set("currentUser", user)

			go func() {
				m.r.Expire(fmt.Sprintf("token-%d", claims.ID), time.Minute*service.LogOF)
			}()
			return next(c)
		}
	}
}

func (m authMiddleware) JWT(secret string) echo.MiddlewareFunc {
	config := MW.JWTConfig{
		ErrorHandler: func(err error) error {
			return &echo.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "Not authorized",
			}
		},
		SigningKey: []byte(secret),
		Claims:     &service.JWTClaim{},
	}
	return MW.JWTWithConfig(config)
}
