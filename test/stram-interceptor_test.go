package test

import (
	"fmt"
	"net"
	"testing"

	"github.com/rafaelsouzaribeiro/jwt-auth/pkg/middleware"
	"google.golang.org/grpc"
)

func TestStreamInterceptorBearer(t *testing.T) {
	errChan := make(chan error)

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(StreamInterceptor(t, errChan)),
	)

	lis, err := net.Listen("tcp", "localhost:30015")

	if err != nil {
		panic(err)
	}

	go func() {
		err := grpcServer.Serve(lis)
		errChan <- err
	}()

	err = <-errChan
	if err != nil {
		t.Errorf("Error during server serving: %v", err)
	}

}

func StreamInterceptor(t *testing.T, errChan chan error) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		cre, err := middleware.NewCredential(60, "secret", nil)

		if err != nil {
			errChan <- err
		}

		if cre == nil {
			errChan <- fmt.Errorf("error creating credentials")
		}

		claims := map[string]interface{}{
			"username": "Rafael",
		}
		token, err := cre.CreateToken(claims)

		if err != nil {
			errChan <- err
		}

		if token == "" {
			errChan <- fmt.Errorf("error creating token")
		}

		err = cre.VerifyToken(token)

		if err != nil {
			errChan <- err
		}

		return handler(srv, ss)
	}
}
