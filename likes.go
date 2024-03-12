package main

type Like struct {
	UserID int `json:"id"`
}

type UpdateLike struct {
}

type GetRating struct {
}

type GetWishlist struct {
}

func NewLike(id int) *Like {
	return &Like{
		UserID: id,
	}
}
