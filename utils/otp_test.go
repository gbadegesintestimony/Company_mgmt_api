package utils

import "testing"

func TestOTPGeneration(t *testing.T) {
	otp := GenerateOTP()
	if len(otp) != 6 {
		t.Fatal("OTP must be 6 digits")
	}
}
