package service

import (
	"Nix_trainee_practic/config"
	"Nix_trainee_practic/internal/http/requests"
	"Nix_trainee_practic/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/upper/db/v4"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/sync/errgroup"
	"time"
)

const (
	refresh = 48
	access  = 2
	LogOF   = 10
)

//go:generate mockery --dir . --name AuthService --output ./mocks
type AuthService interface {
	Register(user models.User) (models.User, error)
	Login(user requests.LoginAuth) (string, string, int64, error)
	ValidateJWT(tokenUID string, userID int64, isRefresh bool) (models.User, error)
}

type authService struct {
	userService UserService
	config      config.Configuration
	r           *redis.Client
}

func NewAuth(us UserService, cf config.Configuration, red *redis.Client) AuthService {
	return authService{
		userService: us,
		config:      cf,
		r:           red,
	}
}

func (a authService) Register(user models.User) (models.User, error) {
	_, err := a.userService.FindByEmail(user.Email)
	if err == nil {
		return models.User{}, fmt.Errorf("auth service error register invalid credentials user exist")
	} else if !errors.Is(err, db.ErrNoMoreRows) {
		return models.User{}, fmt.Errorf("auth service error register")
	}
	user, err = a.userService.Save(user)
	if err != nil {
		return models.User{}, fmt.Errorf("auth service error register save user: %w", err)
	}
	return user, nil
}

func (a authService) Login(user requests.LoginAuth) (string, string, int64, error) {
	u, err := a.userService.FindByEmail(user.Email)
	if err != nil {
		if errors.Is(err, db.ErrNoMoreRows) {
			return "", "", 0, fmt.Errorf("auth service error login, invalid credentials user not exist: %w", err)
		}
		return "", "", 0, fmt.Errorf("auth service error login user invalid email or password: %w", err)
	}
	valid := a.checkPasswordHash(user.Password, u.Password)
	if !valid {
		return "", "", 0, fmt.Errorf("auth service error login user invalid email or password: %w", err)
	}
	accessToken, accessUID, exp, err := createToken(u, access, a.config.AccessSecret)
	if err != nil {
		return "", "", 0, fmt.Errorf("auth service error login: %w", err)
	}
	refreshToken, refreshUID, _, err := createToken(u, refresh, a.config.RefreshSecret)
	if err != nil {
		return "", "", 0, fmt.Errorf("auth service error login: %w", err)
	}

	tokensJSON, err := json.Marshal(RedisToken{
		AccessID:  accessUID,
		RefreshID: refreshUID,
	})
	if err != nil {
		return "", "", 0, fmt.Errorf("auth service error, couldn't marshal token pair, %w", err)
	}

	a.r.Set(fmt.Sprintf("token-%d", u.ID), string(tokensJSON), time.Minute*LogOF)
	return accessToken, refreshToken, exp, err
}

func (a authService) ValidateJWT(tokenUID string, userID int64, isRefresh bool) (user models.User, err error) {
	var g errgroup.Group
	g.Go(func() error {
		tokensJSON, err := a.r.Get(fmt.Sprintf("token-%d", userID)).Result()
		if err != nil {
			return fmt.Errorf("auth service error validate token: %w", err)
		}
		var redisToken RedisToken
		err = json.Unmarshal([]byte(tokensJSON), &redisToken)

		var uid string
		if isRefresh {
			uid = redisToken.RefreshID
		} else {
			uid = redisToken.AccessID
		}

		if err != nil || uid != tokenUID {
			return fmt.Errorf("token not exist: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		user, err = a.userService.FindByID(userID)
		if err != nil {
			return fmt.Errorf("auth service error validate jwt invalid credentials user not exist, %w", err)
		}
		return nil
	})

	err = g.Wait()

	return user, err
}

func createToken(user models.User, expireTime int, secret string) (string, string, int64, error) {
	exp := time.Now().Add(time.Hour * time.Duration(expireTime)).Unix()
	uid := uuid.New().String()
	claimsAccess := JWTClaim{
		Name: user.Name,
		ID:   user.ID,
		UID:  uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsAccess)
	t, err := token.SignedString([]byte(secret))

	return t, uid, exp, err
}

func (a authService) checkPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

type JWTClaim struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
	UID  string `json:"uid"`
	jwt.StandardClaims
}

type RedisToken struct {
	AccessID  string `json:"access"`
	RefreshID string `json:"refresh"`
}
