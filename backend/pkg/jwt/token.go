package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims represents the JWT claims
type Claims struct {
	UserID uuid.UUID `json:"userId"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

// TokenManager handles JWT token operations
type TokenManager struct {
	secretKey []byte
	issuer    string
}

// NewTokenManager creates a new TokenManager
func NewTokenManager(secretKey, issuer string) *TokenManager {
	return &TokenManager{
		secretKey: []byte(secretKey),
		issuer:    issuer,
	}
}

// GenerateToken generates a new JWT token
func (tm *TokenManager) GenerateToken(userID uuid.UUID, email, role string, duration time.Duration) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    tm.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(tm.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (tm *TokenManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return tm.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
