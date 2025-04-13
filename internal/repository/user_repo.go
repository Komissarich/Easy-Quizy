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

func (u *UserRepository) SaveUser(ctx context.Context, email string, passHash []byte) error {
	return u.doSaveUser(ctx, email, passHash, entity.UserRole)
}

func (u *UserRepository) SaveAdmin(ctx context.Context, email string, passHash []byte) error {
	return u.doSaveUser(ctx, email, passHash, entity.AdminRole)
}

func (u *UserRepository) doSaveUser(ctx context.Context, email string, passHash []byte, role entity.Role) error {
	sql := `INSERT INTO users (email, pass_hash, role) VALUES ($1, $2, $3)`

	_, err := u.db.ExecContext(ctx, sql, email, passHash, role)
	if err != nil {

		return fmt.Errorf("can't to save user: %w", err)
	}

	return nil
}

func (u *UserRepository) User(ctx context.Context, email string) (*entity.User, error) {
	sql := `
	SELECT 
    u.id, u.email, u.pass_hash, u.role,
	FROM users u
	WHERE u.email = $1
	GROUP BY u.id, u.email, u.pass_hash, u.role;
	`

	rows, err := u.db.QueryContext(ctx, sql, email)
	if err != nil {
		return nil, fmt.Errorf("can't to find user by email: %w", err)
	}
	defer rows.Close()

	var user entity.User
	if rows.Next() {

		err := rows.Scan(&user.ID, &user.Email, &user.PassHash, &user.Role)
		if err != nil {
			return nil, fmt.Errorf("can't to scan user: %w", err)
		}

		return &user, nil
	}

	return nil, entity.ErrUserNotFound
}

func (u *UserRepository) Admins(ctx context.Context) ([]entity.Admin, error) {
	sql := `SELECT id, email FROM users WHERE role = $1`

	rows, err := u.db.QueryContext(ctx, sql, entity.AdminRole)
	if err != nil {
		return nil, fmt.Errorf("can't find admins: %w", err)
	}
	defer rows.Close()

	var admins = make([]entity.Admin, 0)

	for rows.Next() {
		var admin entity.Admin

		err = rows.Scan(&admin.UserID, &admin.Email)
		if err != nil {
			return nil, fmt.Errorf("can't scan admin: %w", err)
		}

		admins = append(admins, admin)
	}

	return admins, nil
}

func (u *UserRepository) DeleteAdmin(ctx context.Context, userID uint64) error {
	sql := `DELETE FROM users WHERE id = $1`

	_, err := u.db.ExecContext(ctx, sql, userID)
	if err != nil {
		return fmt.Errorf("can't delete user: %w", err)
	}

	return nil
}

func (u *UserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	sql := `SELECT id, email, role, pass_hash FROM users WHERE id = $1`

	var user entity.User

	err := u.db.QueryRowContext(ctx, sql, id).Scan(&user.ID, &user.Email, &user.Role, &user.PassHash)
	if err != nil {
		return nil, fmt.Errorf("can't find user by id: %w", err)
	}

	return &user, nil
}
