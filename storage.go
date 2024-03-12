package main

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Storage interface {
	// Create
	CreateUser(int) error
	CreateMedia(int, string) error
	CreateLike(*Like) error

	//Update
	UpdateLike(*UpdateLike) error

	// Get
	GetUserLikes(int) (*GetUserLikes, error)
	GetMediaLikes(int, string) (*GetMediaLikes, error)
	GetAverage(int, string) (*GetRating, error)
	GetWishlist(int) (*GetWishlist, error)

	//Delete
	DeleteUser(int) error
	DeleteMedia(int, string) error
	DeleteLike(int, int, string) error

	//Close Session
	CloseSession()
}

type Neo4jStore struct {
	session neo4j.SessionWithContext
	ctx     context.Context
	driver  neo4j.DriverWithContext
}

func NewNeo4jStore() (*Neo4jStore, error) {
	ctx := context.Background()
	dbUri := "neo4j://localhost:7000"
	dbUser := "neo4j"
	dbPassword := "0900pass"
	driver, _ := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))

	err := driver.VerifyConnectivity(ctx)
	if err != nil {
		return nil, err
	}

	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	return &Neo4jStore{
		session: session,
		ctx:     ctx,
		driver:  driver,
	}, nil
}

func (s *Neo4jStore) CloseSession() {
	s.driver.Close(s.ctx)
	s.session.Close(s.ctx)
}

// Create Functions
func (s *Neo4jStore) CreateUser(i int) error {

	_, err := s.session.ExecuteWrite(s.ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(s.ctx, "CREATE (u:User {id_user: $id})", map[string]interface{}{"id": i})
		if err != nil {
			return nil, err
		}

		if result.Next(s.ctx) {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *Neo4jStore) CreateMedia(i int, tp string) error {

	query := "CREATE (m:Movie {id_movie: $id})"

	if tp == "SON" {
		query = "CREATE (s:Song {id_song: $id})"
	} else if tp == "BOO" {
		query = "CREATE (b:Book {id_book: $id})"
	}

	_, err := s.session.ExecuteWrite(s.ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(s.ctx, query, map[string]interface{}{"id": i})
		if err != nil {
			return nil, err
		}

		if result.Next(s.ctx) {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *Neo4jStore) CreateLike(l *Like) error {

	query := "MATCH (m:Movie), (u:User) WHERE m.id_movie = $id_media AND u.id_user = $id_user CREATE (u)-[r:LK {rating: $r, wishlist: $w, like_type:'LK', media_type: 'MOV', media_id: $id_media, user_id: $id_user}]->(m)"

	if l.LikeType == "DLK" {
		if l.MediaType == "SON" {
			query = "MATCH (s:Song), (u:User) WHERE s.id_song = $id_media AND u.id_user = $id_user CREATE (u)-[r:DLK {rating: $r, wishlist: $w, like_type:'DLK', media_type: 'SON', media_id: $id_media, user_id: $id_user}]->(s)"
		} else if l.MediaType == "BOO" {
			query = "MATCH (b:Book), (u:User) WHERE b.id_book = $id_media AND u.id_user = $id_user CREATE (u)-[r:DLK {rating: $r, wishlist: $w, like_type:'DLK', media_type: 'BOO', media_id: $id_media, user_id: $id_user}]->(b)"
		} else {
			query = "MATCH (m:Movie), (u:User) WHERE m.id_movie = $id_media AND u.id_user = $id_user CREATE (u)-[r:DLK {rating: $r, wishlist: $w, like_type:'DLK', media_type: 'MOV', media_id: $id_media, user_id: $id_user}]->(m)"
		}
	} else if l.LikeType == "BLK" {
		if l.MediaType == "SON" {
			query = "MATCH (s:Song), (u:User) WHERE s.id_song = $id_media AND u.id_user = $id_user CREATE (u)-[r:BLK {rating: $r, wishlist: $w, like_type:'BLK', media_type: 'SON', media_id: $id_media, user_id: $id_user}]->(s)"
		} else if l.MediaType == "BOO" {
			query = "MATCH (b:Book), (u:User) WHERE b.id_book = $id_media AND u.id_user = $id_user CREATE (u)-[r:BLK {rating: $r, wishlist: $w, like_type:'BLK', media_type: 'BOO', media_id: $id_media, user_id: $id_user}]->(b)"
		} else {
			query = "MATCH (m:Movie), (u:User) WHERE m.id_movie = $id_media AND u.id_user = $id_user CREATE (u)-[r:BLK {rating: $r, wishlist: $w, like_type:'BLK', media_type: 'MOV', media_id: $id_media, user_id: $id_user}]->(m)"
		}
	} else {
		if l.MediaType == "SON" {
			query = "MATCH (s:Song), (u:User) WHERE s.id_song = $id_media AND u.id_user = $id_user CREATE (u)-[r:LK {rating: $r, wishlist: $w, like_type:'LK', media_type: 'SON', media_id: $id_media, user_id: $id_user}]->(s)"
		} else if l.MediaType == "BOO" {
			query = "MATCH (b:Book), (u:User) WHERE b.id_book = $id_media AND u.id_user = $id_user CREATE (u)-[r:LK {rating: $r, wishlist: $w, like_type:'LK', media_type: 'BOO', media_id: $id_media, user_id: $id_user}]->(b)"
		}
	}

	_, err := s.session.ExecuteWrite(s.ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(s.ctx, query, map[string]interface{}{"id_media": l.MediaID, "id_user": l.UserID, "r": l.Rating, "w": l.Wishlist})
		if err != nil {
			return nil, err
		}

		if result.Next(s.ctx) {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return err
	}

	return nil
}

// Update Functions
func (s *Neo4jStore) UpdateLike(*UpdateLike) error {
	//@todo
	//queryLK := "MATCH (:User {id_user: $id_user})-[r]-(:Movie {id_movie: $id_media}) RETURN r"
	return nil
}

// Delete Functions
func (s *Neo4jStore) DeleteUser(i int) error {
	_, err := s.session.ExecuteWrite(s.ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(s.ctx, "MATCH (u:User) WHERE u.id_user = $id DELETE u", map[string]interface{}{"id": i})
		if err != nil {
			return nil, err
		}

		if result.Next(s.ctx) {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *Neo4jStore) DeleteMedia(i int, tp string) error {
	query := "MATCH (m:Movie) WHERE m.id_movie = $id DELETE m"

	if tp == "SON" {
		query = "MATCH (s:Song) WHERE s.id_song = $id DELETE s"
	} else if tp == "BOO" {
		query = "MATCH (b:Book) WHERE b.id_book = $id DELETE b"
	}

	_, err := s.session.ExecuteWrite(s.ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(s.ctx, query, map[string]interface{}{"id": i})
		if err != nil {
			return nil, err
		}

		if result.Next(s.ctx) {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *Neo4jStore) DeleteLike(user_id int, media_id int, tp string) error {
	queryLK := "MATCH (:Movie {id_movie: $id_media})-[r:LK]-(:User {id_user: $id_user}) DELETE r"

	if tp == "SON" {
		queryLK = "MATCH (:Song {id_song: $id_media})-[r:LK]-(:User {id_user: $id_user}) DELETE r"
	} else if tp == "BOO" {
		queryLK = "MATCH (:Book {id_book: $id_media})-[r:LK]-(:User {id_user: $id_user}) DELETE r"
	}

	queryDLK := "MATCH (:Movie {id_movie: $id_media})-[r:DLK]-(:User {id_user: $id_user}) DELETE r"

	if tp == "SON" {
		queryDLK = "MATCH (:Song {id_song: $id_media})-[r:DLK]-(:User {id_user: $id_user}) DELETE r"
	} else if tp == "BOO" {
		queryDLK = "MATCH (:Book {id_book: $id_media})-[r:DLK]-(:User {id_user: $id_user}) DELETE r"
	}

	queryBLK := "MATCH (:Movie {id_movie: $id_media})-[r:BLK]-(:User {id_user: $id_user}) DELETE r"

	if tp == "SON" {
		queryBLK = "MATCH (:Song {id_song: $id_media})-[r:BLK]-(:User {id_user: $id_user}) DELETE r"
	} else if tp == "BOO" {
		queryBLK = "MATCH (:Book {id_book: $id_media})-[r:BLK]-(:User {id_user: $id_user}) DELETE r"
	}

	_, err := s.session.ExecuteWrite(s.ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(s.ctx, queryLK, map[string]interface{}{"id_media": media_id, "id_user": user_id})
		if err != nil {
			return nil, err
		}

		if result.Next(s.ctx) {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return err
	}

	_, err = s.session.ExecuteWrite(s.ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(s.ctx, queryDLK, map[string]interface{}{"id_media": media_id, "id_user": user_id})
		if err != nil {
			return nil, err
		}

		if result.Next(s.ctx) {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return err
	}

	_, err = s.session.ExecuteWrite(s.ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(s.ctx, queryBLK, map[string]interface{}{"id_media": media_id, "id_user": user_id})
		if err != nil {
			return nil, err
		}

		if result.Next(s.ctx) {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return err
	}

	return nil
}

// Get Functions
func (s *Neo4jStore) GetUserLikes(i int) (*GetUserLikes, error) {
	queryLK := "MATCH (:User {id_user: $id_user})-[r:LK]-(n) RETURN r as relation"
	queryDLK := "MATCH (:User {id_user: $id_user})-[r:DLK]-(n) RETURN r as relation"
	queryBLK := "MATCH (:User {id_user: $id_user})-[r:BLK]-(n) RETURN r as relation"

	var results []neo4j.Relationship
	var movies []LikeRelation
	var songs []LikeRelation
	var books []LikeRelation

	_, err := s.session.ExecuteWrite(s.ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(s.ctx, queryDLK, map[string]interface{}{"id_user": i})
		if err != nil {
			return nil, err
		}

		for result.Next(s.ctx) {
			results = append(results, result.Record().AsMap()["relation"].(neo4j.Relationship))
		}

		return nil, result.Err()
	})

	if err != nil {
		return nil, err
	}

	_, errLK := s.session.ExecuteWrite(s.ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, errLK := transaction.Run(s.ctx, queryLK, map[string]interface{}{"id_user": i})
		if errLK != nil {
			return nil, errLK
		}

		for result.Next(s.ctx) {
			results = append(results, result.Record().AsMap()["relation"].(neo4j.Relationship))
		}

		return nil, result.Err()
	})

	if errLK != nil {
		return nil, errLK
	}

	_, errBLK := s.session.ExecuteWrite(s.ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, errBLK := transaction.Run(s.ctx, queryBLK, map[string]interface{}{"id_user": i})
		if errBLK != nil {
			return nil, errBLK
		}

		for result.Next(s.ctx) {
			results = append(results, result.Record().AsMap()["relation"].(neo4j.Relationship))
		}

		return nil, result.Err()
	})

	if errBLK != nil {
		return nil, errBLK
	}

	for r := 0; r < len(results); r++ {
		props := results[r].Props
		mediaType := props["media_type"]
		like := NewLikeRelation(i, props["media_id"], mediaType, props["like_type"], props["wishlist"], props["rating"])

		if mediaType == "MOV" {
			movies = append(movies, *like)
		} else if mediaType == "SON" {
			songs = append(songs, *like)
		} else {
			books = append(books, *like)
		}
	}

	return &GetUserLikes{
		UserID: i,
		Movies: movies,
		Songs:  songs,
		Books:  books,
	}, nil
}

func (s *Neo4jStore) GetMediaLikes(i int, tp string) (*GetMediaLikes, error) {
	queryLK := "MATCH (:Movie {id_movie: $id})-[r:LK]-(n) RETURN r as relation"
	queryDLK := "MATCH (:Movie {id_movie: $id})-[r:DLK]-(n) RETURN r as relation"
	queryBLK := "MATCH (:Movie {id_movie: $id})-[r:BLK]-(n) RETURN r as relation"

	if tp == "SON" {
		queryLK = "MATCH (:Song {id_song: $id})-[r:LK]-(n) RETURN r as relation"
		queryDLK = "MATCH (:Song {id_song: $id})-[r:DLK]-(n) RETURN r as relation"
		queryBLK = "MATCH (:Song {id_song: $id})-[r:BLK]-(n) RETURN r as relation"
	} else if tp == "BOO" {
		queryLK = "MATCH (:Book {id_book: $id})-[r:LK]-(n) RETURN r as relation"
		queryDLK = "MATCH (:Book {id_book: $id})-[r:DLK]-(n) RETURN r as relation"
		queryBLK = "MATCH (:Book {id_book: $id})-[r:BLK]-(n) RETURN r as relation"
	}

	var results []neo4j.Relationship
	var likes []LikeRelation
	var sumRating float64

	_, err := s.session.ExecuteWrite(s.ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(s.ctx, queryDLK, map[string]interface{}{"id": i})
		if err != nil {
			return nil, err
		}

		for result.Next(s.ctx) {
			results = append(results, result.Record().AsMap()["relation"].(neo4j.Relationship))
		}

		return nil, result.Err()
	})

	if err != nil {
		return nil, err
	}

	_, errLK := s.session.ExecuteWrite(s.ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, errLK := transaction.Run(s.ctx, queryLK, map[string]interface{}{"id": i})
		if errLK != nil {
			return nil, errLK
		}

		for result.Next(s.ctx) {
			results = append(results, result.Record().AsMap()["relation"].(neo4j.Relationship))
		}

		return nil, result.Err()
	})

	if errLK != nil {
		return nil, errLK
	}

	_, errBLK := s.session.ExecuteWrite(s.ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, errBLK := transaction.Run(s.ctx, queryBLK, map[string]interface{}{"id": i})
		if errBLK != nil {
			return nil, errBLK
		}

		for result.Next(s.ctx) {
			results = append(results, result.Record().AsMap()["relation"].(neo4j.Relationship))
		}

		return nil, result.Err()
	})

	if errBLK != nil {
		return nil, errBLK
	}

	sumRating = 0.0
	for r := 0; r < len(results); r++ {
		props := results[r].Props
		mediaType := props["media_type"]
		rating := props["rating"]
		like := NewLikeRelation(i, props["media_id"], mediaType, props["like_type"], props["wishlist"], rating)
		likes = append(likes, *like)

		if rating != -1 {
			sumRating = sumRating + rating.(float64)
		}
	}

	return &GetMediaLikes{
		Likes:     likes,
		AvgRating: sumRating / float64(len(results)),
	}, nil
}

func (s *Neo4jStore) GetAverage(i int, tp string) (*GetRating, error) {
	queryLK := "MATCH (:Movie {id_movie: $id})-[r:LK]-(n) RETURN r as relation"
	queryDLK := "MATCH (:Movie {id_movie: $id})-[r:DLK]-(n) RETURN r as relation"
	queryBLK := "MATCH (:Movie {id_movie: $id})-[r:BLK]-(n) RETURN r as relation"

	if tp == "SON" {
		queryLK = "MATCH (:Song {id_song: $id})-[r:LK]-(n) RETURN r as relation"
		queryDLK = "MATCH (:Song {id_song: $id})-[r:DLK]-(n) RETURN r as relation"
		queryBLK = "MATCH (:Song {id_song: $id})-[r:BLK]-(n) RETURN r as relation"
	} else if tp == "BOO" {
		queryLK = "MATCH (:Book {id_book: $id})-[r:LK]-(n) RETURN r as relation"
		queryDLK = "MATCH (:Book {id_book: $id})-[r:DLK]-(n) RETURN r as relation"
		queryBLK = "MATCH (:Book {id_book: $id})-[r:BLK]-(n) RETURN r as relation"
	}

	var results []neo4j.Relationship
	var sumRating float64

	_, err := s.session.ExecuteWrite(s.ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(s.ctx, queryDLK, map[string]interface{}{"id": i})
		if err != nil {
			return nil, err
		}

		for result.Next(s.ctx) {
			results = append(results, result.Record().AsMap()["relation"].(neo4j.Relationship))
		}

		return nil, result.Err()
	})

	if err != nil {
		return nil, err
	}

	_, errLK := s.session.ExecuteWrite(s.ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, errLK := transaction.Run(s.ctx, queryLK, map[string]interface{}{"id": i})
		if errLK != nil {
			return nil, errLK
		}

		for result.Next(s.ctx) {
			results = append(results, result.Record().AsMap()["relation"].(neo4j.Relationship))
		}

		return nil, result.Err()
	})

	if errLK != nil {
		return nil, errLK
	}

	_, errBLK := s.session.ExecuteWrite(s.ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, errBLK := transaction.Run(s.ctx, queryBLK, map[string]interface{}{"id": i})
		if errBLK != nil {
			return nil, errBLK
		}

		for result.Next(s.ctx) {
			results = append(results, result.Record().AsMap()["relation"].(neo4j.Relationship))
		}

		return nil, result.Err()
	})

	if errBLK != nil {
		return nil, errBLK
	}

	sumRating = 0.0
	for r := 0; r < len(results); r++ {
		props := results[r].Props
		rating := props["rating"]
		if rating != -1 {
			sumRating = sumRating + rating.(float64)
		}
	}

	return &GetRating{
		MediaID:   i,
		MediaType: tp,
		AvgRating: sumRating / float64(len(results)),
	}, nil
}

func (s *Neo4jStore) GetWishlist(i int) (*GetWishlist, error) {
	queryLK := "MATCH (:User {id_user: $id_user})-[r:LK {wishlist: true}]-(n) RETURN r as relation"
	queryDLK := "MATCH (:User {id_user: $id_user})-[r:DLK {wishlist: true}]-(n) RETURN r as relation"
	queryBLK := "MATCH (:User {id_user: $id_user})-[r:BLK {wishlist: true}]-(n) RETURN r as relation"

	var results []neo4j.Relationship
	var movies []int64
	var songs []int64
	var books []int64

	_, err := s.session.ExecuteWrite(s.ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(s.ctx, queryDLK, map[string]interface{}{"id_user": i})
		if err != nil {
			return nil, err
		}

		for result.Next(s.ctx) {
			results = append(results, result.Record().AsMap()["relation"].(neo4j.Relationship))
		}

		return nil, result.Err()
	})

	if err != nil {
		return nil, err
	}

	_, errLK := s.session.ExecuteWrite(s.ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, errLK := transaction.Run(s.ctx, queryLK, map[string]interface{}{"id_user": i})
		if errLK != nil {
			return nil, errLK
		}

		for result.Next(s.ctx) {
			results = append(results, result.Record().AsMap()["relation"].(neo4j.Relationship))
		}

		return nil, result.Err()
	})

	if errLK != nil {
		return nil, errLK
	}

	_, errBLK := s.session.ExecuteWrite(s.ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, errBLK := transaction.Run(s.ctx, queryBLK, map[string]interface{}{"id_user": i})
		if errBLK != nil {
			return nil, errBLK
		}

		for result.Next(s.ctx) {
			results = append(results, result.Record().AsMap()["relation"].(neo4j.Relationship))
		}

		return nil, result.Err()
	})

	if errBLK != nil {
		return nil, errBLK
	}

	for r := 0; r < len(results); r++ {
		props := results[r].Props
		mediaType := props["media_type"]
		mediaID := props["media_id"].(int64)

		if mediaType == "MOV" {
			movies = append(movies, mediaID)
		} else if mediaType == "SON" {
			songs = append(songs, mediaID)
		} else {
			books = append(books, mediaID)
		}
	}

	return &GetWishlist{
		UserID: i,
		Movies: movies,
		Songs:  songs,
		Books:  books,
	}, nil
}
