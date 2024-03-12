package main

type User struct {
	UserID int `json:"id"`
}

type GetUserLikes struct {
}

func NewUser(id int) *User {
	return &User{
		UserID: id,
	}
}
