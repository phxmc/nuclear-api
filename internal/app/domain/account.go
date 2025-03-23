package domain

import "time"

type Account struct {
	Id         string
	Email      string
	TelegramId *int64
	Perms      int
	CreatedAt  time.Time
}
