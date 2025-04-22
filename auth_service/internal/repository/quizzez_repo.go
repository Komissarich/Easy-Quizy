package repository

import (
	"database/sql"
	"fmt"
)

type QuizzezRepository struct {
	db *sql.DB
}

func NewQuizzesRepo(DB *sql.DB) *QuizzezRepository {
	return &QuizzezRepository{
		db: DB,
	}
}

func (q *QuizzezRepository) AddFavoriteQuiz(userID uint64, quizID string) error {
	_, err := q.db.Exec(
		"INSERT INTO favorite_quizzes (user_id, quiz_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		userID,
		quizID,
	)

	if err != nil {
		return fmt.Errorf("can't add favorite quiz: %w", err)
	}

	return nil
}

func (q *QuizzezRepository) RemoveFavoriteQuiz(userID uint64, quizID string) error {
	_, err := q.db.Exec(
		"DELETE FROM favorite_quizzes WHERE user_id = $1 AND quiz_id = $2",
		userID,
		quizID,
	)

	if err != nil {
		return fmt.Errorf("failse to remove favorite quiz: %w", err)
	}

	return nil
}

func (q *QuizzezRepository) GetFavoriteQuizzes(userID uint64) ([]string, error) {
	rows, err := q.db.Query(
		"SELECT quiz_id FROM favorite_quizzes WHERE user_id = $1 ORDER BY created_at DESC",
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get favorite quizzes: %w", err)
	}
	defer rows.Close()

	var quizIDs []string
	for rows.Next() {
		var quizID string
		if err := rows.Scan(&quizID); err != nil {
			return nil, fmt.Errorf("failed to scan favorite quiz: %w", err)
		}
		quizIDs = append(quizIDs, quizID)
	}

	return quizIDs, nil
}
