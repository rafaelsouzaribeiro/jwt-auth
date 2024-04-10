package middleware

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCredentialTime(t *testing.T) {
	c, err := NewCredential(0, "secret", nil)
	assert.Equal(t, errors.New("expiration time must be greater than zero"), err)
	assert.Nil(t, c)

}

func TestNewCredentialSecret(t *testing.T) {
	c, err := NewCredential(60, "", nil)
	assert.Equal(t, errors.New("the secret key cannot be empty"), err)
	assert.Nil(t, c)
}

func TestVeridyToken(t *testing.T) {
	c, err := NewCredential(60, "secret", nil)
	assert.Nil(t, err)

	claims := map[string]interface{}{
		"username": "Rafael",
	}

	token, err := c.CreateToken(claims)
	assert.Nil(t, err)
	assert.NotNil(t, token)

	err = c.VerifyToken("")
	assert.NotNil(t, err)

	err = c.VerifyToken(token)
	assert.Nil(t, err)

}

func TestCreate(t *testing.T) {
	c, err := NewCredential(60, "secret", nil)
	assert.Nil(t, err)

	claims := map[string]interface{}{
		"username": "Rafael",
	}
	token, err := c.CreateToken(claims)
	assert.Nil(t, err)
	assert.NotNil(t, token)

}

func TestTokenExpired(t *testing.T) {
	c, err := NewCredential(3, "secret", nil)
	assert.Nil(t, err)

	claims := map[string]interface{}{
		"username": "Rafael",
	}
	token, err := c.CreateToken(claims)
	assert.Nil(t, err)
	assert.NotNil(t, token)

	assert.False(t, c.TokenExpired(token))

	time.Sleep(4 * time.Second)
	assert.True(t, c.TokenExpired(token))

}

func TestExtracClaims(t *testing.T) {
	cre, err := NewCredential(1, "secretkey", nil)
	assert.Nil(t, err)
	assert.NotNil(t, cre)

	claims := map[string]interface{}{
		"lastname":  "Fernando",
		"firstname": "Rafael",
		// ... other claims
	}

	token, err := cre.CreateToken(claims)
	assert.Nil(t, err)
	assert.NotNil(t, token)

	claims, err = cre.ExtractClaims(token)
	assert.Nil(t, err)
	assert.NotNil(t, claims)

	assert.Equal(t, claims["lastname"], "Fernando")
	assert.Equal(t, claims["firstname"], "Rafael")

}
