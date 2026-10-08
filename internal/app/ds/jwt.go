package ds

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID uint   `json:"user_id"`
	Login  string `json:"login"`
	Role   string `json:"role"`
}
