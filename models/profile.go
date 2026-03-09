package models

type Profile struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Phone      string `json:"phone"`
	JobTitle   string `json:"job_title"`
	Department string `json:"department"`
}
