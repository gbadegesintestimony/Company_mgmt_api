package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

func GenerateOTP() string {
	// crypto/rand, not math/rand — this code gates email verification and
	// password resets, so it must not be predictable.
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return fmt.Sprintf("%06d", n.Int64())
}

func HashOTP(code string) string {
	// Implementation for hashing the OTP code
	hash, _ := bcrypt.GenerateFromPassword([]byte(code), bcryptCost)

	return string(hash)
}

func VerifyOTP(hash, code string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(code))
}
