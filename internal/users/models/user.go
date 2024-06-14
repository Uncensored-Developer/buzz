package models

import "time"

type User struct {
	ID        int
	Email     string
	Password  string
	Name      string
	Gender    string
	Dob       time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
