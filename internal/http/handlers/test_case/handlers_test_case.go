package test_case

import (
	"Nix_trainee_practic/internal/http/validators"
	"Nix_trainee_practic/internal/service"
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http/httptest"
	"strings"
	"time"
)

type TestCase struct {
	TestName    string
	Request     Request
	RequestBody interface{}
	HandlerFunc func(c echo.Context) error
	Expected    ExpectedResponse
}

type Request struct {
	Method    string
	Url       string
	PathParam *PathParam
}

type PathParam struct {
	Name  string
	Value string
}

type ExpectedResponse struct {
	StatusCode int
	BodyPart   string
}

func PrepareContextFromTestCase(test TestCase) (c echo.Context, recorder *httptest.ResponseRecorder) {
	e := echo.New()
	e.Validator = validators.NewValidator()
	requestJson, _ := json.Marshal(test.RequestBody)
	request := httptest.NewRequest(test.Request.Method, test.Request.Url, strings.NewReader(string(requestJson)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder = httptest.NewRecorder()
	c = e.NewContext(request, recorder)

	if test.Request.PathParam != nil {
		c.SetParamNames(test.Request.PathParam.Name)
		c.SetParamValues(test.Request.PathParam.Value)
	}

	return
}

func Token() *jwt.Token {
	exp := time.Now().Add(time.Hour * 2).Unix()
	claimsAccess := &service.JWTClaim{
		Name: "Name",
		ID:   int64(1),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}
	tokenReturn := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsAccess)
	return tokenReturn
}
