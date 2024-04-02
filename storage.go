package main

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Storage interface {
	// Create
	CreateUser(int) error
	CreateMedia(string, string) error
	SetLike(*Like) error
	AddToWishlist(int, string, string) error
	SetAverage(int, string, string, float64) error

	// Get
	GetUserLikes(int, string, string) (*GetUserLikes, error)
	GetMediaLikes(string, string, string) (*GetMediaLikes, error)
	GetSpecificLike(int, string, string) (*LikeRelation, error)
	GetAverage(string, string) (float64, error)
	GetRating(string, string, int) (float64, error)
	GetWishlist(int, string) (*GetWishlist, error)

	//Delete
	DeleteUser(int) error
	DeleteMedia(string, string) error
	DeleteLike(int, string, string) error
	RemoveFromWishlist(int, string, string) error

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
	dbUri := "neo4j://neo4j:7687" //"neo4j://localhost:7000" --> for local | "neo4j://neo4j:7687"  ---> for docker
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

func (s *Neo4jStore) CreateMedia(i string, tp string) error {

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

func (s *Neo4jStore) SetLike(l *Like) error {

	query := `
	MERGE (n:User {id_user: $id_user})
	MERGE (m:Movie {id_movie: $id_media})
	MERGE (n)-[r:PREF]->(m)
	ON CREATE
		SET
			r.type = $type
			r.media_id = $id_media
			r.media_type = MOV
			r.user_id = $id_user
	ON MATCH
		SET
			r.type = $type
			r.media_id = $id_media
			r.media_type = MOV
			r.user_id = $id_user
	`

	if l.MediaType == "SON" {
		query = `
		MERGE (n:User {id_user: $id_user})
		MERGE (m:Song {id_song: $id_media})
		MERGE (n)-[r:PREF]->(m)
		ON CREATE
			SET
				r.type = $type
				r.media_id = $id_media
				r.media_type = SON
				r.user_id = $id_user
		ON MATCH
			SET
				r.type = $type
				r.media_id = $id_media
				r.media_type = SON
				r.user_id = $id_user
		`
	} else if l.MediaType == "BOO" {
		query = `
		MERGE (n:User {id_user: $id_user})
		MERGE (m:Book {id_book: $id_media})
		MERGE (n)-[r:PREF]->(m)
		ON CREATE
			SET
				r.type = $type
				r.media_id = $id_media
				r.media_type = BOO
				r.user_id = $id_user
		ON MATCH
			SET
				r.type = $type
				r.media_id = $id_media
				r.media_type = BOO
				r.user_id = $id_user
		`
	}

	_, err := s.session.ExecuteWrite(s.ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(s.ctx, query, map[string]interface{}{"id_media": l.MediaID, "id_user": l.UserID, "type": l.LikeType})
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

// Delete Functions
func (s *Neo4jStore) DeleteUser(i int) error {
	_, err := s.session.ExecuteWrite(s.ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(s.ctx, "MATCH (u:User) WHERE u.id_user = $id DETACH DELETE u", map[string]interface{}{"id": i})
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

func (s *Neo4jStore) DeleteMedia(i string, tp string) error {
	query := "MATCH (m:Movie) WHERE m.id_movie = $id DETACH DELETE m"

	if tp == "SON" {
		query = "MATCH (s:Song) WHERE s.id_song = $id DETACH DELETE s"
	} else if tp == "BOO" {
		query = "MATCH (b:Book) WHERE b.id_book = $id DETACH DELETE b"
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

func (s *Neo4jStore) DeleteLike(user_id int, media_id string, tp string) error {
	queryLK := "MATCH (:Movie {id_movie: $id_media})-[r:PREF]-(:User {id_user: $id_user}) DELETE r"

	if tp == "SON" {
		queryLK = "MATCH (:Song {id_song: $id_media})-[r:PREF]-(:User {id_user: $id_user}) DELETE r"
	} else if tp == "BOO" {
		queryLK = "MATCH (:Book {id_book: $id_media})-[r:PREF]-(:User {id_user: $id_user}) DELETE r"
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

	return nil
}

// Get Functions
func (s *Neo4jStore) GetUserLikes(i int, media string, tp string) (*GetUserLikes, error) {
	queryLK := "MATCH (:User {id_user: $id_user})-[r:PREF]-(n) RETURN r as relation"

	if media == "SON" {
		queryLK = "MATCH (:User {id_user: $id_user})-[r:PREF]-(:Song) RETURN r as relation"
	} else if media == "BOO" {
		queryLK = "MATCH (:User {id_user: $id_user})-[r:PREF]-(:Book) RETURN r as relation"
	} else if media == "MOV" {
		queryLK = "MATCH (:User {id_user: $id_user})-[r:PREF]-(:Movie) RETURN r as relation"
	}

	if tp == "LK" {
		queryLK = "MATCH (:User {id_user: $id_user})-[r:PREF {type: LK}]-(:Song) RETURN r as relation"
	} else if tp == "DLK" {
		queryLK = "MATCH (:User {id_user: $id_user})-[r:PREF {type: DLK}]-(:Book) RETURN r as relation"
	}

	var results []neo4j.Relationship
	var movies []LikeRelation
	var songs []LikeRelation
	var books []LikeRelation

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

	for r := 0; r < len(results); r++ {
		props := results[r].Props
		mediaType := props["media_type"]
		like := NewLikeRelation(i, props["media_id"], mediaType, props["type"])

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

func (s *Neo4jStore) GetMediaLikes(i string, media string, tp string) (*GetMediaLikes, error) {

	queryLK := "MATCH (:Movie {id_movie: $id})-[r:PREF]-(n) RETURN r as relation"

	if tp == "LK" {
		queryLK = "MATCH (:Movie {id_movie: $id})-[r:PREF {type: LK}]-(n) RETURN r as relation"
	} else if tp == "DLK" {
		queryLK = "MATCH (:Movie {id_movie: $id})-[r:PREF {type: DLK}]-(n) RETURN r as relation"
	}

	if media == "SON" {
		queryLK = "MATCH (:Song {id_song: $id})-[r:PREF]-(n) RETURN r as relation"

		if tp == "LK" {
			queryLK = "MATCH (:Song {id_song: $id})-[r:PREF {type: LK}]-(n) RETURN r as relation"
		} else if tp == "DLK" {
			queryLK = "MATCH (:Song {id_song: $id})-[r:PREF {type: DLK}]-(n) RETURN r as relation"
		}
	} else if media == "BOO" {
		queryLK = "MATCH (:Book {id_book: $id})-[r:PREF]-(n) RETURN r as relation"

		if tp == "LK" {
			queryLK = "MATCH (:Book {id_book: $id})-[r:PREF {type: LK}]-(n) RETURN r as relation"
		} else if tp == "DLK" {
			queryLK = "MATCH (:Book {id_book: $id})-[r:PREF {type: DLK}]-(n) RETURN r as relation"
		}
	}

	var results []neo4j.Relationship
	var likes []LikeRelation

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

	for r := 0; r < len(results); r++ {
		props := results[r].Props
		mediaType := props["media_type"]
		like := NewLikeRelation(props["user_id"], i, mediaType, props["type"])
		likes = append(likes, *like)
	}

	return &GetMediaLikes{
		Likes: likes,
	}, nil
}

func (s *Neo4jStore) GetSpecificLike(i int, media_id string, media string) (*LikeRelation, error) {

	queryLK := "MATCH (:Movie {id_movie: $id})-[r:PREF]-(:User {id_user: $user_id}) RETURN r as relation"

	if media == "SON" {
		queryLK = "MATCH (:Song {id_song: $id})-[r:PREF]-(:User {id_user: $user_id}) RETURN r as relation"
	} else if media == "BOO" {
		queryLK = "MATCH (:Book {id_book: $id})-[r:PREF]-(:User {id_user: $user_id}) RETURN r as relation"
	}

	var results []neo4j.Relationship

	_, errLK := s.session.ExecuteWrite(s.ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, errLK := transaction.Run(s.ctx, queryLK, map[string]interface{}{"id": media_id, "user_id": i})
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

	props := results[0].Props
	mediaType := props["media_type"]
	like := NewLikeRelation(i, media_id, mediaType, props["type"])

	return like, nil
}

func (s *Neo4jStore) GetAverage(i string, tp string) (float64, error) {
	queryLK := "MATCH (:Movie {id_movie: $id})-[r:RTE]-(n) RETURN r as relation"

	if tp == "SON" {
		queryLK = "MATCH (:Song {id_song: $id})-[r:RTE]-(n) RETURN r as relation"
	} else if tp == "BOO" {
		queryLK = "MATCH (:Book {id_book: $id})-[r:RTE]-(n) RETURN r as relation"
	}

	var results []neo4j.Relationship
	var sumRating float64

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
		return 0.0, errLK
	}

	sumRating = 0.0
	for r := 0; r < len(results); r++ {
		props := results[r].Props
		rating := props["rating"]
		if rating != -1 {
			sumRating = sumRating + rating.(float64)
		}
	}

	return sumRating / float64(len(results)), nil
}

func (s *Neo4jStore) GetRating(i string, tp string, u int) (float64, error) {
	queryLK := "MATCH (:Movie {id_movie: $id})-[r:RTE]-(:User {id_user:$user_id}) RETURN r as relation"

	if tp == "SON" {
		queryLK = "MATCH (:Song {id_song: $id})-[r:RTE]-(:User {id_user:$user_id}) RETURN r as relation"
	} else if tp == "BOO" {
		queryLK = "MATCH (:Book {id_book: $id})-[r:RTE]-(:User {id_user:$user_id}) RETURN r as relation"
	}

	var results []neo4j.Relationship

	_, errLK := s.session.ExecuteWrite(s.ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, errLK := transaction.Run(s.ctx, queryLK, map[string]interface{}{"id": i, "user_id": u})
		if errLK != nil {
			return nil, errLK
		}

		for result.Next(s.ctx) {
			results = append(results, result.Record().AsMap()["relation"].(neo4j.Relationship))
		}

		return nil, result.Err()
	})

	if errLK != nil {
		return 0.0, errLK
	}

	props := results[0].Props

	return props["rating"].(float64), nil
}

func (s *Neo4jStore) SetAverage(i int, md string, tp string, rate float64) error {
	query := `
	MERGE (n:User {id_user: $id_user})
	MERGE (m:Movie {id_movie: $id_media})
	MERGE (n)-[r:RTE]->(m)
	ON CREATE
		SET
			r.media_id = $id_media
			r.media_type = MOV
			r.user_id = $id_user
	ON MATCH
		SET
			r.rating = $rate
			r.media_id = $id_media
			r.media_type = MOV
			r.user_id = $id_user
	`

	if tp == "SON" {
		query = `
		MERGE (n:User {id_user: $id_user})
		MERGE (m:Song {id_song: $id_media})
		MERGE (n)-[r:RTE]->(m)
		ON CREATE
			SET
				r.rating = $rate
				r.media_id = $id_media
				r.media_type = SON
				r.user_id = $id_user
		ON MATCH
			SET
				r.rating = $rate
				r.media_id = $id_media
				r.media_type = SON
				r.user_id = $id_user
		`
	} else if tp == "BOO" {
		query = `
		MERGE (n:User {id_user: $id_user})
		MERGE (m:Book {id_book: $id_media})
		MERGE (n)-[r:RTE]->(m)
		ON CREATE
			SET
				r.rating = $rate
				r.media_id = $id_media
				r.media_type = BOO
				r.user_id = $id_user
		ON MATCH
			SET
				r.rating = $rate
				r.media_id = $id_media
				r.media_type = BOO
				r.user_id = $id_user
		`
	}

	_, err := s.session.ExecuteWrite(s.ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(s.ctx, query, map[string]interface{}{"id_media": md, "id_user": i, "rate": rate})
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

func (s *Neo4jStore) GetWishlist(i int, tp string) (*GetWishlist, error) {
	queryLK := "MATCH (:User {id_user: $id_user})-[r:WSH]-(n) RETURN r as relation"

	if tp == "SON" {
		queryLK = "MATCH (:User {id_user: $id_user})-[r:WSH]-(:Song) RETURN r as relation"
	} else if tp == "BOO" {
		queryLK = "MATCH (:User {id_user: $id_user})-[r:WSH]-(:Book) RETURN r as relation"
	} else if tp == "MOV" {
		queryLK = "MATCH (:User {id_user: $id_user})-[r:WSH]-(:Movie) RETURN r as relation"
	}

	var results []neo4j.Relationship
	var movies []string
	var songs []string
	var books []string

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

	for r := 0; r < len(results); r++ {
		props := results[r].Props
		mediaType := props["media_type"]
		mediaID := props["media_id"].(string)

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

func (s *Neo4jStore) AddToWishlist(i int, md string, tp string) error {
	query := `
	MERGE (n:User {id_user: $id_user})
	MERGE (m:Movie {id_movie: $id_media})
	MERGE (n)-[r:WSH]->(m)
	ON CREATE
		SET
			r.media_id = $id_media
			r.media_type = MOV
			r.user_id = $id_user
	ON MATCH
		SET
			r.media_id = $id_media
			r.media_type = MOV
			r.user_id = $id_user
	`

	if tp == "SON" {
		query = `
		MERGE (n:User {id_user: $id_user})
		MERGE (m:Song {id_song: $id_media})
		MERGE (n)-[r:WSH]->(m)
		ON CREATE
			SET
				r.media_id = $id_media
				r.media_type = SON
				r.user_id = $id_user
		ON MATCH
			SET
				r.media_id = $id_media
				r.media_type = SON
				r.user_id = $id_user
		`
	} else if tp == "BOO" {
		query = `
		MERGE (n:User {id_user: $id_user})
		MERGE (m:Book {id_book: $id_media})
		MERGE (n)-[r:WSH]->(m)
		ON CREATE
			SET
				r.media_id = $id_media
				r.media_type = BOO
				r.user_id = $id_user
		ON MATCH
			SET
				r.media_id = $id_media
				r.media_type = BOO
				r.user_id = $id_user
		`
	}

	_, err := s.session.ExecuteWrite(s.ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(s.ctx, query, map[string]interface{}{"id_media": md, "id_user": i})
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

func (s *Neo4jStore) RemoveFromWishlist(user_id int, media_id string, tp string) error {
	queryLK := "MATCH (:Movie {id_movie: $id_media})-[r:WSH]-(:User {id_user: $id_user}) DELETE r"

	if tp == "SON" {
		queryLK = "MATCH (:Song {id_song: $id_media})-[r:WSH]-(:User {id_user: $id_user}) DELETE r"
	} else if tp == "BOO" {
		queryLK = "MATCH (:Book {id_book: $id_media})-[r:WSH]-(:User {id_user: $id_user}) DELETE r"
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

	return nil
}
