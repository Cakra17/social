package models

type Response struct {
	Success bool `json:"status"`
	Message string `json:"message,omitempty"`
	Data any `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool `json:"status"`
	Message string `json:"message"`
}