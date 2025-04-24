package repository

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"quizzes/internal/config"
	v1 "quizzes/pkg/api/v1"
	"quizzes/pkg/logger"
	"quizzes/pkg/postgres"
	"strings"

	pb "quizzes/pkg/authapi/v1"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Repository struct {
	pool *pgxpool.Pool
}
type IDGenerator struct {
	pool    *pgxpool.Pool
	charset string
}

// NewIDGenerator создает новый генератор
func NewIDGenerator(pool *pgxpool.Pool) *IDGenerator {
	return &IDGenerator{
		pool:    pool,
		charset: "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
	}
}

// GenerateID генерирует уникальный 5-символьный ключ
func (g *IDGenerator) GenerateID(ctx context.Context) (string, error) {
	// Максимум попыток, чтобы избежать бесконечного цикла
	const maxAttempts = 10000
	for attempt := 0; attempt < maxAttempts; attempt++ {
		key, err := g.generateRandomKey(5)
		if err != nil {
			return "", err
		}

		// Проверяем уникальность в базе данных
		unique, err := g.isUniqueInDB(ctx, key)
		if err != nil {
			return "", err
		}
		if unique {
			return key, nil
		}
	}

	return "", fmt.Errorf("unable to create unique code after %d atempts", maxAttempts)
}

// generateRandomKey генерирует случайный ключ заданной длины
func (g *IDGenerator) generateRandomKey(length int) (string, error) {
	result := make([]byte, length)
	charsetLength := len(g.charset)

	for i := 0; i < length; i++ {
		// Генерируем случайный индекс для символа из charset
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(charsetLength)))
		if err != nil {
			return "", err
		}
		result[i] = g.charset[idx.Int64()]
	}

	return string(result), nil
}

// isUniqueInDB проверяет уникальность ключа в таблице quizzes
func (g *IDGenerator) isUniqueInDB(ctx context.Context, key string) (bool, error) {
	var count int
	err := g.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM quizzes WHERE Quiz_ID = $1
	`, key).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
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
) (string, string, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	generator := NewIDGenerator(r.pool)
	quizID, err := generator.GenerateID(ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to create an id: %w", err)
	}
	uuID := uuid.New().String()
	fmt.Println(questions)
	_, err = tx.Exec(ctx,
		"INSERT INTO quizzes (Quiz_ID, Name, Author, Image_ID, Description, UUID) VALUES ($1, $2, $3, $4, $5, $6)",
		quizID, name, author, &image_id, &description, uuID)
	if err != nil {
		return "", "", fmt.Errorf("failed to insert quiz: %w", err)
	}
	for _, q := range questions {
		questionID := uuid.New().String()
		_, err = tx.Exec(ctx,
			"INSERT INTO questions (Question_ID, Quiz_ID, Question_text, Image_ID) VALUES ($1, $2, $3, $4)",
			questionID, quizID, q.QuestionText, &q.ImageId)
		if err != nil {
			return "", "", fmt.Errorf("failed to insert question: %w", err)
		}

		for _, a := range q.Answer {
			answerID := uuid.New().String()
			_, err = tx.Exec(ctx,
				"INSERT INTO answers (Answer_ID, Question_ID, Answer_text, Is_correct) VALUES ($1, $2, $3, $4)",
				answerID, questionID, a.AnswerText, a.IsCorrect)
			if err != nil {
				return "", "", fmt.Errorf("failed to insert answer: %w", err)
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return quizID, uuID, nil
}

func (r *Repository) GetQuiz(
	ctx context.Context,
	quizID string,
) (*v1.GetQuizResponse, error) {
	var name, author, image_id, description string
	fmt.Println("GET SOME QUIZ", quizID)
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
	return &v1.GetQuizResponse{
		Name:        name,
		Author:      author,
		ImageId:     &image_id,
		Description: &description,
		Question:    questions,
	}, nil
}

func (r *Repository) GetQuizByAuthor(
	ctx context.Context,
	author string,
) (*v1.GetQuizByAuthorResponse, error) {
	fmt.Println("HERE WE GO", author)
	rows, err := r.pool.Query(ctx,
		"SELECT Quiz_ID, Name, Image_ID, Description FROM quizzes WHERE Author = $1", author)
	if err != nil {
		return nil, fmt.Errorf("failed to get questions: %w", err)
	}
	defer rows.Close()
	var quizzesbyauthor []*v1.GetQuizResponse
	for rows.Next() {
		var quiz_ids, name, image_id, description string
		err = rows.Scan(&quiz_ids, &name, &image_id, &description)
		if err != nil {
			return nil, fmt.Errorf("failed to scan quiz: %w", err)
		}
		fmt.Println(quiz_ids)
		quiz, err := r.GetQuiz(ctx, quiz_ids)
		if err != nil {
			return nil, fmt.Errorf("failed to get quiz: %w", err)
		}
		quizzesbyauthor = append(quizzesbyauthor, quiz)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating questions: %w", rows.Err())
	}
	author_quizzes := &v1.GetQuizzes{Quizzes: quizzesbyauthor}
	conn, err := grpc.NewClient("auth_service:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth: %w", err)
	}
	defer conn.Close()
	client := pb.NewAuthServiceClient(conn)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
	}

	token := strings.TrimPrefix(authHeaders[0], "Bearer ")
	if token == "" {
		return nil, status.Error(codes.Unauthenticated, "invalid auth token format")
	}
	meta := metadata.Pairs("authorization", "Bearer "+token)
	ctx = metadata.NewOutgoingContext(context.Background(), meta)
	a := &pb.GetFavoriteQuizzesRequest{}
	var quizbyfav []*v1.GetQuizResponse
	response, err := client.GetFavoriteQuizzes(ctx, a)
	if err != nil {
		return nil, fmt.Errorf("failed to get quiz: %w", err)
	}
	for _, i := range response.QuizIds {
		quiz, err := r.GetQuiz(ctx, i)
		if err != nil {
			return nil, fmt.Errorf("failed to get quiz: %w", err)
		}
		quizbyfav = append(quizbyfav, quiz)
	}
	favouritequizzes := &v1.GetQuizzes{Quizzes: quizbyfav}
	var res []*v1.GetQuizzes
	res = append(res, author_quizzes)
	res = append(res, favouritequizzes)
	return &v1.GetQuizByAuthorResponse{AuthorQuizzes: res}, nil
}
