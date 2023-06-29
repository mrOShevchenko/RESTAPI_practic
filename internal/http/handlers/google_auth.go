package handlers

import (
	"Nix_trainee_practic/config"
	"Nix_trainee_practic/internal/http/requests"
	"Nix_trainee_practic/internal/http/response"
	"Nix_trainee_practic/internal/service"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
)

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

type Oauth struct {
	us service.UserService
	as service.AuthService
}

func NewOauth(u service.UserService, a service.AuthService) Oauth {
	return Oauth{
		us: u,
		as: a,
	}
}

func (o Oauth) GetInfo(ctx echo.Context) error {
	googleConfig := config.LoadOAUTHConfiguration()
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("Google oauth error: %s", err))
	}
	state := base64.URLEncoding.EncodeToString(b)
	url := googleConfig.AuthCodeURL(state)
	log.Println(url)
	err = ctx.Redirect(http.StatusTemporaryRedirect, url)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("Google oauth error, could not redirect url: %s", err))
	}
	return response.MessageResponse(ctx, http.StatusOK, "Success")
}

func (o Oauth) CallBackRegister(ctx echo.Context) error {
	cfg := config.LoadOAUTHConfiguration()

	token, err := cfg.Exchange(context.Background(), ctx.FormValue("code"))
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, fmt.Sprintf("Google oauth error, code exchange wrong: %s", err))
	}
	resp, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, fmt.Sprintf("Google oauth error, failed gatting user info: %s", err))
	}
	defer resp.Body.Close()
	var usr requests.RegisterOauth2
	err = json.NewDecoder(resp.Body).Decode(&usr)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "Could not decode user data")
	}
	usr.Password = usr.ID + usr.Email + config.LoadOAUTHConfiguration().ClientID
	userFromRegister := usr.RegisterToUser()
	user, err := o.as.Register(userFromRegister)
	if err != nil {
		if strings.HasSuffix(err.Error(), "invalid credentials user exist") {
			user, err = o.us.FindByEmail(userFromRegister.Email)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("User not exist: %s", err))
			}
			if user.ID == 0 || (bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userFromRegister.Password)) != nil) {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
			}
		}
	}
	accessToken, refreshToken, exp, err := o.as.Login(requests.LoginAuth{
		Email:    userFromRegister.Email,
		Password: userFromRegister.Password,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	res := response.NewLoginResponse(accessToken, refreshToken, exp)
	return response.Response(ctx, http.StatusOK, res)
}
