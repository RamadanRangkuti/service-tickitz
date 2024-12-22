package pkg

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pilinux/argon2"
)

func VerifyHash(hash string, password string) bool {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	result, _ := argon2.ComparePasswordAndHash(hash, os.Getenv("SECRETKEY"), password)
	return result
}

func GenerateHash(password string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	result, _ := argon2.CreateHash(password, os.Getenv("SECRETKEY"), argon2.DefaultParams)
	return result
}
