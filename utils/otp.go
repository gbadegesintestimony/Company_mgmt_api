package utils

import (
	"fmt"
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

func GenerateOTP() string {
	// Implementation for generating a random OTP code
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func HashOTP(code string) string {
	// Implementation for hashing the OTP code
	hash, _ := bcrypt.GenerateFromPassword([]byte(code), bcryptCost)

	return string(hash)
}

func VerifyOTP(hash, code string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(code))
}
