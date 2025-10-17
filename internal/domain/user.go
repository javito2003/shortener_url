package domain

import "time"

type User struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Password  string
	CreatedAt time.Time
}
