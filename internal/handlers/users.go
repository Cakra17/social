package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/cakra17/social/internal/models"
	"github.com/cakra17/social/internal/store"
	"github.com/cakra17/social/internal/utils"
	"github.com/cakra17/social/pkg/jwt"
	"github.com/cakra17/social/pkg/validation"
	"github.com/google/uuid"
)

type UserHandler struct {
	userRepo store.UserRepo
	jwtAuthenticator *jwt.JWTAuthenticator
}

type UserHandlerConfig struct {
	UserRepo store.UserRepo
	JWTAuthenticator *jwt.JWTAuthenticator
}

func NewUserHandler(cfg UserHandlerConfig) UserHandler {
	return UserHandler{
		userRepo: cfg.UserRepo,
		jwtAuthenticator: cfg.JWTAuthenticator,
	}
}

func(h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var payload models.RegisterPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := validation.Validate(&payload); err != nil {
		log.Print(err)
		http.Error(w, "The payload is not valid", http.StatusInternalServerError)
		return
	}

	id, err := uuid.NewV7()
	if err != nil {
		log.Println("Failed to generate id")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := &models.User{
		ID: id.String(),
		Username: payload.Username,
		Email: payload.Email,
		Password: hashedPassword,
	}

	ctx := r.Context()

	err = h.userRepo.CreateUser(ctx, user)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := models.Response{
		Status: "success",
		Data: user,
	}

	jsonBytes, err := json.Marshal(res)
	if err != nil {
		log.Println("Failed to encode to json")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)
}

func(h *UserHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	var payload models.LoginPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := validation.Validate(&payload); err != nil {
		log.Print(err)
		http.Error(w, "The payload is not valid", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	user, err := h.userRepo.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if ok := utils.ComparePassword(payload.Password, user.Password); !ok {
		log.Println(err)
		http.Error(w, "Password Incorrect", http.StatusInternalServerError)
		return
	}

	token, err := h.jwtAuthenticator.GenerateToken(jwt.JWTUser{
		ID: user.ID,
		Email: user.Email,
	}, 5 * time.Hour)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := models.Response{
		Status: "success",
		Message: "success to login",
		Data: models.AuthResponse{
			AccessToken: token,
		},
	}

	jsonBytes, err := json.Marshal(res)
	if err != nil {
		log.Println("Failed to encode to json")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var payload models.UpdateUserPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := validation.Validate(&payload); err != nil {
		log.Print(err)
		http.Error(w, "The payload is not valid", http.StatusInternalServerError)
		return
	}

	id := r.PathValue("id")

	ctx := r.Context()
	err := h.userRepo.UpdateUser(ctx, &payload, id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}	

	res := models.Response{
		Status: "success",
		Message: "Data Updated successfully",
	}

	jsonBytes, err := json.Marshal(res)
	if err != nil {
		log.Println("Failed to encode to json")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ctx := r.Context()
	err := h.userRepo.Delete(ctx, id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}	

	res := models.Response{
		Status: "success",
		Message: "Data deleted successfully",
	}

	jsonBytes, err := json.Marshal(res)
	if err != nil {
		log.Println("Failed to encode to json")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

