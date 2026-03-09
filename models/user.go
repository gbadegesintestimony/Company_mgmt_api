package models

import "time"

type User struct {
	ID             string
	CompanyID      string
	Email          string
	PasswordHash   string
	FirstName      *string
	LastName       *string
	Phone          *string
	JobTitle       *string
	Department     *string
	IsActive       bool
	Role           string
	CreatedAt      time.Time
	UpdatedAt      *time.Time
	DeletedAt      *time.Time
	EmailVerifiedAt *time.Time
}