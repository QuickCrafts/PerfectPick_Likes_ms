package main

type Media struct {
	UserID int `json:"id"`
}

func NewMedia(id int) *Media {
	return &Media{
		UserID: id,
	}
}
