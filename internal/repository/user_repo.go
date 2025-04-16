package repository

import (
	"context"
	"database/sql"
	"eazy-quizy-auth/internal/entity"
	"fmt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) SaveUser(ctx context.Context, email string, passHash []byte) (uint64, error) {
	return u.doSaveUser(ctx, email, passHash)
}

func (u *UserRepository) doSaveUser(ctx context.Context, email string, passHash []byte) (uint64, error) {
	sql := `INSERT INTO users (email, pass_hash) VALUES ($1, $2) RETURNING id`

	var id uint64
	err := u.db.QueryRowContext(ctx, sql, email, passHash).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("can't to save user: %w", err)
	}

	return id, nil
}

func (u *UserRepository) User(ctx context.Context, email string) (*entity.User, error) {
	sql := `
	SELECT 
    u.id, u.email, u.pass_hash,
	FROM users u
	WHERE u.email = $1
	GROUP BY u.id, u.email, u.pass_hash;
	`

	rows, err := u.db.QueryContext(ctx, sql, email)
	if err != nil {
		return nil, fmt.Errorf("can't to find user by email: %w", err)
	}
	defer rows.Close()

	var user entity.User
	if rows.Next() {

		err := rows.Scan(&user.ID, &user.Email, &user.Password)
		if err != nil {
			return nil, fmt.Errorf("can't to scan user: %w", err)
		}

		return &user, nil
	}

	return nil, entity.ErrUserNotFound
}

func (u *UserRepository) FindByID(ctx context.Context, id uint64) (*entity.User, error) {
	sql := `SELECT id, email, pass_hash FROM users WHERE id = $1`

	var user entity.User

	err := u.db.QueryRowContext(ctx, sql, id).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("can't find user by id: %w", err)
	}

	return &user, nil
}

func (u *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	sql := `SELECT id, email, pass_hash FROM users WHERE email = $1`

	var user entity.User

	err := u.db.QueryRowContext(ctx, sql, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("can't find user by email: %w", err)
	}

	return &user, nil
}

func (u *UserRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	sql := `UPDATE users SET email = $1, pass_hash = $2, WHERE id = $3`

	_, err := u.db.ExecContext(ctx, sql, user.Email, user.Password, user.ID)
	if err != nil {
		return fmt.Errorf("can't update user: %w", err)
	}

	return nil
}
