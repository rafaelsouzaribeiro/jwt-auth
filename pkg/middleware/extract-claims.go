package middleware

import "github.com/golang-jwt/jwt"

func (c *Credential) ExtractClaims(tokenStr string) (jwt.MapClaims, error) {
	token, err := c.ParseJWT(tokenStr)

	if err != nil {
		return nil, err
	}

	claims, _ := token.Claims.(jwt.MapClaims)

	return claims, nil
}
