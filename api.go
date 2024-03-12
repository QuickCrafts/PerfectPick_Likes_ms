package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// API Structure

type APIServer struct {
	listenAddr string
	store      Storage
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
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

	// return WriteJSON(w, http.StatusBadRequest, "Guard failed") // 400
	// return WriteJSON(w, http.StatusBadRequest, "User and Media already has a relation") // 400
	// return WriteJSON(w, http.StatusInternalServerError, error) // 500

	return WriteJSON(w, http.StatusCreated, like) // 201
}

func (s *APIServer) handleUpdateLike(w http.ResponseWriter, r *http.Request) error {
	// @todo
	// Need body
	params := mux.Vars(r)

	// return WriteJSON(w, http.StatusBadRequest, "Guard failed") // 400
	// return WriteJSON(w, http.StatusBadRequest, "User id not provided") // 400
	// return WriteJSON(w, http.StatusBadRequest, "Media id not provided") // 400
	// return WriteJSON(w, http.StatusBadRequest, "Media type not provided") // 400
	// return WriteJSON(w, http.StatusNotFound, "Relation not found") // 404
	// return WriteJSON(w, http.StatusInternalServerError, error) // 500

	return WriteJSON(w, http.StatusCreated, params) //201
}

func (s *APIServer) handleDeleteLike(w http.ResponseWriter, r *http.Request) error {
	// @todo
	params := mux.Vars(r)

	// return WriteJSON(w, http.StatusBadRequest, "User id not provided") // 400
	// return WriteJSON(w, http.StatusBadRequest, "Media id not provided") // 400
	// return WriteJSON(w, http.StatusBadRequest, "Media type not provided") // 400
	// return WriteJSON(w, http.StatusNotFound, "Relation not found") // 404
	// return WriteJSON(w, http.StatusInternalServerError, error) // 500

	return WriteJSON(w, http.StatusNoContent, params) // 204
}

// /likes/user Functions

func (s *APIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)

	if params["id"] == "" {
		return WriteJSON(w, http.StatusBadRequest, "User id not provided") // 400
	}

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, err) // 500
	}

	if err := s.store.CreateUser(id); err != nil {
		return WriteJSON(w, http.StatusInternalServerError, err) // 500
	}

	return WriteJSON(w, http.StatusCreated, "User created") // 201
}

func (s *APIServer) handleGetUserLikes(w http.ResponseWriter, r *http.Request) error {
	// @todo
	params := mux.Vars(r)

	// return WriteJSON(w, http.StatusBadRequest, "User id not provided") // 400
	// return WriteJSON(w, http.StatusNotFound, "User not found") // 404
	// return WriteJSON(w, http.StatusInternalServerError, error) // 500

	return WriteJSON(w, http.StatusOK, params)
}

func (s *APIServer) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)

	if params["id"] == "" {
		return WriteJSON(w, http.StatusBadRequest, "User id not provided") // 400
	}

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, err) // 500
	}

	if err := s.store.DeleteUser(id); err != nil {
		return WriteJSON(w, http.StatusInternalServerError, err) // 500
	}

	return WriteJSON(w, http.StatusNoContent, "")
}

// /likes/media Functions

func (s *APIServer) handleCreateMedia(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)

	if params["id"] == "" {
		return WriteJSON(w, http.StatusBadRequest, "Media id not provided") // 400
	}

	if params["media_type"] == "" {
		return WriteJSON(w, http.StatusBadRequest, "Media type not provided") // 400
	}

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, err) // 500
	}

	if err := s.store.CreateMedia(id, params["media_type"]); err != nil {
		return WriteJSON(w, http.StatusInternalServerError, err) // 500
	}

	return WriteJSON(w, http.StatusCreated, "Media created") // 201
}

func (s *APIServer) handleGetMediaLikes(w http.ResponseWriter, r *http.Request) error {
	// @todo
	params := mux.Vars(r)

	// return WriteJSON(w, http.StatusBadRequest, "Media id not provided") // 400
	// return WriteJSON(w, http.StatusBadRequest, "Media type not provided") // 400
	// return WriteJSON(w, http.StatusNotFound, "Media not found") // 404
	// return WriteJSON(w, http.StatusInternalServerError, error) // 500

	return WriteJSON(w, http.StatusOK, params)
}

func (s *APIServer) handleDeleteMedia(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)

	if params["id"] == "" {
		return WriteJSON(w, http.StatusBadRequest, "Media id not provided") // 400
	}

	if params["media_type"] == "" {
		return WriteJSON(w, http.StatusBadRequest, "Media type not provided") // 400
	}

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, "Media id not provided") // 400
	}

	if err := s.store.DeleteMedia(id, params["media_type"]); err != nil {
		return WriteJSON(w, http.StatusInternalServerError, err) // 500
	}

	return WriteJSON(w, http.StatusNoContent, "")
}

// /likes/average Functions

func (s *APIServer) handleAverage(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "GET" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	// @todo
	params := mux.Vars(r)

	// return WriteJSON(w, http.StatusBadRequest, "Media id not provided") // 400
	// return WriteJSON(w, http.StatusBadRequest, "Media type not provided") // 400
	// return WriteJSON(w, http.StatusNotFound, "Media not found") // 404
	// return WriteJSON(w, http.StatusInternalServerError, error) // 500

	return WriteJSON(w, http.StatusOK, params)
}

// /likes/wishlist Functions

func (s *APIServer) handleWishlist(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "GET" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	// @todo
	params := mux.Vars(r)

	// return WriteJSON(w, http.StatusBadRequest, "User id not provided") // 400
	// return WriteJSON(w, http.StatusNotFound, "User not found") // 404
	// return WriteJSON(w, http.StatusInternalServerError, error) // 500

	return WriteJSON(w, http.StatusOK, params)
}
