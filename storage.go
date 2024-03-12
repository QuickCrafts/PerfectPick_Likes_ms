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

func (s *Neo4jStore) CreateLike(*Like) error {
	return nil
}

// Update Functions
func (s *Neo4jStore) UpdateLike(*UpdateLike) error {
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
	return nil
}

// Get Functions
func (s *Neo4jStore) GetUserLikes(int) (*GetUserLikes, error) {
	return nil, nil
}

func (s *Neo4jStore) GetMediaLikes(i int, tp string) (*GetMediaLikes, error) {
	return nil, nil
}

func (s *Neo4jStore) GetAverage(i int, tp string) (*GetRating, error) {
	return nil, nil
}

func (s *Neo4jStore) GetWishlist(int) (*GetWishlist, error) {
	return nil, nil
}
