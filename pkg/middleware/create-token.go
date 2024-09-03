package middleware

import (
	"time"

	"github.com/golang-jwt/jwt"
)

func (c *Credential) CreateToken(claims map[string]interface{}) (string, error) {

	if _, ok := claims["exp"]; !ok {
		claims["exp"] = time.Now().Add(time.Second * time.Duration(c.ExpireIn)).Unix() // Default expiration (1 hour)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))

	tokenString, err := token.SignedString(c.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
