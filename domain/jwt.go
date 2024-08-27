package domain

import (
	"github.com/dgrijalva/jwt-go"
)

type JWTClaim struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	Exp    int64  `json:"exp"`
	jwt.StandardClaims
}
