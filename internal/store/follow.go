package store

import (
	"context"
	"database/sql"

	"github.com/cakra17/social/internal/models"
)

type FollowRepo struct {
	db *sql.DB
}

func NewFollowRepo(db *sql.DB) FollowRepo {
	return FollowRepo{ db: db }
}

func (r *FollowRepo) Follow(ctx context.Context, f models.Follow) error {
	query := `
		INSERT INTO followers (
			id, followers_id, followee_id
		) VALUES (
			$1, $2, $3
		)
	`
	_, err := r.db.ExecContext(ctx, query, f.ID, f.FollowerID, f.FolloweeID)
	if err != nil {
		return err
	}
	return nil
}

func (r *FollowRepo) GetFollowers(ctx context.Context, userId string) ([]models.Follower, error) {
	var followers []models.Follower
	query := `
		SELECT f.id, u.id AS user_id, u.username 
		FROM followers f 
		INNER JOIN users u ON u.id = f.followers_id 
		WHERE f.followee_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var follower models.Follower
		if err := rows.Scan(&follower.ID, &follower.UserID, &follower.Username); err != nil {
			return nil, err
		}
		followers = append(followers, follower)
	}

	return followers, nil
}

func (r *FollowRepo) GetFollowing(ctx context.Context, userId string) ([]models.Follower, error) {
	var followers []models.Follower
	query := `
		SELECT f.id, u.id AS user_id, u.username 
		FROM followers f 
		INNER JOIN users u ON u.id = f.followee_id 
		WHERE f.followers_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var follower models.Follower
		if err := rows.Scan(&follower.ID, &follower.UserID, &follower.Username); err != nil {
			return nil, err
		}
		followers = append(followers, follower)
	}

	return followers, nil
}

func (r *FollowRepo) Unfollow(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `DELETE FROM followers WHERE id = $1`
	_, err = r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}