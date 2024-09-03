package middleware

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Credential) VerifyToken(tokenString string) error {

	token, err := c.ParseJWT(tokenString)

	if err != nil {
		return err
	}

	if !token.Valid {
		return status.Error(codes.Unauthenticated, "Invalid token")
	}

	return nil
}
