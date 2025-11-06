package models

import (
	"encoding/json"
	"time"
)

type User struct {
	ID        string 			`json:"id"`
	Username  string			`json:"username"`
	Email			string			`json:"email"`
	Password  string			`json:"-"`
	CreatedAt *time.Time	`json:"created_at"`
}

func (u User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

type RegisterPayload struct {
	Username  string			`json:"username" validate:"required"`
	Email			string			`json:"email" validate:"required,email,max=255"`
	Password  string			`json:"password" validate:"required,min=8,max=30"`
}

type LoginPayload struct {
	Email			string			`json:"email" validate:"required,email,max=255"`
	Password  string			`json:"password" validate:"required,min=8,max=30"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
}

type UpdateUserPayload struct {
	Username  string			`json:"username" validate:"required"`
	Email			string			`json:"email" validate:"required,email,max=255"`
}

