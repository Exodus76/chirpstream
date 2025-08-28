package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

type CreateLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func authMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		next(w, r, p)
	}
}

// router mux same stuff
func (h *Handler) RegisterRoutes(router *httprouter.Router) {
	router.POST("/api/user/register", h.handleCreateUser)
	router.POST("/api/user/login", h.handleUserLogin)

	router.GET("/api/user/getUser/:id", authMiddleware(h.handleGetuser))
}

func (h *Handler) handleCreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.service.RegisterNewUser(&req)
	if err != nil {
		log.Printf("Error registering new user %v \n", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) handleUserLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req CreateLoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "handle login user")
}

func (h *Handler) handleGetuser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	userId := p.ByName("id")

	fmt.Fprintf(w, "handle get user %v", userId)
}
