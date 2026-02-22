package chirps

import (
	"chirpstream/internal/auth"
	"chirpstream/pkg/response"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gocql/gocql"
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
	router.PUT("/api/chirp/update", auth.AuthMiddleware(h.handleUpdateChirpContent))
	router.DELETE("/api/chirp/delete", auth.AuthMiddleware(h.handleDeleteChirp))

	router.GET("/api/chirp/getChirpById/:chirpId", auth.AuthMiddleware(h.handleGetChirpById))
	router.GET("/api/chirp/getChirpsByUserId/:userId", auth.AuthMiddleware(h.handleGetChirpsByUserId))
}

func (h *Handler) handleCreateChirp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req CreateChirpRequest

	ctx := r.Context()

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error: failed decoding json body %v\n", err)
		response.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, ok := ctx.Value("user").(*auth.CustomClaim)
	if !ok {
		log.Printf("Error: failed getting user from context\n")
		response.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := h.service.CreateChirp(ctx, req.Content, int(user.UserID))
	if err != nil {
		log.Printf("Error: failed creating chirp %v\n", err)
		response.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusCreated, map[string]string{"message": "Chirp created successfully"})
}

func (h *Handler) handleUpdateChirpContent(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req UpdateChirpRequest

	ctx := r.Context()

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, ok := ctx.Value("user").(*auth.CustomClaim)
	if !ok {
		response.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	chirpId := r.URL.Query().Get("chirpId")
	chirpUUID, err := gocql.ParseUUID(chirpId)
	if err != nil {
		log.Printf("Error: failed parsing chirpId parameter %v\n", err)
		response.Error(w, "Something went wrong", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateChirp(ctx, int(user.UserID), chirpUUID, req.Content)
	if err != nil {
		log.Printf("Error: failed updating chirp %v\n", err)
		response.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Chirp updated successfully"})
}

func (h *Handler) handleDeleteChirp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()

	user, ok := ctx.Value("user").(*auth.CustomClaim)
	if !ok {
		log.Printf("Error: failed getting user from context\n")
		response.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	chirpId := r.URL.Query().Get("chirpId")
	chirpUUID, err := gocql.ParseUUID(chirpId)
	if err != nil {
		log.Printf("Error: failed parsing chirpId parameter %v\n", err)
		response.Error(w, "Something went wrong", http.StatusBadRequest)
		return
	}
	err = h.service.DeleteChirp(ctx, int(user.UserID), chirpUUID)
	if err != nil {
		log.Printf("Error: failed deleting chirp %v\n", err)
		response.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "Chirp deleted successfully"})
}

func (h *Handler) handleGetChirpById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var chirp *Chirp
	ctx := r.Context()

	param := p.ByName("chirpId")

	chirpId, err := gocql.ParseUUID(param)
	// chirpId, err := strconv.Atoi(param)
	if err != nil {
		log.Printf("Error: failed parsing paramter to UUID: %v\n", err)
		response.Error(w, "Invalid param value", http.StatusBadRequest)
		return
	}

	user, ok := r.Context().Value("user").(*auth.CustomClaim)
	if !ok {
		log.Printf("Error: failed getting user from context\n")
		response.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	chirp, err = h.service.GetChirpById(ctx, int(user.UserID), chirpId)
	if err != nil {
		log.Printf("Error: cant get chirp %v\n", err)
		response.Error(w, "No chirp found", http.StatusNotFound)
		return
	}

	payload, err := json.Marshal(chirp)
	if err != nil {
		log.Printf("Error: marshalling JSON: %v\n", err)
		response.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusOK, payload)
}

func (h *Handler) handleGetChirpsByUserId(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()

	userId, err := strconv.Atoi(p.ByName("userId"))
	if err != nil {
		log.Printf("Error: failed parsing userId parameter %v\n", err)
		response.Error(w, "Invalid userId parameter", http.StatusBadRequest)
		return
	}

	var ps []byte
	if pageStateBase := r.URL.Query().Get("pageState"); pageStateBase != "" {
		ps, err = base64.StdEncoding.DecodeString(pageStateBase)
		if err != nil {
			log.Printf("Error: failed decoding pageState parameter %v\n", err)
			response.Error(w, "Invalid pageState parameter", http.StatusBadRequest)
			return
		}
	}

	limitInt := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limitInt = parsed
		}
	}

	chirps, nextPageState, err := h.service.GetChirpsByUserId(ctx, userId, ps, limitInt)
	if err != nil {
		log.Printf("Error: failed getting chirps by user id %v\n", err)
		response.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	nextStateStr := ""
	if len(nextPageState) > 0 {
		nextStateStr = base64.StdEncoding.EncodeToString(nextPageState)
	}

	response.JSON(w, http.StatusOK, struct {
		Chirps        []Chirp `json:"chirps"`
		NextPageState string  `json:"nextPageState,omitempty"`
	}{
		Chirps:        chirps,
		NextPageState: nextStateStr,
	})

}
