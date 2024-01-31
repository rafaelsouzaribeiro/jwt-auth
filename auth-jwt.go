package authjwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
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
			"exp":      time.Now().Add(time.Hour * time.Duration(c.ExpireIn)).Unix(),
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
