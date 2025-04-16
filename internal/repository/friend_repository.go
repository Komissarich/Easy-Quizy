package repository

import (
	"context"
	"database/sql"
	"eazy-quizy-auth/internal/entity"
	"fmt"
)

type FriendRepository struct {
	db *sql.DB
}

func NewFriendRepository(db *sql.DB) *FriendRepository {
	return &FriendRepository{db: db}
}

func (r *FriendRepository) AddFriend(ctx context.Context, userID uint64, friendID string) error {
	query := `INSERT INTO friends (user_id, friend_id) VALUES ($1, $2), ($2, $1)`
	_, err := r.db.ExecContext(ctx, query, userID, friendID)
	return err
}

func (r *FriendRepository) RemoveFriend(ctx context.Context, userID uint64, friendID string) error {
	query := `DELETE FROM friends WHERE (user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1)`
	_, err := r.db.ExecContext(ctx, query, userID, friendID)
	return err
}

func (r *FriendRepository) GetFriendIDs(ctx context.Context, userID uint64) ([]string, error) {
	query := `SELECT friend_id FROM friends WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friendIDs []string
	for rows.Next() {
		var friendID string
		if err := rows.Scan(&friendID); err != nil {
			return nil, err
		}
		friendIDs = append(friendIDs, friendID)
	}

	return friendIDs, nil
}

func (r *FriendRepository) CheckFriendship(ctx context.Context, userID uint64, friendID string) (bool, error) {
	query := `SELECT COUNT(*) FROM friends WHERE user_id = $1 AND friend_id = $2`

	var count int

	err := r.db.QueryRowContext(ctx, query, userID, friendID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *FriendRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	query := `SELECT id, email, username FROM users WHERE id = $1`

	var user entity.User

	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Email, &user.Username)
	if err != nil {
		return nil, fmt.Errorf("can't find user by id: %w", err)
	}

	return &user, nil
}
