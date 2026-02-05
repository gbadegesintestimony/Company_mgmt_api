package models

import "time"

type OTP struct {
	ID        string
	UserID    string
	CodeHash  string
	Purpose   string
	ExpiresAt time.Time
}
