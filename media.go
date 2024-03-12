package main

type Media struct {
	UserID int `json:"id"`
}

type GetMediaLikes struct {
}

type RequestMediaLikes struct {
}

func NewMedia(id int) *Media {
	return &Media{
		UserID: id,
	}
}
