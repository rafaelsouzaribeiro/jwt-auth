package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Credential) UnaryInterceptorBearer(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

	_, methodName := i.ExtractServiceMethod(info.FullMethod)

	if i.Contains(i.DeniedMethods, methodName) {
		return handler(ctx, req)
	}

	token, err := i.GetToken(ctx)

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	err = i.VerifyToken(token)

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	return handler(ctx, req)
}
