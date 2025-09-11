package store

import (
	"context"
	"database/sql"

	"github.com/cakra17/social/internal/models"
	"github.com/lib/pq"
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
			id, content, title, 
			user_id,tags
		) VALUES (
			$1, $2, $3, $4, $5
		) RETURNING created_at, updated_at
	`
	err := r.db.QueryRowContext(
		ctx, query, 
		post.ID, 
		post.Content, 
		post.Title, 
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}