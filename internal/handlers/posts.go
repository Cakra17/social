package handlers

import (
	"net/http"

	"github.com/cakra17/social/internal/store"
)

type PostHandler struct {
	postRepo store.PostRepo
}

type PostHandlerConfig struct {
	PostRepo store.PostRepo
}

func NewPostHandler(cfg PostHandlerConfig) PostHandler {
	return PostHandler{
		postRepo: cfg.PostRepo,
	}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {

} 