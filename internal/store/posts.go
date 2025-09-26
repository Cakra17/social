package store

import (
	"context"
	"database/sql"

	"github.com/cakra17/social/internal/models"
)

type PostRepo struct {
	db *sql.DB
}

func NewPostRepo(db *sql.DB) PostRepo {
	return PostRepo{ db: db }
}

func (r *PostRepo) Create(ctx context.Context, post *models.Post) error {
	query := `
		INSERT INTO posts (
			id, caption, media, 
			user_id
		) VALUES (
			$1, $2, $3, $4
		) RETURNING created_at, updated_at
	`
	err := r.db.QueryRowContext(
		ctx, query, 
		post.ID, 
		post.Media, 
		post.UserID,
	).Scan(
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}
