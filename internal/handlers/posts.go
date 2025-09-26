package handlers

import (
	// "encoding/json"
	"encoding/base64"
	"io"
	"log"
	"net/http"

	"github.com/cakra17/social/internal/models"
	"github.com/cakra17/social/internal/store"
	"github.com/google/uuid"
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
	caption := r.FormValue("caption")
	userid := r.FormValue("userID")

	// max 10 MB
	r.ParseMultipartForm(10 << 20)

	media, _, err := r.FormFile("media")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer media.Close()

	fileBytes, err := io.ReadAll(media)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encodedFile := base64.StdEncoding.EncodeToString(fileBytes)

	id, err := uuid.NewV7()
	if err != nil {
		log.Println("Failed to generate id")
		return
	}
	
	post := &models.Post{
		ID: id.String(),
		Caption: caption,
		Media: encodedFile,
		UserID: userid,
	}

	err = h.postRepo.Create(r.Context(), post)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
} 