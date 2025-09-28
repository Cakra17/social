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

func (r *PostRepo) Update(ctx context.Context, post *models.Post) error {
	query := `
		UPDATE posts SET media = $1, caption = $2 WHERE id = $3 AND user_id = $4
	`
	_, err := r.db.ExecContext(ctx, query, post.Media, post.Caption, post.ID, post.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostRepo) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM posts WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
