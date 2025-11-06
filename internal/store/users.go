package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/cakra17/social/internal/models"
)

type UserRepo struct {
	db *sql.DB
}

const (
	defaultTimeout = 5 * time.Second
)

func NewUserRepo(db *sql.DB) UserRepo {
	return UserRepo{ db: db }
}

func (r *UserRepo) CreateUser(ctx context.Context, user *models.User) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

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

func (r *UserRepo) GetUserById(ctx context.Context, id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	query := `
		SELECT id, username, email, created_at
		FROM users WHERE id = $1
	`

	user := &models.User{}
	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
	)
	
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, errors.New("User Not Found!")
		default:
			return nil, err
		}
	}
	return user, nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

query := `
		SELECT id, username, email, password
		FROM users WHERE email = $1
	`

	user := &models.User{}
	row := r.db.QueryRowContext(ctx, query, email)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
	)
	
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, errors.New("User Not Found!")
		default:
			return nil, err
		}
	}
	return user, nil
}

func (r *UserRepo) UpdateUser(ctx context.Context, user *models.UpdateUserPayload, ID string) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		UPDATE users SET username = $1, email = $2 WHERE id = $3
	`
	_, err = tx.ExecContext(ctx, query, user.Username, user.Email, ID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `DELETE FROM users WHERE id = $1`
	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}