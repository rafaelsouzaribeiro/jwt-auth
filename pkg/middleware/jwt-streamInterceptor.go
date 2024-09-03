package middleware

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Credential) JwtStreamInterceptor(token string) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		err := c.VerifyToken(token)

		if err != nil {
			return status.Error(codes.Unauthenticated, "Invalid token")
		}

		if c.TokenExpired(token) {
			return status.Error(codes.Unauthenticated, "Token Expired")
		}

		return handler(srv, ss)
	}
}
