package middleware

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func GetToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return "", status.Errorf(codes.Unauthenticated, "missing authorization header")
	}

	parts := strings.Split(authHeader[0], " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", status.Errorf(codes.Unauthenticated, "invalid authorization format")
	}

	return parts[1], nil

}
