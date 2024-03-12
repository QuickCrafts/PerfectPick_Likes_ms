package main

type Like struct {
	UserID    int    `json:"user_id"`
	MediaID   int    `json:"media_id"`
	MediaType string `json:"media_type"` // 'MOV' | 'BOO' | 'SON'
	LikeType  string `json:"like_type"`  // 'LK' | 'DLK' | 'BLK'
	Wishlist  bool   `json:"wishlist"`

	//Optional Attributes
	Rating float64 `json:"rating"`
}

type UpdateLike struct {
}

type GetRating struct {
	MediaID   int     `json:"id"`
	MediaType string  `json:"type"` // 'MOV' | 'SON' | 'BOO'
	AvgRating float64 `json:"avg_rating"`
}

type GetWishlist struct {
	UserID int     `json:"id"`
	Movies []int64 `json:"movies"`
	Songs  []int64 `json:"songs"`
	Books  []int64 `json:"books"`
}

func NewLike(id int, media int, mtype string, ltype string, w bool, r float64) *Like {
	return &Like{
		UserID:    id,
		MediaID:   media,
		MediaType: mtype,
		LikeType:  ltype,
		Wishlist:  w,
		Rating:    r,
	}
}
