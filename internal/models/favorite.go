package models

import "time"

type Favorite struct {
	ID        string `json:"id" db:"id"`
	PostId    string `json:"post_id" db:"post_id"`
	UserId    string `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}