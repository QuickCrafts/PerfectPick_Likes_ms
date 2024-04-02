package main

type LikeRelation struct {
	UserID    any `json:"user_id"`
	MediaID   any `json:"media_id"`
	MediaType any `json:"type"`      // 'MOV' | 'BOO' | 'SON'
	LikeType  any `json:"like_type"` // 'LK' | 'DLK'
}

type GetUserLikes struct {
	UserID int            `json:"id"`
	Movies []LikeRelation `json:"movies"`
	Songs  []LikeRelation `json:"songs"`
	Books  []LikeRelation `json:"books"`
}

func NewLikeRelation(id any, media any, mtype any, ltype any) *LikeRelation {
	return &LikeRelation{
		UserID:    id,
		MediaID:   media,
		MediaType: mtype,
		LikeType:  ltype,
	}
}
