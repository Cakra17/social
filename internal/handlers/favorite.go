package handlers

import (
	"net/http"

	"github.com/cakra17/social/internal/models"
	"github.com/cakra17/social/internal/store"
	"github.com/cakra17/social/internal/utils"
	"github.com/cakra17/social/pkg/jwt"
	"github.com/google/uuid"
)

type FavoriteHandler struct {
	favoriteRepo store.FavoriteRepo
	logger *utils.Logger
	jwtAuthenticator *jwt.JWTAuthenticator
}

type FavoriteHandlerConfig struct {
	FavoriteRepo store.FavoriteRepo
	Logger *utils.Logger
	JWTAuthenticator *jwt.JWTAuthenticator
}

func NewFavoriteHandler(cfg FavoriteHandlerConfig) FavoriteHandler {
	return FavoriteHandler{
		favoriteRepo: cfg.FavoriteRepo,
		logger: cfg.Logger,
		jwtAuthenticator: cfg.JWTAuthenticator,
	}
}

func (h *FavoriteHandler) AddFavorite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims, ok := h.jwtAuthenticator.GetClaims(ctx)
	if !ok {
		h.logger.Error("Authetication Error", "claim not found")
		return
	}

	userID := claims["userId"].(string)
	postID := r.PathValue("postId")

	favorite := models.Favorite{
		ID: uuid.Must(uuid.NewV7()).String(),
		PostId: postID,
		UserId: userID,
	}

	err := h.favoriteRepo.Add(ctx, &favorite)
	if err != nil {
		h.logger.Error("Favorite Handler Error", "Failed add post to favorite", err.Error())
		utils.WriteError(w, utils.CustomError{
			Code: http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	utils.WriteJson(w, utils.CustomSuccess{
		Code: http.StatusCreated,
		Message: "Added to your favorite",
		Data: favorite,
	})
} 

func (h *FavoriteHandler) GetFavouritePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims, ok := h.jwtAuthenticator.GetClaims(ctx)
	if !ok {
		h.logger.Error("Authetication Error", "claim not found")
		return
	}

	userID := claims["userId"].(string)

	favorites, err := h.favoriteRepo.GetFavouritePost(ctx, userID)
	if err != nil {
		h.logger.Error("Favorite Handler Error", "Failed to get favorite", err.Error())
		utils.WriteError(w, utils.CustomError{
			Code: http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	utils.WriteJson(w, utils.CustomSuccess{
		Code: http.StatusOK,
		Message: "Succes to getlikes data",
		Data: favorites,
	})
}

func (h *FavoriteHandler) DeleteFavorite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	postID := r.PathValue("postId")

	err := h.favoriteRepo.Delete(ctx, postID)
	if err != nil {
		h.logger.Error("Favorite Error", "Failed to delete favorite", err.Error())
		utils.WriteError(w, utils.CustomError{
			Code: http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
