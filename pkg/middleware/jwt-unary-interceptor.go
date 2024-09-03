package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Credential) JwtUnaryInterceptor(token string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		err := c.VerifyToken(token)

		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "Invalid token: "+err.Error())
		}

		if c.TokenExpired(token) {
			return nil, status.Error(codes.Unauthenticated, "Token Expired")
		}

		// If no errors, proceed with the unary handler
		return handler(ctx, req)
	}
}
