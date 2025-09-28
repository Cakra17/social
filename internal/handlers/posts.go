package handlers

import (
	// "encoding/json"
	"encoding/base64"
	"encoding/json"
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

	res := models.Response{
		Status: "success",
		Message: "Post Created successfully",
	}

	jsonBytes, err := json.Marshal(&res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)
} 

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	caption := r.FormValue("caption")
	media := r.FormValue("media")

	post := &models.Post{
		ID: id,
		Caption: caption,
		Media: media,
	}

	err := h.postRepo.Update(r.Context(), post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := models.Response{
		Status: "success",
		Message: "Post updated successfully",
	}

	jsonBytes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := h.postRepo.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := models.Response{
		Status: "success",
		Message: "Data deleted successfully",
	}

	jsonBytes, err := json.Marshal(&res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}