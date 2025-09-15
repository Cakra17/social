package models

import (
	"time"
)

type User struct {
	ID        string 			`json:"id"`
	Username  string			`json:"username"`
	Email			string			`json:"email"`
	Password  string			`json:"-"`
	CreatedAt *time.Time	`json:"created_at"`
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

type UpdateUserPayload struct {
	Username  string			`json:"username" validate:"required"`
	Email			string			`json:"email" validate:"required,email,max=255"`
}

