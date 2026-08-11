package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(userID, companyID, role, secret string) (string, error) {
	// Implementation for generating JWT access token
	claims := jwt.MapClaims{
		"user_id":    userID,
		"company_id": companyID,
		"role":       role,
		"exp":        time.Now().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GenerateRefreshToken returns a cryptographically random, high-entropy
// opaque token. It carries no claims — validity is enforced entirely by the
// hashed lookup in SessionRepository, so it must never be predictable or
// reproducible across calls (unlike a JWT built from a coarse timestamp).
func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func ParseToken(tokenStr, secret string) (jwt.MapClaims, error) {
	// Implementation for parsing and validating JWT token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}
