package main

type User struct {
	UserID int `json:"id"`
}

func NewUser(id int) *User {
	return &User{
		UserID: id,
	}
}
