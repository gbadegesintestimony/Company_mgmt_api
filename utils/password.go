package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

func HashPassword(password string) (string, error) {
	// Implementation for hashing password

	if len(password) < 8 {
		return "", errors.New("password must be atleast 8 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return string(hash), err
	}
	return string(hash), nil
}

func VerifyPassword(password, hash string) error {
	// Implementation for verifying password
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
}
