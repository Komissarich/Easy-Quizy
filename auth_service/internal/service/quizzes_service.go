package service

import (
	"eazy-quizy-auth/internal/repository"
	"eazy-quizy-auth/pkg/logger"
)

type QuizzesService interface {
	AddFavoriteQuiz(userID uint64, quizID string) error
	GetFavoriteQuizzes(userID uint64) ([]string, error)
	RemoveFavoriteQuiz(userID uint64, quizID string) error
}

type quizzesService struct {
	repo repository.QuizzezRepository
	l    *logger.Logger
}

func NewQuizzesService(repo repository.QuizzezRepository, l *logger.Logger) *quizzesService {
	return &quizzesService{
		repo: repo,
		l:    l,
	}
}

func (q *quizzesService) AddFavoriteQuiz(userID uint64, quizID string) error {
	return q.repo.AddFavoriteQuiz(userID, quizID)
}

func (q *quizzesService) GetFavoriteQuizzes(userID uint64) ([]string, error) {
	return q.repo.GetFavoriteQuizzes(userID)
}

func (q *quizzesService) RemoveFavoriteQuiz(userID uint64, quizID string) error {
	return q.repo.RemoveFavoriteQuiz(userID, quizID)
}
