package service

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenType defines whether it's access or refresh
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// MyCustomClaims defines what we encode inside the token
type MyCustomClaims struct {
	UserID string    `json:"user_id"`
	Type   TokenType `json:"type"`
	jwt.RegisteredClaims
}

// GenerateTokens creates both Access (15m) and Refresh (7d) tokens
func GenerateTokens(userID string) (string, string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", "", errors.New("JWT_SECRET not set in .env")
	}
	key := []byte(secret)

	// 1. Create Access Token (Short Lived - 15 mins)
	accessClaims := MyCustomClaims{
		UserID: userID,
		Type:   AccessToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "travel-app",
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(key)
	if err != nil {
		return "", "", err
	}

	// 2. Create Refresh Token (Long Lived - 7 days)
	refreshClaims := MyCustomClaims{
		UserID: userID,
		Type:   RefreshToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "travel-app",
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(key)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ValidateToken parses and validates a token string
func ValidateToken(tokenString string) (*MyCustomClaims, error) {
	secret := os.Getenv("JWT_SECRET")
	key := []byte(secret)

	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify strict signing method to prevent "None" alg attacks
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return key, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
