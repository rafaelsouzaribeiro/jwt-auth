package authjwt

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCredentialTime(t *testing.T) {
	c, err := NewCredential(0, "secret")
	assert.Equal(t, errors.New("o tempo de expiração deve ser maior que zero"), err)
	assert.Nil(t, c)

}

func TestNewCredentialSecret(t *testing.T) {
	c, err := NewCredential(60, "")
	assert.Equal(t, errors.New("a chave secreta não pode estar vazia"), err)
	assert.Nil(t, c)
}

func TestVeridyToken(t *testing.T) {
	c, err := NewCredential(60, "secret")
	assert.Nil(t, err)

	token, err := c.CreateToken("username")
	assert.Nil(t, err)
	assert.NotNil(t, token)

	err = c.VerifyToken("")
	assert.NotNil(t, err)

	err = c.VerifyToken(token)
	assert.Nil(t, err)

}

func TestCreate(t *testing.T) {
	c, err := NewCredential(60, "secret")
	assert.Nil(t, err)

	token, err := c.CreateToken("username")
	assert.Nil(t, err)
	assert.NotNil(t, token)

}

func TestTokenExpired(t *testing.T) {
	c, err := NewCredential(3, "secret")
	assert.Nil(t, err)

	token, err := c.CreateToken("username")
	assert.Nil(t, err)
	assert.NotNil(t, token)

	assert.False(t, c.TokenExpired(token))

	time.Sleep(4 * time.Second)
	assert.True(t, c.TokenExpired(token))

}
