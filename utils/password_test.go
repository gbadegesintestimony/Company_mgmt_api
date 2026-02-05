package utils

import "testing"

// In password_test.go
func TestPasswordHashAndVerify(t *testing.T) {
	hash, err := HashPassword("password123")
	if err != nil {
		t.Fatal(err)
	}

	// Change CheckPassword to VerifyPassword
	// Since VerifyPassword returns an error, we check if err is nil
	err = VerifyPassword("password123", hash)
	if err != nil {
		t.Fatal("password should match but got error:", err)
	}
}
