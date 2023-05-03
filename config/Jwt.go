package config

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var JWT_KEY = []byte(os.Getenv("JWT_KEY"))

type JTWClaim struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// CreateToken is a function to create token
func CreateToken(username string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// exp := jwt.NumericDate(time.Now().Add(time.Hour * 24).Unix())
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		Issuer:    "alterra",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWT_KEY)
}

func VerifyToken(tokenString string) (*JTWClaim, error) {
	// Initialize a new instance of `Claims`
	claims := &JTWClaim{}

	// Parse the JWT token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token's signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return JWT_KEY, nil
	})

	// If there is an error validating the token, return the error
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Return the claims
	return claims, nil
}

func EncryptToken(token string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	key := []byte(os.Getenv("AES_KEY"))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// PKCS#7 padding
	padSize := aes.BlockSize - len(token)%aes.BlockSize
	padding := bytes.Repeat([]byte{byte(padSize)}, padSize)
	paddedToken := append([]byte(token), padding...)

	// CBC mode
	ciphertext := make([]byte, aes.BlockSize+len(paddedToken))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(ciphertext[aes.BlockSize:], paddedToken)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func DecryptToken(encoded string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	key := []byte(os.Getenv("AES_KEY"))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	// CBC mode
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(ciphertext, ciphertext)

	// PKCS#7 unpadding
	padSize := int(ciphertext[len(ciphertext)-1])
	if padSize < 1 || padSize > aes.BlockSize {
		return "", errors.New("invalid padding size")
	}

	return string(ciphertext[:len(ciphertext)-padSize]), nil
}
