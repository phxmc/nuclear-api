package dto

import "time"

type Nickname struct {
	Value     string    `json:"nickname"`
	CreatedAt time.Time `json:"created_at"`
}

type NicknameRequest struct {
	Value string `json:"nickname"`
}
