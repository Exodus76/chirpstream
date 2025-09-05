package chirps

import (
	"chirpstream/internal/auth"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	service Service
}

type CreateChirpRequest struct {
	Content string `json:"content"`
}

type UpdateChirpRequest struct {
	Content string `json:"content"`
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router *httprouter.Router) {
	router.POST("/api/chirp/create", auth.AuthMiddleware(h.handleCreateChirp))
	router.POST("/api/chirp/update", auth.AuthMiddleware(h.handleUpdateChirp))

	router.POST("/api/chirp/getChripById", auth.AuthMiddleware(h.handleGetChirpById))

	router.POST("/api/chirp/getChirpsByUserId", auth.AuthMiddleware(h.handleGetChirpsByUserId))
}

func (h *Handler) handleCreateChirp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req CreateChirpRequest

	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	defer r.Body.Close()

	user := r.Context().Value("user").(*auth.CustomClaim)

	//do the jwt stuuff before implementing this
	err := h.service.CreateChirp(ctx, req.Content, int(user.UserID))
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusBadRequest)
	}
}

func (h *Handler) handleUpdateChirp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req UpdateChirpRequest

	_ = r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}
}

func (h *Handler) handleGetChirpById(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

func (h *Handler) handleGetChirpsByUserId(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}
