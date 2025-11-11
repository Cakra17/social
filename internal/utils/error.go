package utils

import (
	"encoding/json"
	"net/http"

	. "github.com/cakra17/social/internal/models"
)

type NewError struct {
	Code    int
	Message string
}

var (
	ErrNoTokenProvided = NewError{Code: http.StatusUnauthorized, Message: "No token provided"}
	ErrTokenMalformed = NewError{Code: http.StatusUnauthorized, Message: "Token Malformed"}
	ErrTokenExpires = NewError{Code: http.StatusUnauthorized, Message: "Token expires, please login again"}
	ErrPayloadMalformed = NewError{Code: http.StatusBadRequest, Message: "Payload Malformed"}
	ErrFailedToCreateUser = NewError{Code: http.StatusInternalServerError, Message: "Failed to Create User"}
	ErrCredentialExist = NewError{Code: http.StatusConflict, Message: "Credentials already used"}
	ErrUserNotFound = NewError{Code: http.StatusNotFound, Message: "User not found"}
	ErrWrongPassword = NewError{Code: http.StatusBadRequest, Message: "Wrong password"}
	ErrInvalidUploadedFile = NewError{Code: http.StatusBadRequest, Message: "Invalid uploaded file"}
	ErrInvalidFileSize = NewError{Code: http.StatusBadRequest, Message: "Invalid file size, max 5mb"}
	ErrInvalidFileType = NewError{Code: http.StatusBadRequest, Message: "Invalid file type"}
  ErrInvalidPayload = NewError{Code: http.StatusBadRequest, Message: "Invalid Payload"}
	ErrFailedToUploadPhoto = NewError{Code: http.StatusInternalServerError, Message: "Failed to upload photo"}
  ErrFailedToCreatePost = NewError{Code: http.StatusInternalServerError, Message: "Failed to create post"}
)

func CreateNewError(code int, message string) NewError {
	return NewError{
		Code: code,
		Message: message,
	}
}

func JsonError(ne NewError, w http.ResponseWriter) {
	e := ErrorResponse{
		Message: ne.Message,
	}
	json, _ := getBytes(e)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(ne.Code)
	w.Write(json)
}

func getBytes(e any) ([]byte, error) {
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}
