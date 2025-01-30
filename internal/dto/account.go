package dto

type Account struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Perms int    `json:"perms"`
}
