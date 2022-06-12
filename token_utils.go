package main

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	Uuid string
}

type CustomClaims struct {
	User
	jwt.RegisteredClaims
}

const ACCESS_TOKEN_TTL = time.Minute * 10

var accessKey = []byte(os.Getenv("ACCESS_KEY"))
var refreshKey = []byte(os.Getenv("REFRESH_KEY"))

func generateAccessToken(user User) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS512, CustomClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ACCESS_TOKEN_TTL)),
		},
	})

	tkn, err := t.SignedString(accessKey)
	if err != nil {
		return "", err
	}
	return tkn, nil
}

func authorizeAccessToken(tokenStr string) (*CustomClaims, error) {

	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return accessKey, nil
	})
	if err != nil {
		return nil, err
	}
	tkn, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, err
	}
	if err = tkn.Valid(); err != nil {
		return nil, err
	}
	return tkn, nil
}

func generateRefreshToken(user User) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS512, CustomClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	})

	tkn, err := t.SignedString(refreshKey)
	if err != nil {
		return "", err
	}
	return tkn, nil
}

func authorizeRefreshToken(tokenStr string) (*CustomClaims, error) {

	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return refreshKey, nil
	})
	if err != nil {
		return nil, err
	}

	tkn, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("bad claims")
	}

	if err = tkn.Valid(); err != nil {
		return nil, err
	}
	return tkn, nil
}
