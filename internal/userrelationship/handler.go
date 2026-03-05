package userrelationship

import (
	"chirpstream/internal/auth"
	"chirpstream/pkg/response"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router *httprouter.Router) {
	router.POST("/api/user/follow/:followingId", auth.AuthMiddleware(h.handleFollowUser))
	router.POST("/api/user/unfollow/:followingId", auth.AuthMiddleware(h.handleUnfollowUser))
	router.GET("/api/user/following", auth.AuthMiddleware(h.handleGetFollowing))
	router.GET("/api/user/followers", auth.AuthMiddleware(h.handleGetFollowers))
}

func (h *Handler) handleFollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	followingId, err := strconv.Atoi(ps.ByName("followingId"))
	if err != nil {
		log.Printf("Error: failed parsing followingId parameter %v\n", err)
		response.Error(w, "Invalid followingId", http.StatusBadRequest)
		return
	}

	userId, ok := auth.GetUserIdFromContext(ctx)
	if !ok {
		log.Printf("Error: failed getting user from context\n")
		response.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if followingId == userId {
		log.Printf("Error: user cannot follow themselves\n")
		response.Error(w, "You cannot follow yourself", http.StatusBadRequest)
		return
	}

	err = h.service.FollowUser(ctx, userId, followingId)
	if err != nil {
		log.Printf("Error: failed following user %v\n", err)
		response.Error(w, "Failed to follow user", http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "User followed successfully"})
}

func (h *Handler) handleUnfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	followingId, err := strconv.Atoi(ps.ByName("followingId"))
	if err != nil {
		log.Printf("Error: failed parsing followingId parameter %v\n", err)
		response.Error(w, "Invalid followingId", http.StatusBadRequest)
		return
	}

	userId, ok := auth.GetUserIdFromContext(ctx)
	if !ok {
		log.Printf("Error: failed getting user from context\n")
		response.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = h.service.UnfollowUser(ctx, userId, followingId)
	if err != nil {
		log.Printf("Error: failed unfollowing user %v\n", err)
		response.Error(w, "Failed to unfollow user", http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "User unfollowed successfully"})
}

func (h *Handler) handleGetFollowing(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	userId, ok := auth.GetUserIdFromContext(ctx)
	if !ok {
		log.Printf("Error: failed getting user from context\n")
		response.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	following, err := h.service.GetFollowing(ctx, userId)
	if err != nil {
		log.Printf("Error: failed getting following %v\n", err)
		response.Error(w, "Failed to get following", http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusOK, following)
}

func (h *Handler) handleGetFollowers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	userId, ok := auth.GetUserIdFromContext(ctx)
	if !ok {
		log.Printf("Error: failed getting user from context\n")
		response.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	followers, err := h.service.GetFollowers(ctx, userId)
	if err != nil {
		log.Printf("Error: failed getting following %v\n", err)
		response.Error(w, "Failed to get following", http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusOK, followers)
}
