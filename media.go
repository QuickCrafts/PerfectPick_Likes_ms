package main

type GetMediaLikes struct {
	Likes     []LikeRelation `json:"likes"`
	AvgRating float64        `json:"avg_rating"`
}
