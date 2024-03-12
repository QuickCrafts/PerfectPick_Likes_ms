package main

type Like struct {
	UserID int `json:"id"`
}

func NewLike(id int) *Like {
	return &Like{
		UserID: id,
	}
}
