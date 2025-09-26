package models

import "time"

type Post struct {
	ID        string    `json:"id"`
	Caption   string    `json:"caption"`
	Media 		string		`json:"Media"`
	UserID    string    `json:"user_id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
