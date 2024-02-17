package authjwt

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Credential struct {
	SecretKey []byte
	ExpireIn  int
}

func NewCredential(Expire int, SecretKey string) (*Credential, error) {
	if Expire <= 0 {
		return nil, errors.New("expiration time must be greater than zero")
	}

	if SecretKey == "" {
		return nil, errors.New("the secret key cannot be empty")
	}

	c := &Credential{
		SecretKey: []byte(SecretKey),
		ExpireIn:  Expire,
	}

	return c, nil
}

func (c *Credential) CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Second * time.Duration(c.ExpireIn)).Unix(),
		})

	tokenString, err := token.SignedString(c.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (c *Credential) VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return c.SecretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return status.Error(codes.Unauthenticated, "Token inválido")
	}

	return nil
}

func (c *Credential) TokenExpired(tokenString string) bool {
	token, _ := jwt.Parse(tokenString, nil)
	claims, _ := token.Claims.(jwt.MapClaims)
	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
	return expirationTime.Before(time.Now())
}

// jwtStreamInterceptor é o interceptor gRPC para autenticação JWT em streams
func (c *Credential) JwtStreamInterceptor(token string) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		err := c.VerifyToken(token)

		if err != nil {
			return status.Error(codes.Unauthenticated, "Token inválido")
		}

		if c.TokenExpired(token) {
			return status.Error(codes.Unauthenticated, "Token Expirou")
		}

		// Se não houver erro, chame o manipulador de chamada de streaming
		return handler(srv, ss)
	}
}

// JwtUnaryInterceptor é o interceptor gRPC para autenticação JWT em unary
func (c *Credential) JwtUnaryInterceptor(token string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		err := c.VerifyToken(token)

		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "Token inválido: "+err.Error())
		}

		if c.TokenExpired(token) {
			return nil, status.Error(codes.Unauthenticated, "Token Expirou")
		}

		// If no errors, proceed with the unary handler
		return handler(ctx, req)
	}
}

func (i *Credential) StreamInterceptorBearer(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	token, err := GetToken(ss.Context())

	if err != nil {
		return status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	err = i.VerifyToken(token)

	if err != nil {
		return status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	return handler(srv, ss)
}

func (i *Credential) UnaryInterceptorBearer(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

	token, err := GetToken(ctx)

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	err = i.VerifyToken(token)

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	// Continue with the handler if authentication is successful
	return handler(ctx, req)
}

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
