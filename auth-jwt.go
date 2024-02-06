package authjwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Credential struct {
	SecretKey []byte
	ExpireIn  int
}

func NewCredential(Expire int, SecretKey string) (*Credential, error) {
	if Expire <= 0 {
		return nil, errors.New("o tempo de expiração deve ser maior que zero")
	}

	if SecretKey == "" {
		return nil, errors.New("a chave secreta não pode estar vazia")
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
		return fmt.Errorf("invalid token")
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
