package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cakra17/social/internal/models"
	"github.com/cakra17/social/internal/utils"
)

type LikesRepo struct {
	db *sql.DB
	logger *utils.Logger
}

func NewLikesRepo(db *sql.DB, lg *utils.Logger) LikesRepo {
	return LikesRepo{db: db, logger: lg}
}

func (r *LikesRepo) Like(ctx context.Context, likes models.Likes) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction error: %s", err.Error())
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				r.logger.Error("Database Error", "rollback failed", err.Error())
			}
		}
	}()

	query := `
		INSERT INTO likes (id, post_id, user_id) VALUES ($1, $2, $3)
	`

	_, err = r.db.ExecContext(ctx, query, likes.ID, likes.PostId, likes.UserId)
	if err != nil {
		return fmt.Errorf("Failed to add like: %s", err.Error())
	}

	return nil
}

func (r *LikesRepo) GetLikes(ctx context.Context, postId string) (models.Likes,error) {
	var likes models.Likes
	query := `
		SELECT u.id, u.username FROM users u INNER JOIN likes l ON u.id = l.user_id WHERE l.post_id = $1
	`
	row, err := r.db.QueryContext(ctx, query, postId)	
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Likes{}, fmt.Errorf("Failed post's likes not found")
		}
		return models.Likes{}, fmt.Errorf("Failed to get like data: %s", err.Error())
	}
	defer row.Close()	

	for row.Next() {
		var user models.User 
		err := row.Scan(&user.ID, &user.Username)
		if err != nil {
			return models.Likes{}, err
		}
		likes.Users = append(likes.Users, user)
	}

	likes.LikesCount = len(likes.Users)
	return likes, nil
}

func (r *LikesRepo) Unlike(ctx context.Context, postId string) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout) 
	defer cancel()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				r.logger.Error("Database Error", "rollback failed", err.Error())
			}
		}
	}()

	query := `DELETE FROM likes WHERE id = $1`
	_, err = r.db.ExecContext(ctx, query, postId)
	if err != nil {
		return err
	}

	if err != tx.Commit() {
		return fmt.Errorf("Failed to unlike: %s", err.Error())
	}

	return nil
}