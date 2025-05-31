package utils

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret []byte

func init() {
	// Try to get secret from environment variable first
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Generate a random secret if none provided
		secret = generateRandomSecret()
		log.Printf("Generated new JWT secret (consider setting JWT_SECRET environment variable)")
	}
	jwtSecret = []byte(secret)
}

func generateRandomSecret() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to a default secret if random generation fails
		return "goforum-default-secret-key-change-in-production-2024"
	}
	return hex.EncodeToString(bytes)
}

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID int, username, role string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "goforum",
			Subject:   username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})

	if err != nil {
		log.Printf("JWT validation error: %v", err)
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
