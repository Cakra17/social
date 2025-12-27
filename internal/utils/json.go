package utils

import (
	"encoding/json"
	"net/http"
)

type CustomError struct {
	Code    int
	Message string
}

type CustomSuccess struct {
	Code int
	Message string
	Data any
}

var (
	ErrNoTokenProvided = CustomError{Code: http.StatusUnauthorized, Message: "No token provided"}
	ErrTokenMalformed = CustomError{Code: http.StatusUnauthorized, Message: "Token Malformed"}
	ErrTokenNotContainsInfo = CustomError{Code: http.StatusUnauthorized, Message: "Bearer token not contains user info"}
	ErrTokenExpires = CustomError{Code: http.StatusUnauthorized, Message: "Token expires, please login again"}
	ErrPayloadMalformed = CustomError{Code: http.StatusBadRequest, Message: "Payload Malformed"}
	ErrFailedToCreateUser = CustomError{Code: http.StatusInternalServerError, Message: "Failed to Create User"}
	ErrCredentialExist = CustomError{Code: http.StatusConflict, Message: "Credentials already used"}
	ErrUserNotFound = CustomError{Code: http.StatusNotFound, Message: "User not found"}
	ErrWrongPassword = CustomError{Code: http.StatusBadRequest, Message: "Wrong password"}
	ErrInvalidUploadedFile = CustomError{Code: http.StatusBadRequest, Message: "Invalid uploaded file"}
	ErrInvalidFileSize = CustomError{Code: http.StatusBadRequest, Message: "Invalid file size, max 5mb"}
	ErrInvalidFileType = CustomError{Code: http.StatusBadRequest, Message: "Invalid file type"}
  ErrInvalidPayload = CustomError{Code: http.StatusBadRequest, Message: "Invalid Payload"}
	ErrFailedToUploadPhoto = CustomError{Code: http.StatusInternalServerError, Message: "Failed to upload photo"}
  ErrFailedToCreatePost = CustomError{Code: http.StatusInternalServerError, Message: "Failed to create post"}
)

type Response struct {
	Message string `json:"message,omitempty"`
	Data any `json:"data,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func ParseBody(r *http.Request, payload any) error {
	return json.NewDecoder(r.Body).Decode(&payload)
}

func WriteJson(w http.ResponseWriter, successResponse CustomSuccess) {
	var res Response
	if successResponse.Message == "" {
		res = Response{
			Data: successResponse.Data,
		}
	} else {
		res = Response{
			Message: successResponse.Message,
			Data: successResponse.Data,
		}
	}

	resByte, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(successResponse.Code)
	w.Write(resByte)
}

func WriteError(w http.ResponseWriter, errorResponse CustomError) {
	errRes := ErrorResponse{
		Message: errorResponse.Message,
	}

	errBytes, _ := json.Marshal(errRes)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errorResponse.Code)
	w.Write(errBytes)
}