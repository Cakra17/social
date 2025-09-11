package store

import (
	"context"
	"database/sql"

	"github.com/cakra17/social/internal/models"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepo {
	return UserRepo{ db: db }
}

func (r *UserRepo) CreateUser(ctx context.Context, user *models.User) error {
	query := `
	INSERT INTO users (
		id, username, email, password
	) VALUES (
		$1, $2, $3, $4
	) RETURNING created_at
	`
	err := r.db.QueryRowContext(
		ctx, query,
		user.ID,
		user.Username,
		user.Email,
		user.Password,
	).Scan(&user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}