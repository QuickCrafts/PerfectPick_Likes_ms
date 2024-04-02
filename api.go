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
	router.HandleFunc("/likes/user/{id}", makeHTTPHandleFunc(s.handleUser)).Queries("media_type", "{media_type}", "preference", "{preference}")
	router.HandleFunc("/likes/user/{id}", makeHTTPHandleFunc(s.handleUser)).Queries("media_type", "{media_type}")
	router.HandleFunc("/likes/user/{id}", makeHTTPHandleFunc(s.handleUser)).Queries("preference", "{preference}")
	router.HandleFunc("/likes/user/{id}", makeHTTPHandleFunc(s.handleUser))
	router.HandleFunc("/likes/media/{id}", makeHTTPHandleFunc(s.handleMedia)).Queries("media_type", "{media_type}", "preference", "{preference}")
	router.HandleFunc("/likes/media/{id}", makeHTTPHandleFunc(s.handleMedia)).Queries("media_type", "{media_type}")
	router.HandleFunc("/likes/rate/{id}", makeHTTPHandleFunc(s.handleRate)).Queries("media_type", "{media_type}")
	router.HandleFunc("/likes/rate/{id}", makeHTTPHandleFunc(s.handleRate)).Queries("media_type", "{media_type}", "user_id", "{user_id}")
	router.HandleFunc("/likes/wishlist/{id}", makeHTTPHandleFunc(s.handleWishlist)).Queries("media_type", "{media_type}")
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
	if r.Method == "GET" {
		return s.handleSpecificLike(w, r)
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

func (s *APIServer) handleRate(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		return s.handleCreateRate(w, r)
	}
	if r.Method == "GET" {
		return s.handleAverage(w, r)
	}
	if r.Method == "PUT" {
		return s.handleUpdateRate(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleWishlist(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		return s.handleSetMediaWish(w, r)
	}
	if r.Method == "GET" {
		return s.handleGetWishlist(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

// /likes Functions

func (s *APIServer) handleCreateLike(w http.ResponseWriter, r *http.Request) error {
	createLike := new(Like)

	if err := json.NewDecoder(r.Body).Decode(createLike); err != nil {
		return WriteJSON(w, http.StatusBadRequest, "Guard failed") // 400
	}

	like := NewLike(createLike.UserID, createLike.MediaID, createLike.MediaType, createLike.LikeType)
	if err := s.store.SetLike(like); err != nil {
		return WriteJSON(w, http.StatusInternalServerError, err) // 500
	}

	return WriteJSON(w, http.StatusCreated, "Relation created") // 201
}

func (s *APIServer) handleUpdateLike(w http.ResponseWriter, r *http.Request) error {
	createLike := new(Like)

	if err := json.NewDecoder(r.Body).Decode(createLike); err != nil {
		return WriteJSON(w, http.StatusBadRequest, "Guard failed") // 400
	}

	like := NewLike(createLike.UserID, createLike.MediaID, createLike.MediaType, createLike.LikeType)
	if err := s.store.SetLike(like); err != nil {
		return WriteJSON(w, http.StatusInternalServerError, err) // 500
	}

	return WriteJSON(w, http.StatusCreated, "Relation updated") // 201
}

func (s *APIServer) handleDeleteLike(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)

	if params["user_id"] == "" {
		return WriteJSON(w, http.StatusBadRequest, "User id not provided") // 400
	}

	if params["media_id"] == "" {
		return WriteJSON(w, http.StatusBadRequest, "Media id not provided") // 400
	}

	if params["media_type"] == "" {
		return WriteJSON(w, http.StatusBadRequest, "Media type not provided") // 400
	}

	user_id, err := strconv.Atoi(params["user_id"])
	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, err) // 500
	}

	if err := s.store.DeleteLike(user_id, params["media_id"], params["media_type"]); err != nil {
		return WriteJSON(w, http.StatusInternalServerError, err) // 500
	}

	return WriteJSON(w, http.StatusNoContent, params) // 204
}

func (s *APIServer) handleSpecificLike(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)

	if params["user_id"] == "" {
		return WriteJSON(w, http.StatusBadRequest, "User id not provided") // 400
	}

	if params["media_id"] == "" {
		return WriteJSON(w, http.StatusBadRequest, "Media id not provided") // 400
	}

	if params["media_type"] == "" {
		return WriteJSON(w, http.StatusBadRequest, "Media type not provided") // 400
	}

	user_id, err := strconv.Atoi(params["user_id"])
	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, err) // 500
	}

	result, err := s.store.GetSpecificLike(user_id, params["media_id"], params["media_type"])

	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, err) // 500
	}

	return WriteJSON(w, http.StatusOK, result)
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
	params := mux.Vars(r)

	if params["id"] == "" {
		return WriteJSON(w, http.StatusBadRequest, "User id not provided") // 400
	}

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, err) // 500
	}

	result, err := s.store.GetUserLikes(id, params["media_type"], params["preference"])

	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, err) // 500
	}

	return WriteJSON(w, http.StatusOK, result)
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

	if err := s.store.CreateMedia(params["id"], params["media_type"]); err != nil {
		return WriteJSON(w, http.StatusInternalServerError, err) // 500
	}

	return WriteJSON(w, http.StatusCreated, "Media created") // 201
}

func (s *APIServer) handleGetMediaLikes(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)

	if params["id"] == "" {
		return WriteJSON(w, http.StatusBadRequest, "Media id not provided") // 400
	}

	if params["media_type"] == "" {
		return WriteJSON(w, http.StatusBadRequest, "Media type not provided") // 400
	}

	result, err := s.store.GetMediaLikes(params["id"], params["media_type"], params["preference"])

	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, err) // 500
	}

	return WriteJSON(w, http.StatusOK, result)
}

func (s *APIServer) handleDeleteMedia(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)

	if params["id"] == "" {
		return WriteJSON(w, http.StatusBadRequest, "Media id not provided") // 400
	}

	if params["media_type"] == "" {
		return WriteJSON(w, http.StatusBadRequest, "Media type not provided") // 400
	}

	if err := s.store.DeleteMedia(params["id"], params["media_type"]); err != nil {
		return WriteJSON(w, http.StatusInternalServerError, err) // 500
	}

	return WriteJSON(w, http.StatusNoContent, "")
}

// /likes/rate Functions

func (s *APIServer) handleAverage(w http.ResponseWriter, r *http.Request) error {

	params := mux.Vars(r)

	if params["id"] == "" {
		return WriteJSON(w, http.StatusBadRequest, "Media id not provided") // 400
	}

	if params["media_type"] == "" {
		return WriteJSON(w, http.StatusBadRequest, "Media type not provided") // 400
	}

	user_id, err := strconv.Atoi(params["user_id"])
	if err == nil {
		result, err := s.store.GetRating(params["id"], params["media_type"], user_id)

		if err != nil {
			return WriteJSON(w, http.StatusInternalServerError, err) // 500
		}

		return WriteJSON(w, http.StatusOK, result)
	} else {
		result, err := s.store.GetAverage(params["id"], params["media_type"])

		if err != nil {
			return WriteJSON(w, http.StatusInternalServerError, err) // 500
		}

		return WriteJSON(w, http.StatusOK, result)
	}
}

func (s *APIServer) handleCreateRate(w http.ResponseWriter, r *http.Request) error {

	rating := new(Rate)

	if err := json.NewDecoder(r.Body).Decode(rating); err != nil {
		return WriteJSON(w, http.StatusBadRequest, "Guard failed") // 400
	}

	params := mux.Vars(r)

	if params["id"] == "" {
		return WriteJSON(w, http.StatusBadRequest, "Media id not provided") // 400
	}

	if params["media_type"] == "" {
		return WriteJSON(w, http.StatusBadRequest, "Media type not provided") // 400
	}

	if params["user_id"] == "" {
		return WriteJSON(w, http.StatusBadRequest, "User id not provided") // 400
	}

	user_id, errUser := strconv.Atoi(params["user_id"])
	if errUser != nil {
		return WriteJSON(w, http.StatusInternalServerError, errUser) // 500
	}

	if err := s.store.SetAverage(user_id, params["id"], params["media_type"], rating.Rating); err != nil {
		return WriteJSON(w, http.StatusInternalServerError, err) // 500
	}

	return WriteJSON(w, http.StatusCreated, "Rate added") // 201
}

func (s *APIServer) handleUpdateRate(w http.ResponseWriter, r *http.Request) error {

	rating := new(Rate)

	if err := json.NewDecoder(r.Body).Decode(rating); err != nil {
		return WriteJSON(w, http.StatusBadRequest, "Guard failed") // 400
	}

	params := mux.Vars(r)

	if params["id"] == "" {
		return WriteJSON(w, http.StatusBadRequest, "Media id not provided") // 400
	}

	if params["media_type"] == "" {
		return WriteJSON(w, http.StatusBadRequest, "Media type not provided") // 400
	}

	if params["user_id"] == "" {
		return WriteJSON(w, http.StatusBadRequest, "User id not provided") // 400
	}

	user_id, errUser := strconv.Atoi(params["user_id"])
	if errUser != nil {
		return WriteJSON(w, http.StatusInternalServerError, errUser) // 500
	}

	if err := s.store.SetAverage(user_id, params["id"], params["media_type"], rating.Rating); err != nil {
		return WriteJSON(w, http.StatusInternalServerError, err) // 500
	}

	return WriteJSON(w, http.StatusCreated, "Rate updated") // 201
}

// /likes/wishlist Functions

func (s *APIServer) handleGetWishlist(w http.ResponseWriter, r *http.Request) error {

	params := mux.Vars(r)

	if params["id"] == "" {
		return WriteJSON(w, http.StatusBadRequest, "User id not provided") // 400
	}

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, err) // 500
	}

	result, err := s.store.GetWishlist(id, params["media_type"])

	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, err) // 500
	}

	return WriteJSON(w, http.StatusOK, result)
}

func (s *APIServer) handleSetMediaWish(w http.ResponseWriter, r *http.Request) error {

	wish := new(ChangeWishlist)

	if err := json.NewDecoder(r.Body).Decode(wish); err != nil {
		return WriteJSON(w, http.StatusBadRequest, "Guard failed") // 400
	}

	params := mux.Vars(r)

	if params["id"] == "" {
		return WriteJSON(w, http.StatusBadRequest, "User id not provided") // 400
	}

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, "User id not provided") // 400
	}

	if wish.Type == "ADD" {

		if err := s.store.AddToWishlist(id, wish.MediaID, wish.MediaType); err != nil {
			return WriteJSON(w, http.StatusInternalServerError, err) // 500
		}

		return WriteJSON(w, http.StatusCreated, "Media added to user wishlist") // 201

	} else if wish.Type == "RMV" {

		if err := s.store.RemoveFromWishlist(id, wish.MediaID, wish.MediaType); err != nil {
			return WriteJSON(w, http.StatusInternalServerError, err) // 500
		}

		return WriteJSON(w, http.StatusCreated, "Media removed to user wishlist") // 201

	} else {
		return WriteJSON(w, http.StatusBadRequest, "Action type not allowed") // 400
	}
}
