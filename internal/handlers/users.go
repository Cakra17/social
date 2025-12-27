package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cakra17/social/internal/models"
	"github.com/cakra17/social/internal/store"
	"github.com/cakra17/social/internal/utils"
	. "github.com/cakra17/social/internal/utils"
	"github.com/cakra17/social/pkg/jwt"
	"github.com/cakra17/social/pkg/validation"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type UserHandler struct {
	userRepo store.UserRepo
	redis *redis.Client
	jwtAuthenticator *jwt.JWTAuthenticator
	logger *utils.Logger
}

type UserHandlerConfig struct {
	UserRepo store.UserRepo
	Redis *redis.Client
	JWTAuthenticator *jwt.JWTAuthenticator
	Logger *utils.Logger
}

func NewUserHandler(cfg UserHandlerConfig) UserHandler {
	return UserHandler{
		userRepo: cfg.UserRepo,
		redis: cfg.Redis,
		jwtAuthenticator: cfg.JWTAuthenticator,
		logger: cfg.Logger,
	}
}

func(h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var payload models.RegisterPayload

	if err := utils.ParseBody(r, &payload); err != nil {
		h.logger.Error("User Handler Error", "Failed to decode payload", err.Error())
		WriteError(w, ErrPayloadMalformed)
		return
	}

	if err := validation.Validate(&payload); err != nil {
		h.logger.Error("User Handler Error", "Failed to validate payload", err)
		WriteError(w, ErrInvalidPayload)
		return
	}

	ctx := r.Context()

	user, _ := h.userRepo.GetUserByEmail(ctx, payload.Email)
	if user != nil {
		h.logger.Error("User Handler Error", "Failed to create user", "Credentials already used")
		WriteError(w, ErrCredentialExist)
		return
	}

	id, err := uuid.NewV7()
	if err != nil {
		h.logger.Error("User Handler Error", "Failed to create user", "Can't generate uuid")
		WriteError(w, ErrFailedToCreateUser)
		return
	}

	hashedPassword, err := HashPassword(payload.Password)
	if err != nil {
		h.logger.Error("User Handler Error", "Failed to create user", "Can't hash password")
		WriteError(w, ErrFailedToCreateUser)
		return
	}

	user = &models.User{
		ID: id.String(),
		Username: payload.Username,
		Email: payload.Email,
		Password: hashedPassword,
	}
	
	err = h.userRepo.CreateUser(ctx, user)
	if err != nil {
		h.logger.Error("User Handler Error", "Failed to create user", err.Error())
		WriteError(w, ErrFailedToCreateUser)
		return
	}

	WriteJson(w, CustomSuccess{
		Code: http.StatusCreated,
		Data: user,
	})
}

func(h *UserHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	var payload models.LoginPayload

	if err := utils.ParseBody(r, &payload); err != nil {
		h.logger.Error("User Handler Error", "Failed to decode payload", err.Error())
		WriteError(w, ErrPayloadMalformed)
		return
	}

	if err := validation.Validate(&payload); err != nil {
		h.logger.Error("User Handler Error", "Failed to validate payload", err)
		WriteError(w, ErrInvalidPayload)
		return
	}

	ctx := r.Context()
	user, err := h.userRepo.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		h.logger.Error("User Handler Error", "Failed to login", "Account doesn't exist")
		WriteError(w, ErrUserNotFound)
		return
	}

	h.redis.Set(ctx, user.ID, user, 30 * time.Second)

	if ok := ComparePassword(payload.Password, user.Password); !ok {
		h.logger.Error("User Handler Error", "Failed to login", "Wrong password")
		WriteError(w, ErrWrongPassword)
		return
	}

	token, err := h.jwtAuthenticator.GenerateToken(jwt.JWTUser{
		ID: user.ID,
		Email: user.Email,
	})

	if err != nil {	
		h.logger.Error("User Handler Error", "Failed to generate token", err.Error())
		WriteError(w, CustomError{
			Code: http.StatusInternalServerError, 
			Message: err.Error(),
		})
		return
	}

	WriteJson(w, CustomSuccess{
		Code: http.StatusOK,
		Message: "success to login",
		Data: models.AuthResponse{
			AccessToken: token,
		},
	})
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	var user *models.User
	
	ctx := r.Context()
	claims, ok := h.jwtAuthenticator.GetClaims(ctx)
	if !ok {
		h.logger.Error("User Handler Error", "Failed get claims")
		WriteError(w, ErrTokenExpires)
		return
	}

	userID, _ := claims["userId"].(string)

	s, err := h.redis.Get(ctx, userID).Result()
	if err != nil {
		user, err = h.userRepo.GetUserById(ctx, userID)
		if err != nil {
			WriteError(w, ErrUserNotFound)
			return
		}
		err :=h.redis.Set(ctx, userID, user, time.Minute).Err()
		if err != nil {
			h.logger.Error("User Handler Error", "Failed to save in redis", err.Error())
		}
	} else {
		if err := json.Unmarshal([]byte(s), &user); err != nil {
			return
		}
	}

	WriteJson(w, CustomSuccess{
		Code: http.StatusOK,
		Data: user,
	})
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var payload models.UpdateUserPayload

	if err := utils.ParseBody(r, &payload); err != nil {
		h.logger.Error("User Handler Error", "Failed to decode payload", err.Error())
		WriteError(w, ErrPayloadMalformed)
		return
	}

	if err := validation.Validate(&payload); err != nil {
		h.logger.Error("User Handler Error", "Failed to validate payload", err)
		WriteError(w, ErrInvalidPayload)
		return
	}

	id := r.PathValue("id")

	ctx := r.Context()
	err := h.userRepo.UpdateUser(ctx, &payload, id)
	if err != nil {
		h.logger.Error("User Handler Error", "Failed to update user", err.Error())
		WriteError(w, CustomError{
			Code :http.StatusInternalServerError, 
			Message: "Failed to update user",
		})
		return
	}	

	WriteJson(w, CustomSuccess{
		Code: http.StatusOK,
		Message: "Data Updated successfully",
	})
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ctx := r.Context()
	err := h.userRepo.Delete(ctx, id)
	if err != nil {
		h.logger.Error("User Handler Error", "Failed to delete user", err.Error())
		WriteError(w, CustomError{
			Code: http.StatusInternalServerError,
			Message: "Failed to delete user",
		})
		return
	}	

	WriteJson(w, CustomSuccess{
		Code: http.StatusOK,
		Message: "Data deleted successfully",
	})
}
