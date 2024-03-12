package main

type Storage interface {
	// Create
	CreateUser(*User) error
	CreateMedia(*Media) error
	CreateLike(*Like) error

	//Update
	UpdateLike(*UpdateLike) error

	// Get
	GetUserLikes(int) (*GetUserLikes, error)
	GetMediaLikes(*RequestMediaLikes) (*GetMediaLikes, error)
	GetAverage(*RequestMediaLikes) (*GetRating, error)
	GetWishlist(int) (*GetWishlist, error)

	//Delete
	DeleteUser(int) error
	DeleteMedia(int) error
	DeleteLike(int) error
}
