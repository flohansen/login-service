package security

import (
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type HashEngine interface {
	HashPassword(password []byte) ([]byte, error)
}

type BcryptEngine struct {
}

type JwtConfig struct {
	SignKey string
}

type JwtClaims struct {
	UserId   int    `json:"userId"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func NewBcryptEngine() *BcryptEngine {
	return &BcryptEngine{}
}

func GenerateToken(id int, username string, signingMethod jwt.SigningMethod, key interface{}) (string, error) {
	claims := JwtClaims{
		UserId:   id,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(signingMethod, claims)
	signedToken, err := token.SignedString(key)
	return signedToken, err
}

func (b *BcryptEngine) HashPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, 8)
}
