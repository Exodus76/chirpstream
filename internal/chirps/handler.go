package chirps

import (
	"chirpstream/internal/auth"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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

	router.GET("/api/chirp/getChirpById/:chirpId", auth.AuthMiddleware(h.handleGetChirpById))

	router.GET("/api/chirp/getChirpsByUserId/:userId", auth.AuthMiddleware(h.handleGetChirpsByUserId))
}

func (h *Handler) handleCreateChirp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req CreateChirpRequest

	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error: failed decoding json body %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	user := r.Context().Value("user").(*auth.CustomClaim)

	err := h.service.CreateChirp(ctx, req.Content, int(user.UserID))
	if err != nil {
		log.Printf("Error: failed creating chirp %v\n", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) handleUpdateChirp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req UpdateChirpRequest

	_ = r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
}

func (h *Handler) handleGetChirpById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var chirp *ChirpWithLikes
	var err error
	ctx := r.Context()

	chirpIdParam := p.ByName("chirpId")

	id, err := strconv.Atoi(chirpIdParam)
	if err != nil {
		log.Printf("Error: failed converting paramter to int: %v\n", err)
		http.Error(w, "Invalid param value", http.StatusBadRequest)
		return
	}

	chirp, err = h.service.GetChirpWithLikesById(ctx, id)
	if err != nil {
		log.Printf("Error: cant get chirp %v\n", err)
		http.Error(w, "No chirp found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	payload, err := json.Marshal(chirp)
	if err != nil {
		log.Printf("Error: marshalling JSON: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

func (h *Handler) handleGetChirpsByUserId(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}
