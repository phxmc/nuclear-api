package domain

import "time"

type Pass struct {
	Id        string
	From      *time.Time
	To        *time.Time
	CreatedAt time.Time
}
