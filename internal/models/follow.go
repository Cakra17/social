package models

type Follower struct {
	ID string `json:"follow_id" db:"id"`
	UserID string `json:"user_id" db:"user_id"`
	Username string `json:"username" db:"username"`
}

type FollowPayload struct {
	FolloweeID string	`json:"followee_id" validate:"required"`
	FollowerID string	`json:"follower_id" validate:"required"`
}

type Follow struct {
	ID string
	FolloweeID string	`db:"followee_id"`
	FollowerID string	`db:"followers_id"`
}