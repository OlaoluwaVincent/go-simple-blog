// utils/jwt.go

package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/olaoluwavincent/full-course/internal/env"
)

var JWTSecret = []byte(env.GetEnvString("JWT_SECRET", "")) // Load from env!

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT token for the user
// Expires in 24 hours by default (customize as needed)
func GenerateToken(userID int64, username, email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    env.GetEnvString("APP_NAME", "MyApp"),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

func ParseToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
