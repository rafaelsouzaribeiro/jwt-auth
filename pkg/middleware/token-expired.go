package middleware

import (
	"time"

	"github.com/golang-jwt/jwt"
)

func (c *Credential) TokenExpired(tokenString string) bool {
	token, _ := jwt.Parse(tokenString, nil)
	claims, _ := token.Claims.(jwt.MapClaims)
	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
	return expirationTime.Before(time.Now())
}
