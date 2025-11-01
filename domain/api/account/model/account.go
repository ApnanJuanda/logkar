package model

type Account struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	IsSeller bool   `json:"is_seller"`
}
