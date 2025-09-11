package user

import (
	"chirpstream/internal/auth"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token  string `json:"token"`
	UserID int64  `json:"id"`
}

// router mux same stuff
func (h *Handler) RegisterRoutes(router *httprouter.Router) {
	router.POST("/api/user/register", h.handleCreateUser)
	router.POST("/api/user/login", h.handleUserLogin)

	router.GET("/api/user/getUser/:id", auth.AuthMiddleware(h.handleGetuser))
}

func (h *Handler) handleCreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req CreateUserRequest

	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	err := h.service.CreateUser(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		log.Printf("ERROR: registering new user %w \n", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := map[string]string{"message": "User created successfully"}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) handleUserLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req CreateLoginRequest

	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: error parsing request %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	user, err := h.service.VerifyUser(ctx, req.Email, req.Password)
	if err != nil {
		log.Printf("ERROR: error verifying %v\n", err)
		http.Error(w, "Something went wrong", http.StatusUnauthorized)
		return
	}

	if user == nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	claims := &auth.CustomClaim{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "Chirpstream",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte("mykey"))
	if err != nil {
		log.Printf("ERROR: error signing string %w", err)
		return
	}

	response := LoginResponse{
		ss,
		user.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", ss)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) handleGetuser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	userId := p.ByName("id")

	fmt.Fprintf(w, "handle get user %v", userId)
}
