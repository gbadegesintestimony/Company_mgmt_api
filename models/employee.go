package models

import (
	"time"
)

type Employee struct {
	ID           string
	CompanyID    string
	Email        string
	FirstName    string
	LastName     string
	PasswordHash string
	Role         string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    string
}
