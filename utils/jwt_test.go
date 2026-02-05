package utils

import "testing"

func TestJWTGenerationAndValidation(t *testing.T) {
	secret := "test-secret"

	// 1. Use GenerateAccessToken (matching your jwt.go)
	// Added "admin" as the role argument required by your function
	token, err := GenerateAccessToken("u1", "c1", "admin", secret)
	if err != nil {
		t.Fatal(err)
	}

	// 2. Use ParseToken (matching your jwt.go)
	claims, err := ParseToken(token, secret)
	if err != nil {
		t.Fatal(err)
	}

	// 3. Access claims using the map key "user_id"
	if claims["user_id"] != "u1" {
		t.Fatal("invalid user id")
	}
}
