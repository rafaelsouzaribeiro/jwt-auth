package middleware

import "errors"

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
