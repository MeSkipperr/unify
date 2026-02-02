package utils

import (
	"errors"
	"time"
	"unify-backend/models"

	"github.com/golang-jwt/jwt/v5"
)

var (
	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 7 * 24 * time.Hour
)

func GenerateAccessToken(userID, name, secret string) (string, error) {
	claims := models.JWTClaims{
		Sub:  userID,
		Name: name,
		Type: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(userID, secret string) (string, error) {
	claims := models.JWTClaims{
		Sub:  userID,
		Type: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func VerifyJWT(tokenStr, secret string) (*models.JWTClaims, error) {
	claims := &models.JWTClaims{}

	token, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	return claims, nil
}

