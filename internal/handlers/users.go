package handlers

import (
	"net/http"

	"github.com/cakra17/social/internal/store"
)

type UserHandler struct {
	userRepo store.UserRepo
}

type UserHandlerConfig struct {
	UserRepo store.UserRepo
}

func NewUserHandler(cfg UserHandlerConfig) UserHandler {
	return UserHandler{
		userRepo: cfg.UserRepo,
	}
}

func(h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hellow"))
}