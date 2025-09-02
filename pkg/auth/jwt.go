package auth

import (
	"errors"
	"tesodev-korpes/pkg/customError"
	"tesodev-korpes/shared/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(config.GetAuthConfig().JWTSecret)

type Claims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

func GenerateJWT(Id string) (string, error) {
	claims := &Claims{
		ID: Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func VerifyJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, customError.NewUnauthorized(customError.MissingAuthToken)
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token is invalid or expired")
	}

	return claims, nil
}
