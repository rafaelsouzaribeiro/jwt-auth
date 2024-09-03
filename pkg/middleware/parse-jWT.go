package middleware

import "github.com/golang-jwt/jwt"

func (c *Credential) ParseJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return c.SecretKey, nil
	})
}
