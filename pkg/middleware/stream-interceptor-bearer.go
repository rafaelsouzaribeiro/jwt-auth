package middleware

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Credential) StreamInterceptorBearer(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

	_, methodName := i.ExtractServiceMethod(info.FullMethod)

	if i.Contains(i.DeniedMethods, methodName) {
		return handler(srv, ss)
	}

	token, err := i.GetToken(ss.Context())

	if err != nil {
		return status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	err = i.VerifyToken(token)

	if err != nil {
		return status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	return handler(srv, ss)
}
