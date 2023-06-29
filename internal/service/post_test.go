package service

import (
	"github.com/golang-jwt/jwt"
	"time"
)

func token() *jwt.Token {
	exp := time.Now().Add(time.Hour * access).Unix()
	claimsAccess := &JWTClaim{
		Name: "Name",
		ID:   int64(1),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}
	tokenReturn := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsAccess)
	return tokenReturn
}
