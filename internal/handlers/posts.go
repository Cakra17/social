package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cakra17/social/internal/models"
	"github.com/cakra17/social/internal/store"
	"github.com/google/uuid"
)

const (
	MaxUploadSize = 10 << 20
	UploadDir = "./uploads"
)

var allowedType map[string]bool = map[string]bool{
	".jpg": true,
	".jpeg": true,
	".png": true,
}

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

func generateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	timestamp := time.Now().Unix()
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s%d", ext, timestamp)))
	hashStr := hex.EncodeToString(hash[:])[:16]
	return fmt.Sprintf("%d_%s%s", timestamp, hashStr, ext)
}

func deletePhoto(filepath string) error {
	if err := os.Remove(filepath); err != nil {
		return err
	}
	return  nil
}

func uploadPhoto(filepath string, media multipart.File) error {
	dst, err := os.Create(filepath)
	if err != nil {
		log.Println("Failed to save a file")
		return fmt.Errorf("Failed to save a file")
	}
	defer dst.Close()

	_, err = io.Copy(dst, media)
	if err != nil {
		log.Println("Failed to save a file")
		return fmt.Errorf("Failed to save a file")
	}
	return nil
}

func (h *PostHandler) Init() {
	if err := os.MkdirAll(UploadDir, 0755); err != nil {
		log.Printf("Failed to create directory: %s", err.Error())
		return
	}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	caption := r.FormValue("caption")
	userid := r.FormValue("userID")

	// limit file photo size
	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)
	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}	

	media, header, err := r.FormFile("media")
	if err != nil {
		http.Error(w, "Failed retrieving data", http.StatusInternalServerError)
		return
	}
	defer media.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedType[ext] {
		http.Error(w, "File tyoe not allowed", http.StatusBadRequest)
	}

	filename := generateUniqueFilename(header.Filename)
	filepath := filepath.Join(UploadDir, filename)

	err = uploadPhoto(filepath, media)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := uuid.NewV7()
	if err != nil {
		log.Println("Failed to generate id")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	post := &models.Post{
		ID: id.String(),
		Caption: caption,
		Media: filename,
		UserID: userid,
	}

	ctx :=  r.Context()

	err = h.postRepo.Create(ctx, post)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := models.Response{
		Status: "success",
		Message: "Post Created successfully",
		Data: post,
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

	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)
	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}	

	media, header, err := r.FormFile("media")
	if err != nil {
		http.Error(w, "Failed retrieving data", http.StatusInternalServerError)
		return
	}
	defer media.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedType[ext] {
		http.Error(w, "File tyoe not allowed", http.StatusBadRequest)
	}

	filename := generateUniqueFilename(header.Filename)
	filepath := filepath.Join(UploadDir, filename)

	err = uploadPhoto(filepath, media)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	filepath, err = h.postRepo.GetPhoto(ctx, id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Post is not available", http.StatusBadRequest)
		return
	}

	if err := deletePhoto(filepath); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	post := &models.Post{
		ID: id,
		Caption: caption,
		Media: filepath,
	}

	err = h.postRepo.Update(r.Context(), post)
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

	ctx := r.Context()

	filepath, err := h.postRepo.GetPhoto(ctx, id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Post is not available", http.StatusBadRequest)
		return
	}

	err = h.postRepo.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := deletePhoto(filepath); err != nil {
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