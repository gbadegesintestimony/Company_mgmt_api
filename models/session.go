package models

import "time"

type Session struct {
	ID               string
	UserID           string
	CompanyID        string
	Role             string
	RefreshTokenHash string
	ExpiresAt        time.Time
	RevokedAt        *time.Time
}
