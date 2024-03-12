package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// API Structure

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

type APIServer struct {
	listenAddr string
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/likes", makeHTTPHandleFunc(s.handleLikes)).Queries("media_type", "{media_type}", "user_id", "{user_id}", "media_id", "{media_id}")
	router.HandleFunc("/likes", makeHTTPHandleFunc(s.handleLikes))
	router.HandleFunc("/likes/user/{id}", makeHTTPHandleFunc(s.handleUser))
	router.HandleFunc("/likes/media/{id}", makeHTTPHandleFunc(s.handleMedia)).Queries("media_type", "{media_type}")
	router.HandleFunc("/likes/average/{id}", makeHTTPHandleFunc(s.handleAverage)).Queries("media_type", "{media_type}")
	router.HandleFunc("/likes/wishlist/{id}", makeHTTPHandleFunc(s.handleWishlist))

	log.Println("REST API server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

// Routes Handlers

func (s *APIServer) handleLikes(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		return s.handleCreateLike(w, r)
	}
	if r.Method == "PUT" {
		return s.handleUpdateLike(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteLike(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleUser(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		return s.handleCreateUser(w, r)
	}
	if r.Method == "GET" {
		return s.handleGetUserLikes(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteUser(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleMedia(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		return s.handleCreateMedia(w, r)
	}
	if r.Method == "GET" {
		return s.handleGetMediaLikes(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteMedia(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

// /likes Functions

func (s *APIServer) handleCreateLike(w http.ResponseWriter, r *http.Request) error {
	// @todo
	// Need body
	like := NewLike(1)
	return WriteJSON(w, http.StatusOK, like)
}

func (s *APIServer) handleUpdateLike(w http.ResponseWriter, r *http.Request) error {
	// @todo
	// Need body
	params := mux.Vars(r)

	return WriteJSON(w, http.StatusOK, params)
}

func (s *APIServer) handleDeleteLike(w http.ResponseWriter, r *http.Request) error {
	// @todo
	params := mux.Vars(r)

	return WriteJSON(w, http.StatusOK, params)
}

// /likes/user Functions

func (s *APIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	// @todo
	params := mux.Vars(r)

	return WriteJSON(w, http.StatusOK, params)
}

func (s *APIServer) handleGetUserLikes(w http.ResponseWriter, r *http.Request) error {
	// @todo
	params := mux.Vars(r)

	return WriteJSON(w, http.StatusOK, params)
}

func (s *APIServer) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	// @todo
	params := mux.Vars(r)

	return WriteJSON(w, http.StatusOK, params)
}

// /likes/media Functions

func (s *APIServer) handleCreateMedia(w http.ResponseWriter, r *http.Request) error {
	// @todo
	params := mux.Vars(r)

	return WriteJSON(w, http.StatusOK, params)
}

func (s *APIServer) handleGetMediaLikes(w http.ResponseWriter, r *http.Request) error {
	// @todo
	params := mux.Vars(r)

	return WriteJSON(w, http.StatusOK, params)
}

func (s *APIServer) handleDeleteMedia(w http.ResponseWriter, r *http.Request) error {
	// @todo
	params := mux.Vars(r)

	return WriteJSON(w, http.StatusOK, params)
}

// /likes/average Functions

func (s *APIServer) handleAverage(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "GET" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	// @todo
	params := mux.Vars(r)

	return WriteJSON(w, http.StatusOK, params)
}

// /likes/wishlist Functions

func (s *APIServer) handleWishlist(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "GET" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	// @todo
	params := mux.Vars(r)

	return WriteJSON(w, http.StatusOK, params)
}
