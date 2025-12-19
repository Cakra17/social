package models

type Likes struct {
	ID 			string 	`json:"id" db:"id"`
	Users		[]User	`json:"users"`
	LikesCount int `json:"likes_count"`
	PostId 	string 	`json:"post_id" db:"post_id"`
	UserId 	string 	`json:"user_id" db:"user_id"`
}

