package user

import (
	"fmt"
	"net/http"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) http.Handler {
	userMux := http.NewServeMux()

	userMux.HandleFunc("POST /register", h.handleCreateUser)
	userMux.HandleFunc("POST /login", h.handleUserLogin)

	userMux.HandleFunc("GET /getUser", h.handleGetuser)

	return http.StripPrefix("/api/user", userMux)
}

func (h *Handler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "handle create user")
}

func (h *Handler) handleUserLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "handle login user")
}

func (h *Handler) handleGetuser(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("id")

	fmt.Fprintf(w, "handle get user %v", userId)
}
