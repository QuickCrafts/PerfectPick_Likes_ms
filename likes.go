package main

type Like struct {
	UserID    int    `json:"user_id"`
	MediaID   string `json:"media_id"`
	MediaType string `json:"media_type"` // 'MOV' | 'BOO' | 'SON'
	LikeType  string `json:"like_type"`  // 'LK' | 'DLK'
}

type LikeExtended struct {
	UserID    int    `json:"user_id"`
	MediaID   string `json:"media_id"`
	MediaType string `json:"media_type"` // 'MOV' | 'BOO' | 'SON'
	LikeType  string `json:"like_type"`  // 'LK' | 'DLK'
	Wishlist  bool   `json:"wishlist"`

	//Optional Attributes
	Rating float64 `json:"rating"`
}

type Rate struct {
	Rating float64 `json:"rating"`
}

type ChangeWishlist struct {
	MediaID   string `json:"media_id"`
	MediaType string `json:"media_type"` // 'MOV' | 'BOO' | 'SON'
	Type      string `json:"type"`       // 'RMV' | 'ADD'
}

type GetWishlist struct {
	UserID int      `json:"id"`
	Movies []string `json:"movies"`
	Songs  []string `json:"songs"`
	Books  []string `json:"books"`
}

func NewLike(id int, media string, mtype string, ltype string) *Like {
	return &Like{
		UserID:    id,
		MediaID:   media,
		MediaType: mtype,
		LikeType:  ltype,
	}
}
