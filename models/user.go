package models

type User struct {
	ID           string
	CompanyID    string
	Email        string
	PasswordHash string
	Role         string
	IsActive     bool
	CreatedAt    string
	UpdatedAt    string
	// DeletedAt *string
}
