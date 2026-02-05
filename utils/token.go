package utils

import "golang.org/x/crypto/bcrypt"

const refreshTokenCost = 12

func HashRefreshToken(token string) (string, error) {
	// Implementation for hashing refresh token
	hash, err := bcrypt.GenerateFromPassword([]byte(token), refreshTokenCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func VerifyRefreshToken(token, hash string) error {
	// Implementation for verifying refresh token
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(token),
	)
}
