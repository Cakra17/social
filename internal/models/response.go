package models

type Response struct {
	Status string `json:"status"`
	Message string `json:"message,omitempty"`
	Data any `json:"data,omitempty"`
}