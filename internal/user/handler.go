package user

import (
	"chirpstream/internal/auth"
	"chirpstream/pkg/response"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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
		response.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	err := h.service.CreateUser(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		log.Printf("ERROR: registering new user %v \n", err)
		response.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusCreated, map[string]string{"message": "User created successfully"})
}

func (h *Handler) handleUserLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req CreateLoginRequest

	ctx := r.Context()

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: error parsing request %v", err)
		response.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.service.VerifyUser(ctx, req.Email, req.Password)
	if err != nil {
		log.Printf("ERROR: error verifying %v\n", err)
		response.Error(w, "Something went wrong", http.StatusUnauthorized)
		return
	}

	if user == nil {
		response.Error(w, "Invalid credentials", http.StatusUnauthorized)
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
		log.Printf("ERROR: error signing string %v", err)
		return
	}

	payload := &LoginResponse{
		Token:  ss,
		UserID: user.ID,
	}

	response.JSON(w, http.StatusOK, payload)
}

func (h *Handler) handleGetuser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	userId, _ := strconv.Atoi(p.ByName("id"))

	user, err := h.service.GetUserById(ctx, userId)
	if err != nil {
		log.Printf("ERROR: error getting user by id %v\n", err)
		response.Error(w, "Something went wrong", http.StatusInternalServerError)
	}

	response.JSON(w, http.StatusOK, user)
}
