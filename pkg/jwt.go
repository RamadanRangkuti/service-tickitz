package pkg

import (
	"errors"
	"log"
	"os"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/joho/godotenv"
)

func secretKey() []byte {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	var JWT_SECRET []byte = []byte(Md5Hash(os.Getenv("JWTKEY")))
	return JWT_SECRET
}

func GenerateToken(userId int) (string, error) {
	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: secretKey()}, (&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		return "", err
	}
	payload := struct {
		UserId int `json:"userId"`
	}{
		UserId: userId,
	}
	token, err := jwt.Signed(sig).Claims(payload).Serialize()
	if err != nil {
		return "", nil
	}
	return token, nil
}

func VerifyToken(token string) (*jwt.Claims, error) {
	tok, err := jwt.ParseSigned(token, []jose.SignatureAlgorithm{jose.HS256})
	if err != nil {
		return nil, err
	}

	out := jwt.Claims{}
	err = tok.Claims(secretKey(), &out)
	if err != nil {
		return nil, errors.New("token is invalid or expired")
	}
	return &out, nil
}
