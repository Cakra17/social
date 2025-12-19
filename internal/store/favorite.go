package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cakra17/social/internal/models"
	"github.com/cakra17/social/internal/utils"
)

type FavoriteRepo struct {
	db *sql.DB
	logger *utils.Logger
}

func NewFavoriteRepo(db *sql.DB, lg *utils.Logger) FavoriteRepo {
	return FavoriteRepo{
		db: db,
		logger: lg,
	}
}

func (r *FavoriteRepo) Add(ctx context.Context, payload *models.Favorite) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("Failed to begin transaction: %s", err.Error())
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); err != nil {
				r.logger.Error("Database Error", "Failed to rollback", rbErr.Error())
			}
		}
	}()

	query := `
		INSERT INTO favorites (id, post_id, user_id) VALUES ($1, $2, $3) RETURNING created_at
	`
	err = r.db.QueryRowContext(
		ctx, query, payload.ID, payload.PostId, payload.UserId,
	).Scan(&payload.CreatedAt)
	if err != nil {
		return fmt.Errorf("Failed to insert data: %s", err.Error())
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("Failed to commit transaction %s", err.Error())
	}

	return nil
}

func (r *FavoriteRepo) GetFavouritePost(ctx context.Context, userID string) ([]models.Post, error) {
	var posts []models.Post
	query := `
		SELECT 
			p.id, 
			p.caption, 
			p.user_id,
			p.media, 
			p.created_at, 
			p.updated_at 
		FROM posts p INNER JOIN favorites f 
		ON p.user_id = f.user_id 
		WHERE f.user_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return []models.Post{}, fmt.Errorf("Failed to get data: %s", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Caption, &post.UserID, &post.Media, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return []models.Post{}, fmt.Errorf("Failed to scan: %s", err.Error())
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *FavoriteRepo) Delete(ctx context.Context, favoriteID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("Failed to begin transaction: %s", err.Error())
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); err != nil {
				r.logger.Error("Database Error", "Failed to rollback", rbErr.Error())
			}
		}
	}()

	query := `DELETE FROM favorites WHERE id = $1`

	_, err = r.db.ExecContext(ctx, query, favoriteID)
	if err != nil {
		return fmt.Errorf("Failed to delete data: %s", err.Error())
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("Failed to commit transaction %s", err.Error())
	}

	return nil
}