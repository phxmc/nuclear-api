package domain

import "time"

type Pass struct {
	Id        string
	AccountId string
	From      *time.Time
	To        *time.Time
	Active    bool
	CreatedAt time.Time
}
