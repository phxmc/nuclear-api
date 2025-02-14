package dto

import "time"

type Pass struct {
	Id        string     `json:"id"`
	AccountId string     `json:"account_id"`
	From      *time.Time `json:"from"`
	To        *time.Time `json:"to"`
	Active    bool       `json:"active"`
	CreatedAt time.Time  `json:"created_at"`
}

type PassRequest struct {
	AccountId string     `json:"account_id"`
	From      *time.Time `json:"from"`
	To        *time.Time `json:"to"`
}
