package models

import "time"

type Company struct {
	ID        string
	Name      string
	Domain    string
	Status    string
	CreatedAt time.Time
}
