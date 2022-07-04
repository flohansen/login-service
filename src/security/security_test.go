package security

import (
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestGenerateToken(t *testing.T) {
	tokenString, err := GenerateToken(0, "test", jwt.SigningMethodHS256, []byte("supersecretsignkey"))
	if err != nil {
		t.Fatal(err)
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte("supersecretsignkey"), nil
	})

	claims := token.Claims.(jwt.MapClaims)

	assert.Nil(t, err)
	assert.NotNil(t, claims["userId"])
	assert.NotNil(t, claims["username"])
	assert.Equal(t, 0.0, claims["userId"])
	assert.Equal(t, "test", claims["username"])
}

func TestBcryptEngineHashPassword(t *testing.T) {
	engine := NewBcryptEngine()
	hashedPassword, _ := engine.HashPassword([]byte("some password"))
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte("some password"))

	assert.NoError(t, err)
}
