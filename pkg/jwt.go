package pkg

import (
	"errors"
	"log"
	"os"
	"time"

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

type BaseInfo struct {
	IssuedAt jwt.NumericDate `json:"iat"` // Waktu token diterbitkan
}

type TokenPayload struct {
	UserId   int `json:"userId"`   // ID pengguna
	UserRole int `json:"userRole"` // Role pengguna
}

type FullToken struct {
	BaseInfo
	TokenPayload
}

func GenerateToken(userId int, userRole int) (string, error) {
	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: secretKey()}, (&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		return "", err
	}
	tokenData := FullToken{
		BaseInfo: BaseInfo{
			IssuedAt: *jwt.NewNumericDate(time.Now()),
		},
		TokenPayload: TokenPayload{
			UserId:   userId,
			UserRole: userRole,
		},
	}
	token, err := jwt.Signed(sig).Claims(tokenData).Serialize()
	if err != nil {
		return "", nil
	}
	return token, nil
}

func VerifyToken(token string) (*TokenPayload, error) {
	tok, err := jwt.ParseSigned(token, []jose.SignatureAlgorithm{jose.HS256})
	if err != nil {
		return nil, err
	}

	claims := &TokenPayload{}
	err = tok.Claims(secretKey(), &claims)
	if err != nil {
		return nil, errors.New("token is invalid or expired")
	}
	return claims, nil
}
