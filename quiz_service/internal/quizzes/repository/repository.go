package repository

import (
	"context"
	"fmt"
	"quizzes/internal/config"
	v1 "quizzes/pkg/api/v1"
	"quizzes/pkg/logger"
	"quizzes/pkg/postgres"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(ctx context.Context, config *config.Config) *Repository {
	pg, err := postgres.New(ctx, config.Postgres)
	//	fmt.Println(err.Error())
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, fmt.Sprint("failed to create repository", zap.Error(err)))
	} else {
		postgres.InitTables(ctx, pg)
		logger.GetLoggerFromCtx(ctx).Info(ctx, "connected to postgres")
		logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprint("pinging postgres: ", pg.Ping(ctx)))
	}
	return &Repository{
		pool: pg,
	}
}
func (r *Repository) CloseConn() {
	r.pool.Close()
}
func (r *Repository) CreateQuiz(
	ctx context.Context,
	name string,
	author string,
	image_id *string,
	description *string,
	questions []*v1.CreateQuestion,
) (string, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	quizID := uuid.New().String()

	_, err = tx.Exec(ctx,
		"INSERT INTO quizzes (Quiz_ID, Name, Author, Image_ID, Description) VALUES ($1, $2, $3, $4, $5)",
		quizID, name, author, &image_id, &description)
	if err != nil {
		return "", fmt.Errorf("failed to insert quiz: %w", err)
	}
	for _, q := range questions {
		questionID := uuid.New().String()
		_, err = tx.Exec(ctx,
			"INSERT INTO questions (Question_ID, Quiz_ID, Question_text, Image_ID) VALUES ($1, $2, $3, $4)",
			questionID, quizID, q.QuestionText, &q.ImageId)
		if err != nil {
			return "", fmt.Errorf("failed to insert question: %w", err)
		}

		for _, a := range q.Answer {
			answerID := uuid.New().String()
			_, err = tx.Exec(ctx,
				"INSERT INTO answers (Answer_ID, Question_ID, Answer_text, Is_correct) VALUES ($1, $2, $3, $4)",
				answerID, questionID, a.AnswerText, a.IsCorrect)
			if err != nil {
				return "", fmt.Errorf("failed to insert answer: %w", err)
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return quizID, nil
}

func (r *Repository) GetQuiz(
	ctx context.Context,
	quizID string,
) (*v1.GetQuizResponse, error) {
	var name, author, image_id, description string
	err := r.pool.QueryRow(ctx,
		"SELECT Name, Author, Image_ID, Description FROM quizzes WHERE Quiz_ID = $1",
		quizID).Scan(&name, &author, &image_id, &description)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("quiz not found")
		}
		return nil, fmt.Errorf("failed to get quiz: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		"SELECT Question_ID, Question_text, Image_ID FROM questions WHERE Quiz_ID = $1",
		quizID)
	if err != nil {
		return nil, fmt.Errorf("failed to get questions: %w", err)
	}
	defer rows.Close()

	var questions []*v1.CreateQuestion
	for rows.Next() {
		var questionID, questionText, imageID string
		err = rows.Scan(&questionID, &questionText, &imageID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan question: %w", err)
		}

		answerRows, err := r.pool.Query(ctx,
			"SELECT Answer_text, Is_correct FROM answers WHERE Question_ID = $1",
			questionID)
		if err != nil {
			return nil, fmt.Errorf("failed to get answers: %w", err)
		}
		defer answerRows.Close()

		var answers []*v1.CreateAnswer
		for answerRows.Next() {
			var answerText string
			var isCorrect bool
			err = answerRows.Scan(&answerText, &isCorrect)
			if err != nil {
				return nil, fmt.Errorf("failed to scan answer: %w", err)
			}
			answers = append(answers, &v1.CreateAnswer{
				AnswerText: answerText,
				IsCorrect:  isCorrect,
			})
		}
		if answerRows.Err() != nil {
			return nil, fmt.Errorf("error iterating answers: %w", answerRows.Err())
		}

		questions = append(questions, &v1.CreateQuestion{
			QuestionText: questionText,
			ImageId:      &imageID,
			Answer:       answers,
		})
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating questions: %w", rows.Err())
	}
	fmt.Println(name, author, image_id, description, questions)
	return &v1.GetQuizResponse{
		Name:        name,
		Author:      author,
		ImageId:     &image_id,
		Description: &description,
		Question:    questions,
	}, nil
}
