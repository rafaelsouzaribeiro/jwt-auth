package jwtauth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Credential struct {
	SecretKey     []byte
	ExpireIn      int
	DeniedMethods []string
}

func NewCredential(Expire int, SecretKey string, DeniedMethods []string) (*Credential, error) {
	if Expire <= 0 {
		return nil, errors.New("expiration time must be greater than zero")
	}

	if SecretKey == "" {
		return nil, errors.New("the secret key cannot be empty")
	}

	c := &Credential{
		SecretKey:     []byte(SecretKey),
		ExpireIn:      Expire,
		DeniedMethods: DeniedMethods,
	}

	return c, nil
}

func (c *Credential) ExtractClaims(tokenStr string) (jwt.MapClaims, error) {
	token, err := c.ParseJWT(tokenStr)

	if err != nil {
		return nil, err
	}

	claims, _ := token.Claims.(jwt.MapClaims)

	return claims, nil
}

func (c *Credential) CreateToken(claims map[string]interface{}) (string, error) {
	// Ensure mandatory claims like exp are present
	if _, ok := claims["exp"]; !ok {
		claims["exp"] = time.Now().Add(time.Second * time.Duration(c.ExpireIn)).Unix() // Default expiration (1 hour)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))

	tokenString, err := token.SignedString(c.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Deprecated
// func (c *Credential) CreateToken(username string) (string, error) {
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
// 		jwt.MapClaims{
// 			"username": username,
// 			"exp":      time.Now().Add(time.Second * time.Duration(c.ExpireIn)).Unix(),
// 		})

// 	tokenString, err := token.SignedString(c.SecretKey)
// 	if err != nil {
// 		return "", err
// 	}

// 	return tokenString, nil
// }

func (c *Credential) ParseJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return c.SecretKey, nil
	})
}

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

func (c *Credential) TokenExpired(tokenString string) bool {
	token, _ := jwt.Parse(tokenString, nil)
	claims, _ := token.Claims.(jwt.MapClaims)
	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
	return expirationTime.Before(time.Now())
}

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

func (i *Credential) StreamInterceptorBearer(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

	// Extracts the service and method name
	_, methodName := extractServiceMethod(info.FullMethod)

	if contains(i.DeniedMethods, methodName) {
		return handler(srv, ss)
	}

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

	// Extracts the service and method name
	_, methodName := extractServiceMethod(info.FullMethod)

	if contains(i.DeniedMethods, methodName) {
		return handler(ctx, req)
	}

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

func (c *Credential) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", -1)
		err := c.VerifyToken(token)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Access denied: invalid token")
			return
		}

		next(w, r)
	}
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

func extractServiceMethod(fullMethod string) (string, string) {
	parts := strings.Split(fullMethod, "/")

	if len(parts) != 3 {
		return "", ""
	}

	// 1 - Service Name, 2 - Method Name
	return parts[1], parts[2]
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
