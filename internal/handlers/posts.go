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
	"github.com/cakra17/social/internal/utils"
	. "github.com/cakra17/social/internal/utils"
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
	logger *utils.Logger
}

type PostHandlerConfig struct {
	PostRepo store.PostRepo
	Logger *utils.Logger
}

func NewPostHandler(cfg PostHandlerConfig) PostHandler {
	return PostHandler{
		postRepo: cfg.PostRepo,
		logger: cfg.Logger,
	}
}

func generateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	timestamp := time.Now().Unix()
	hash := sha256.Sum256(fmt.Appendf(nil, "%s%d", ext, timestamp))
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
		JsonError(ErrInvalidFileSize, w)
		return
	}	

	media, header, err := r.FormFile("media")
	if err != nil {
		ne := CreateNewError(http.StatusBadRequest, "Failed retrieving data")
		JsonError(ne, w)
		return
	}
	defer media.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedType[ext] {
		JsonError(ErrInvalidFileType, w)
		return
	}

	filename := generateUniqueFilename(header.Filename)
	filepath := filepath.Join(UploadDir, filename)

	err = uploadPhoto(filepath, media)
	if err != nil {
		log.Println(err.Error())
		ne := CreateNewError(http.StatusInternalServerError, "failed to upload image")
		JsonError(ne, w)
		return
	}

	id, err := uuid.NewV7()
	if err != nil {
		JsonError(ErrFailedToCreatePost, w)
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
		log.Println(err.Error())
		JsonError(ErrFailedToCreatePost, w)
		return
	}

	res := models.Response{
		Success: true,
		Message: "Post Created successfully",
		Data: post,
	}

	jsonBytes, err := json.Marshal(&res)
	if err != nil {
    log.Println("Failed to encode to json")
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)
} 

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	caption := r.FormValue("caption")
	userID := r.FormValue("userID")

	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)
	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		JsonError(ErrInvalidFileSize, w)
		return
	}	

	media, header, err := r.FormFile("media")
	if err != nil {
		log.Println(err)
		ne := CreateNewError(http.StatusBadRequest, "Failed retrieving data")
		JsonError(ne, w)
		return
	}
	defer media.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedType[ext] {
		JsonError(ErrInvalidFileType, w)
		return
	}

	newFilename := generateUniqueFilename(header.Filename)
	newFilepath := filepath.Join(UploadDir, newFilename)

	ctx := r.Context()
	oldFilename, err := h.postRepo.GetPhoto(ctx, id)
	if err != nil {
		log.Println(err)
		ne := CreateNewError(http.StatusNotFound, "Post not found")
		JsonError(ne, w)
		return
	}

	oldFilepath := filepath.Join(UploadDir, oldFilename)

	post := &models.Post{
		ID: id,
		Caption: caption,
		Media: newFilename,
		UserID: userID,
	}

	if err := uploadPhoto(newFilepath, media); err != nil {
		ne := CreateNewError(http.StatusInternalServerError,"Failed upload photo")
		JsonError(ne, w)
		return
	}

	err = h.postRepo.Update(ctx, post)
	if err != nil {
		deletePhoto(newFilepath)

		log.Println(err.Error())
		ne := CreateNewError(http.StatusInternalServerError, "Failed to update post")
		JsonError(ne, w)
		return
	}

	if err := deletePhoto(oldFilepath); err != nil {
		ne := CreateNewError(http.StatusInternalServerError, "Failed to update post")
		JsonError(ne, w)
    return
	}

	res := models.Response{
		Success: true,
		Message: "Post updated successfully",
	}

	jsonBytes, err := json.Marshal(res)
	if err != nil {
		log.Println("Failed to encode to json")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ctx := r.Context()

	filename, err := h.postRepo.GetPhoto(ctx, id)
	if err != nil {
		log.Println(err)
		ne := CreateNewError(http.StatusNotFound, "Post not found")
		JsonError(ne, w)
		return
	}

	err = h.postRepo.Delete(r.Context(), id)
	if err != nil {
		log.Println(err.Error())
		ne := CreateNewError(http.StatusInternalServerError, "Failed to delete post")
		JsonError(ne, w)
		return
	}

	filepath := filepath.Join(UploadDir, filename)

	if err := deletePhoto(filepath); err != nil {
		log.Println(err.Error())
		ne := CreateNewError(http.StatusInternalServerError, "Failed to delete post")
		JsonError(ne, w)
		return
	}

	res := models.Response{
		Success: true,
		Message: "Data deleted successfully",
	}

	jsonBytes, err := json.Marshal(&res)
	if err != nil {
		log.Println("Failed to encode to json")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
