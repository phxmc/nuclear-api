package domain

import "time"

type Account struct {
	Id        string
	Email     string
	Perms     int
	CreatedAt time.Time
}
