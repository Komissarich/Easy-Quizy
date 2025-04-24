package repository

import (
	"context"
	"database/sql"
	"eazy-quizy-auth/internal/entity"
	"fmt"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) SaveUser(ctx context.Context, email, username string, passHash []byte) (uint64, error) {
	return u.doSaveUser(ctx, email, username, passHash)
}

func (u *UserRepository) doSaveUser(ctx context.Context, email, username string, passHash []byte) (uint64, error) {
	sql := `INSERT INTO users (email, username, password, created_at) VALUES ($1, $2, $3, $4) RETURNING id`

	var id uint64
	err := u.db.QueryRowContext(
		ctx,
		sql,
		email,
		username,
		passHash,
		time.Now(),
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("can't to save user: %w", err)
	}

	return id, nil
}

func (u *UserRepository) User(ctx context.Context, email string) (*entity.User, error) {
	sql := `SELECT id, username, email, password FROM users WHERE email = $1 GROUP BY id, username, email, password;`

	rows, err := u.db.QueryContext(ctx, sql, email)
	if err != nil {
		return nil, fmt.Errorf("can't to find user by email: %w", err)
	}
	defer rows.Close()

	var user entity.User
	if rows.Next() {

		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
		if err != nil {
			return nil, fmt.Errorf("can't to scan user: %w", err)
		}

		return &user, nil
	}

	return nil, entity.ErrUserNotFound
}

func (r *UserRepository) FindByID(ctx context.Context, id uint64) (*entity.User, error) {
	query := `
        SELECT id, email, username, password, created_at
        FROM users 
        WHERE id = $1`

	var user entity.User

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, id string) (*entity.User, error) {
	query := `
        SELECT id, email, username, password, created_at
        FROM users 
        WHERE username = $1`

	var user entity.User

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &user, nil
}

func (u *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	sql := `SELECT id, username, email, password FROM users WHERE email = $1`

	var user entity.User

	err := u.db.QueryRowContext(ctx, sql, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("can't find user by email: %w", err)
	}

	return &user, nil
}

func (u *UserRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	sql := `
        UPDATE users 
        SET 
            email = COALESCE($1, email),
            username = COALESCE($2, username),
            password = COALESCE($3, password),
        WHERE id = $4
    `

	var email *string
	if user.Email != "" {
		email = &user.Email
	}

	var username *string
	if user.Username != "" {
		username = &user.Username
	}

	var password *string
	if user.Password != "" {
		password = &user.Password
	}

	_, err := u.db.ExecContext(ctx, sql,
		email,
		username,
		password,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user %d: %w", user.ID, err)
	}

	return nil
}
