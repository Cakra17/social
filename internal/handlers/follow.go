package handlers

import (
	"net/http"

	"github.com/cakra17/social/internal/models"
	"github.com/cakra17/social/internal/store"
	"github.com/cakra17/social/internal/utils"
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
	logger *utils.Logger
}

type FollowHandlerConfig struct {
	FollowRepo store.FollowRepo
	Redis *redis.Client
	JWTAuthenticator *jwt.JWTAuthenticator
	Logger *utils.Logger
}

func NewFollowHandler(cfg FollowHandlerConfig) FollowHandler {
	return FollowHandler{
		followRepo: cfg.FollowRepo,
		redis: cfg.Redis,
		jwtAuthenticator: cfg.JWTAuthenticator,
		logger: cfg.Logger,
	}
}

func (h *FollowHandler) Follow(w http.ResponseWriter, r *http.Request) {
	var payload models.FollowPayload

	if err := utils.ParseBody(r, &payload); err != nil {
		WriteError(w, ErrPayloadMalformed)
		return
	}

	if err := validation.Validate(&payload); err != nil {
    WriteError(w, ErrInvalidPayload)
		return
	}

	ctx := r.Context()

	id, err := uuid.NewV7()
	if err != nil {
		h.logger.Error("Favorite Handler Error", "Failed to create id", err.Error())
		WriteError(w, CustomError{
			Code: http.StatusInternalServerError,
			Message: "Failed to follow user",
		})
		return
	}

	follow := models.Follow{
		ID: id.String(),
		FolloweeID: payload.FolloweeID,
		FollowerID: payload.FollowerID,
	}

	err = h.followRepo.Follow(ctx, follow)
	if err != nil {
		h.logger.Error("Favorite Handler Error", "Failed to follow user", err.Error())
		WriteError(w, CustomError{
			Code: http.StatusInternalServerError,
			Message: "Failed to follow user",
		})
		return
	}

	WriteJson(w, CustomSuccess{
		Code: http.StatusCreated,
		Message: "started to follow",
		Data: follow,
	})
}

func (h *FollowHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := h.jwtAuthenticator.GetClaims(ctx)
	if !ok {
    WriteError(w, ErrTokenExpires)
		return
	}

	userId, _ := claims["userId"].(string)

	followers, err := h.followRepo.GetFollowers(ctx, userId)
	if err != nil {
		WriteError(w, CustomError{
			Code: http.StatusInternalServerError,
			Message: "Failed to get followers",
		})
		return
	}

	WriteJson(w, CustomSuccess{
		Code: http.StatusOK,
		Data: followers,
	})
}

func (h *FollowHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := h.jwtAuthenticator.GetClaims(ctx)
	if !ok {
    WriteError(w,ErrTokenExpires)
		return
	}

	userId, _ := claims["userId"].(string)

	following, err := h.followRepo.GetFollowing(ctx, userId)
	if err != nil {
    WriteError(w, CustomError{
			Code: http.StatusInternalServerError, 
			Message: "Failed to get followers",
		})
		return
	}

	WriteJson(w, CustomSuccess{
		Code: http.StatusOK,
		Data: following,
	})
}

func (h *FollowHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ctx := r.Context()
	err := h.followRepo.Unfollow(ctx, id)
	if err != nil {
		WriteError(w, CustomError{
			Code: http.StatusInternalServerError, 
			Message: "Failed to unfollow",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
