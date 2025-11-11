package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cakra17/social/internal/models"
	"github.com/cakra17/social/internal/store"
	. "github.com/cakra17/social/internal/utils"
	"github.com/cakra17/social/pkg/jwt"
	"github.com/cakra17/social/pkg/validation"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type FollowHandler struct {
	followRepo store.FollowRepo
	redis *redis.Client
	jwtAuthenticator *jwt.JWTAuthenticator
}

type FollowHandlerConfig struct {
	FollowRepo store.FollowRepo
	Redis *redis.Client
	JWTAuthenticator *jwt.JWTAuthenticator
}

func NewFollowHandler(cfg FollowHandlerConfig) FollowHandler {
	return FollowHandler{
		followRepo: cfg.FollowRepo,
		redis: cfg.Redis,
		jwtAuthenticator: cfg.JWTAuthenticator,
	}
}

func (h *FollowHandler) Follow(w http.ResponseWriter, r *http.Request) {
	var payload models.FollowPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		JsonError(ErrPayloadMalformed, w)
		return
	}

	if err := validation.Validate(&payload); err != nil {
    JsonError(ErrInvalidPayload, w)
		return
	}

	ctx := r.Context()

	id, err := uuid.NewV7()
	if err != nil {
    ne := CreateNewError(http.StatusInternalServerError, "Failed to follow user")
    JsonError(ne, w)
		return
	}

	follow := models.Follow{
		ID: id.String(),
		FolloweeID: payload.FolloweeID,
		FollowerID: payload.FollowerID,
	}

	err = h.followRepo.Follow(ctx, follow)
	if err != nil {
    ne := CreateNewError(http.StatusInternalServerError, "Failed to follow user")
    JsonError(ne, w)
		return
	}

	res := models.Response{
		Status: "success",
		Message: "started to follow",
	}

	jsonBytes, err := json.Marshal(&res)
	if err != nil {
    log.Println(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *FollowHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := h.jwtAuthenticator.GetClaims(ctx)
	if !ok {
    JsonError(ErrTokenExpires, w)
		return
	}

	userId, _ := claims["userId"].(string)

	followers, err := h.followRepo.GetFollowers(ctx, userId)
	if err != nil {
    ne := CreateNewError(http.StatusInternalServerError, "Failed to get followers")
    JsonError(ne,w)
		return
	}

	res := models.Response {
		Status: "success",
		Data: followers,
	}

	jsonBytes, err := json.Marshal(&res)
	if err != nil {
    log.Println("Failed to encode to json")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *FollowHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := h.jwtAuthenticator.GetClaims(ctx)
	if !ok {
    JsonError(ErrTokenExpires, w)
		return
	}

	userId, _ := claims["userId"].(string)

	following, err := h.followRepo.GetFollowing(ctx, userId)
	if err != nil {
    ne := CreateNewError(http.StatusInternalServerError, "Failed to get followers")
    JsonError(ne,w)
		return
	}

	res := models.Response {
		Status: "success",
		Data: following,
	}

	jsonBytes, err := json.Marshal(&res)
	if err != nil {
    log.Println("Failed to encode to json")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *FollowHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ctx := r.Context()
	err := h.followRepo.Unfollow(ctx, id)
	if err != nil {
    ne := CreateNewError(http.StatusInternalServerError, "Failed to unfollow")
    JsonError(ne, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
