package main

type LikeRelation struct {
	UserID    int `json:"user_id"`
	MediaID   any `json:"id"`
	MediaType any `json:"type"`      // 'MOV' | 'BOO' | 'SON'
	LikeType  any `json:"like_type"` // 'LK' | 'DLK' | 'BLK'
	Wishlist  any `json:"wishlist"`

	//Optional Attributes
	Rating any `json:"rating"`
}

type GetUserLikes struct {
	UserID int            `json:"id"`
	Movies []LikeRelation `json:"movies"`
	Songs  []LikeRelation `json:"songs"`
	Books  []LikeRelation `json:"books"`
}

func NewLikeRelation(id int, media any, mtype any, ltype any, w any, r any) *LikeRelation {
	return &LikeRelation{
		UserID:    id,
		MediaID:   media,
		MediaType: mtype,
		LikeType:  ltype,
		Wishlist:  w,
		Rating:    r,
	}
}
