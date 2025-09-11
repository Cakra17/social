package models

import "time"

type Post struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    string    `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
