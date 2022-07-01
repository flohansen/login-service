package security

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/rds/rdsutils"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type HashEngine interface {
	HashPassword(password []byte) ([]byte, error)
}

type BcryptEngine struct {
}

type Credentials struct {
	Username string
	Password string
}

type CredentialsProvider interface {
	GetCredentials() (Credentials, error)
}

type AwsCredentialsProvider struct {
	host     string
	port     int
	username string
	region   string
}

func NewAwsCredentialsProvider(host string, port int, username string, region string) *AwsCredentialsProvider {
	return &AwsCredentialsProvider{
		host:     host,
		port:     port,
		username: username,
		region:   region,
	}
}

func (acp *AwsCredentialsProvider) GetCredentials() (Credentials, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Credentials{}, err
	}

	credentialsPath := path.Join(homeDir, ".aws", "credentials")
	creds := credentials.NewSharedCredentials(credentialsPath, "default")
	endpoint := fmt.Sprintf("%s:%d", acp.host, acp.port)
	authToken, err := rdsutils.BuildAuthToken(endpoint, acp.region, acp.username, creds)

	return Credentials{Username: acp.username, Password: authToken}, err
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
