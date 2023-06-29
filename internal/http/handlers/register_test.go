package handlers

import (
	"Nix_trainee_practic/internal/http/handlers/test_case"
	"Nix_trainee_practic/internal/http/requests"
	"Nix_trainee_practic/internal/models"
	"Nix_trainee_practic/internal/service"
	"Nix_trainee_practic/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/upper/db/v4"
	"net/http"
	"testing"
)

func TestRegisterHandler_Register(t *testing.T) {
	userMockRequest := requests.RegisterAuth{
		Email:    "test@mail.com",
		Name:     "Name",
		Password: "testPassword",
	}

	userMockDomain := models.User{
		ID:       1,
		Email:    "test@mail.com",
		Name:     "Name",
		Password: "testPassword",
	}

	requestRegister := test_case.Request{
		Method: http.MethodPost,
		Url:    "/register",
	}
	handleSuccessCreate := func(c echo.Context) error {
		mockAuth := func(user requests.RegisterAuth) service.AuthService {
			mock := mocks.NewAuthService(t)
			userDomain := models.User{
				Email:    user.Email,
				Name:     user.Name,
				Password: user.Password,
			}
			mock.On("Register", userDomain).Return(userMockDomain, nil).Times(1)
			return mock
		}(userMockRequest)
		return NewRegister(mockAuth).Register(c)
	}

	handleErrorCreate := func(c echo.Context) error {
		mockAuth := func(user requests.RegisterAuth) service.AuthService {
			mock := mocks.NewAuthService(t)
			userDomain := models.User{
				Email:    user.Email,
				Name:     user.Name,
				Password: user.Password,
			}
			mock.On("Register", userDomain).Return(models.User{}, db.ErrCollectionDoesNotExist).Times(1)
			return mock
		}(userMockRequest)
		return NewRegister(mockAuth).Register(c)
	}

	handleMock := func(c echo.Context) error {
		mockAuth := func() service.AuthService {
			return mocks.NewAuthService(t)
		}()
		return NewRegister(mockAuth).Register(c)
	}

	cases := []test_case.TestCase{
		{
			TestName:    "RegisterUser OK",
			Request:     requestRegister,
			RequestBody: userMockRequest,
			HandlerFunc: handleSuccessCreate,
			Expected: test_case.ExpectedResponse{
				StatusCode: 201,
				BodyPart:   "{\"id\":1,\"email\":\"test@mail.com\",\"name\":\"Name\"}\n"},
		},
		{
			TestName:    "RegisterUser Error",
			Request:     requestRegister,
			RequestBody: userMockRequest,
			HandlerFunc: handleErrorCreate,
			Expected: test_case.ExpectedResponse{
				StatusCode: 500,
				BodyPart:   "{\"code\":500,\"error\":\"Could not save new user: upper: collection does not exist\"}\n"},
		},
		{
			TestName:    "Error decode user data",
			Request:     requestRegister,
			RequestBody: "",
			HandlerFunc: handleMock,
			Expected: test_case.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "{\"code\":400,\"error\":\"Could not decode user data\"}\n"},
		},
		{
			TestName: "Error validate user data",
			Request:  requestRegister,
			RequestBody: requests.RegisterAuth{
				Email:    "test@email.com",
				Password: "testPassword",
				Name:     "",
			},
			HandlerFunc: handleMock,
			Expected: test_case.ExpectedResponse{
				StatusCode: 422,
				BodyPart:   "{\"code\":422,\"error\":\"Could not validate user data\"}\n"},
		},
	}
	for _, test := range cases {
		t.Run(test.TestName, func(t *testing.T) {
			c, recorder := test_case.PrepareContextFromTestCase(test)
			c.Set("user", test_case.Token())

			if assert.NoError(t, test.HandlerFunc(c)) {
				assert.Contains(t, recorder.Body.String(), test.Expected.BodyPart)
				assert.Equal(t, recorder.Code, test.Expected.StatusCode)
			}
		})
	}
}

func TestRegisterHandler_Login(t *testing.T) {
	userMockRequest := requests.LoginAuth{
		Email:    "test@mail.com",
		Password: "testPassword",
	}

	requestRegister := test_case.Request{
		Method: http.MethodPost,
		Url:    "/login",
	}
	handleSuccessLogin := func(c echo.Context) error {
		mockAuth := func(user requests.LoginAuth) service.AuthService {
			mock := mocks.NewAuthService(t)
			mock.On("Login", user).Return("access", "refresh", int64(123), nil).Times(1)
			return mock
		}(userMockRequest)
		return NewRegister(mockAuth).Login(c)
	}

	handleErrorLoginNoMoreRows := func(c echo.Context) error {
		mockAuth := func(user requests.LoginAuth) service.AuthService {
			mock := mocks.NewAuthService(t)
			mock.On("Login", user).Return("", "", int64(0), db.ErrNoMoreRows).Times(1)
			return mock
		}(userMockRequest)
		return NewRegister(mockAuth).Login(c)
	}

	handleErrorLoginInternalServerError := func(c echo.Context) error {
		mockAuth := func(user requests.LoginAuth) service.AuthService {
			mock := mocks.NewAuthService(t)
			mock.On("Login", user).Return("", "", int64(0), db.ErrCollectionDoesNotExist).Times(1)
			return mock
		}(userMockRequest)
		return NewRegister(mockAuth).Login(c)
	}

	handleMock := func(c echo.Context) error {
		mockAuth := func() service.AuthService {
			return mocks.NewAuthService(t)
		}()
		return NewRegister(mockAuth).Login(c)
	}

	cases := []test_case.TestCase{
		{
			TestName:    "LoginUser success",
			Request:     requestRegister,
			RequestBody: userMockRequest,
			HandlerFunc: handleSuccessLogin,
			Expected: test_case.ExpectedResponse{
				StatusCode: 200,
				BodyPart:   "{\"accessToken\":\"access\",\"refreshToken\":\"refresh\",\"exp\":123}\n"},
		},
		{
			TestName:    "LoginUser error login no more rows",
			Request:     requestRegister,
			RequestBody: userMockRequest,
			HandlerFunc: handleErrorLoginNoMoreRows,
			Expected: test_case.ExpectedResponse{
				StatusCode: 404,
				BodyPart:   "{\"code\":404,\"error\":\"Could not login, user not exists: upper: no more rows in this result set\"}\n"},
		},
		{
			TestName:    "LoginUser error login no more rows",
			Request:     requestRegister,
			RequestBody: userMockRequest,
			HandlerFunc: handleErrorLoginInternalServerError,
			Expected: test_case.ExpectedResponse{
				StatusCode: 500,
				BodyPart:   "{\"code\":500,\"error\":\"Could not login user: upper: collection does not exist\"}\n"},
		},
		{
			TestName:    "Error decode user data",
			Request:     requestRegister,
			RequestBody: "",
			HandlerFunc: handleMock,
			Expected: test_case.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "{\"code\":400,\"error\":\"Could not decode user data\"}\n"},
		},
		{
			TestName: "Error validate user data",
			Request:  requestRegister,
			RequestBody: requests.LoginAuth{
				Email:    "test@email.com",
				Password: "",
			},
			HandlerFunc: handleMock,
			Expected: test_case.ExpectedResponse{
				StatusCode: 422,
				BodyPart:   "{\"code\":422,\"error\":\"Could not validate user data\"}\n"},
		},
	}
	for _, test := range cases {
		t.Run(test.TestName, func(t *testing.T) {
			c, recorder := test_case.PrepareContextFromTestCase(test)
			c.Set("user", test_case.Token())

			if assert.NoError(t, test.HandlerFunc(c)) {
				assert.Contains(t, recorder.Body.String(), test.Expected.BodyPart)
				assert.Equal(t, recorder.Code, test.Expected.StatusCode)
			}
		})
	}
}
