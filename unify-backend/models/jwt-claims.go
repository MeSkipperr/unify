package models

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	Sub   string `json:"sub"`             // user id
	Name  string `json:"name,omitempty"`  // user name (optional)
	Type  string `json:"type,omitempty"`  // access | refresh
	jwt.RegisteredClaims
}
