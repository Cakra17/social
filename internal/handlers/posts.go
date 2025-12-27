package handlers

import (
	"crypto/sha256"
	"encoding/hex"
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
		h.logger.Error("Post Handler Error", "Failed to retrive data", err.Error())
		WriteError(w, ErrInvalidFileSize)
		return
	}	

	media, header, err := r.FormFile("media")
	if err != nil {
		h.logger.Error("Post Handler Error", "Failed to retrive data", err.Error())
		WriteError(w, CustomError{
			Code: http.StatusBadRequest,
			Message: "Failed retrieving data",
		})
		return
	}
	defer media.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedType[ext] {
		h.logger.Error("Post Handler Error", "Failed to retrive data", "Invalid file type")
		WriteError(w, ErrInvalidFileType)
		return
	}

	filename := generateUniqueFilename(header.Filename)
	filepath := filepath.Join(UploadDir, filename)

	err = uploadPhoto(filepath, media)
	if err != nil {
		h.logger.Error("Post Handler Error", "Failed to Upload", err.Error())
		WriteError(w, CustomError{
			Code: http.StatusInternalServerError,
			Message: "failed to upload image",
		})
		return
	}

	id, err := uuid.NewV7()
	if err != nil {
		h.logger.Error("Post Handler Error", "Failed to create id", err.Error())
		WriteError(w, ErrFailedToCreatePost)
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
		h.logger.Error("Post Handler Error", "Failed to create post", err.Error())
		log.Println(err.Error())
		WriteError(w, ErrFailedToCreatePost)
		return
	}

	WriteJson(w, CustomSuccess{
		Code: http.StatusCreated,
		Message: "Post Created successfully",
		Data: post,
	})
} 

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	caption := r.FormValue("caption")
	userID := r.FormValue("userID")

	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)
	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		h.logger.Error("Post Handler Error", "Failed to retrive data", err.Error())
		WriteError(w, ErrInvalidFileSize)
		return
	}	

	media, header, err := r.FormFile("media")
	if err != nil {
		h.logger.Error("Post Handler Error", "Failed to retrive data", err.Error())
		WriteError(w, CustomError{
			Code:http.StatusBadRequest, 
			Message: "Failed retrieving data",
		})
		return
	}
	defer media.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedType[ext] {
		h.logger.Error("Post Handler Error", "Failed to retrive data", "Invalid file type")
		WriteError(w, ErrInvalidFileType)
		return
	}

	newFilename := generateUniqueFilename(header.Filename)
	newFilepath := filepath.Join(UploadDir, newFilename)

	ctx := r.Context()
	oldFilename, err := h.postRepo.GetPhoto(ctx, id)
	if err != nil {
		h.logger.Error("Post Handler Error", "Failed to get expected post", err.Error())
		WriteError(w, CustomError{
			Code:http.StatusNotFound, 
			Message: "Post not found",
		})
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
		h.logger.Error("Post Handler Error", "Failed to upload", err.Error())
		WriteError(w, CustomError{
			Code:http.StatusInternalServerError,
			Message: "Failed upload photo",
		})
		return
	}

	err = h.postRepo.Update(ctx, post)
	if err != nil {
		deletePhoto(newFilepath)
		h.logger.Error("Post Handler Error", "Failed to update post", err.Error())
		WriteError(w, CustomError{
			Code: http.StatusInternalServerError,
			Message: "Failed to update post",
		})
		return
	}

	if err := deletePhoto(oldFilepath); err != nil {
		h.logger.Error("Post Handler Error", "Failed to delete old photo", err.Error())
		WriteError(w, CustomError{
			Code:http.StatusInternalServerError,
			Message: "Failed to update post",
		})
    return
	}

	WriteJson(w, CustomSuccess{
		Code: http.StatusOK,
		Message: "Post updated successfully",
	})
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ctx := r.Context()

	filename, err := h.postRepo.GetPhoto(ctx, id)
	if err != nil {
		h.logger.Error("Post Handler Error", "Failed to get expected post", err.Error())
		WriteError(w, CustomError{
			Code: http.StatusNotFound, 
			Message: "Post not found",
		})
		return
	}

	err = h.postRepo.Delete(r.Context(), id)
	if err != nil {
		h.logger.Error("Post Handler Error", "Failed to delete post", err.Error())
		WriteError(w, CustomError{
			Code: http.StatusInternalServerError,
			Message: "Failed to delete post",
		})
		return
	}

	filepath := filepath.Join(UploadDir, filename)

	if err := deletePhoto(filepath); err != nil {
		h.logger.Error("Post Handler Error", "Failed to delete photo", err.Error())
		WriteError(w, CustomError{
			Code: http.StatusInternalServerError, 
			Message: "Failed to delete post",
		})
		return
	}

	WriteJson(w, CustomSuccess{
		Code: http.StatusOK,
		Message: "Data deleted successfully",
	})
}