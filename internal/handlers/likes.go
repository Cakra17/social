package handlers

import (
	"net/http"

	"github.com/cakra17/social/internal/models"
	"github.com/cakra17/social/internal/store"
	"github.com/cakra17/social/internal/utils"
	"github.com/cakra17/social/pkg/jwt"
	"github.com/google/uuid"
)

type LikesHandler struct {
	likesRepo store.LikesRepo
	jwtAuthenticator *jwt.JWTAuthenticator
	logger *utils.Logger
}

type LikesHandlerConfig struct {
	LikesRepo store.LikesRepo
	JWTAuthenticator *jwt.JWTAuthenticator
	Logger *utils.Logger
}

func NewLikesHandler(cfg LikesHandlerConfig) LikesHandler {
	return LikesHandler{
		likesRepo: cfg.LikesRepo,
		jwtAuthenticator: cfg.JWTAuthenticator,
		logger: cfg.Logger,
	}
}

func (h *LikesHandler) Like(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims, ok := h.jwtAuthenticator.GetClaims(ctx)
	if !ok {
		h.logger.Error("Authetication Error", "claim not found")
		return
	}

	userID := claims["userId"].(string)
	postID := r.PathValue("postId")

	likes := models.Likes{
		ID: uuid.Must(uuid.NewV7()).String(),
		PostId: postID,
		UserId: userID,
	}

	err := h.likesRepo.Like(ctx, likes)
	if err != nil {
		h.logger.Error("Like Handler Error", "Failed to liked post", err.Error())
		utils.WriteError(w, utils.CustomError{
			Code: http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	utils.WriteJson(w, utils.CustomSuccess{
		Code: http.StatusCreated,
		Data: likes,
	})
}

func (h *LikesHandler) GetPostLikes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	postID := r.PathValue("postId")

	likes, err := h.likesRepo.GetLikes(ctx, postID)
	if err != nil {
		h.logger.Error("Like Handler Error", "Failed to get liked post", err.Error())
		utils.WriteError(w, utils.CustomError{
			Code: http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	utils.WriteJson(w, utils.CustomSuccess{
		Code: http.StatusOK,
		Message: "Succes to get likes data",
		Data: likes.Users,
	})
}

func (h *LikesHandler) Unlike(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	postID := r.PathValue("likesId")

	err := h.likesRepo.Unlike(ctx, postID)
	if err != nil {
		utils.WriteError(w, utils.CustomError{
			Code: http.StatusInternalServerError,
			Message: err.Error(),
		})
		h.logger.Error("Like Handler Error", "Failed to unlike post", err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}