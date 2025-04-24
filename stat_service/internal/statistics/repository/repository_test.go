package repository

import (
	"context"
	"errors"
	"testing"

	api "quiz_app/pkg/api/v1"
	"quiz_app/pkg/logger"

	"quiz_app/internal/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateStats(t *testing.T) {
	ctx, _ := logger.New(context.Background())
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := &Repository{pg: mock}

	t.Run("successful update", func(t *testing.T) {
		quizID := "quiz1"
		authorID := "author1"
		players := map[string]float32{"player1": 0.8, "player2": 0.9}
		rate := float32(4.5)

		// Ожидаем вызовы для каждого запроса
		mock.ExpectExec("UPDATE stats.quizzes").
			WithArgs(rate, quizID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		mock.ExpectExec("UPDATE stats.authors").
			WithArgs(quizID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		// По одному для каждого игрока
		mock.ExpectExec("UPDATE stats.players").
			WithArgs(players["player1"], "player1").
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		mock.ExpectExec("UPDATE stats.players").
			WithArgs(players["player2"], "player2").
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := repo.UpdateStats(ctx, quizID, authorID, players, rate)
		assert.NoError(t, err)
	})

	t.Run("quiz update error", func(t *testing.T) {
		mock.ExpectExec("UPDATE stats.quizzes").
			WillReturnError(errors.New("update failed"))

		err := repo.UpdateStats(ctx, "quiz1", "author1", map[string]float32{}, 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unable to update quiz statistics")
	})

	// Аналогичные тесты для других ошибок (author update, player update)
}

func TestGetQuizStat(t *testing.T) {
	ctx, _ := logger.New(context.Background())
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := &Repository{pg: mock}

	t.Run("successful get", func(t *testing.T) {
		quizID := "quiz1"
		expected := &api.QuizStat{
			QuizId:      quizID,
			AuthorId:    "author1",
			NumSessions: 10,
			AvgRate:     4.5,
		}

		mock.ExpectQuery("SELECT").
			WithArgs(quizID).
			WillReturnRows(pgxmock.NewRows([]string{"author_id", "num_sessions", "avg_rate"}).
				AddRow(expected.AuthorId, expected.NumSessions, expected.AvgRate))

		result, err := repo.GetQuizStat(ctx, quizID)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT").
			WillReturnError(pgx.ErrNoRows)

		_, err := repo.GetQuizStat(ctx, "nonexistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unable to get quiz statistics")
	})
}

func TestListQuizzes(t *testing.T) {
	ctx, _ := logger.New(context.Background())
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := &Repository{pg: mock}

	t.Run("by avg rate", func(t *testing.T) {
		expected := []*api.QuizStat{
			{
				QuizId:      "quiz1",
				AuthorId:    "author1",
				NumSessions: 5,
				AvgRate:     4.5,
			},
			{
				QuizId:      "quiz2",
				AuthorId:    "author2",
				NumSessions: 10,
				AvgRate:     4.0,
			},
		}

		mock.ExpectQuery("SELECT").
			WillReturnRows(pgxmock.NewRows([]string{"quiz_id", "author_id", "num_sessions", "avg_rate"}).
				AddRow(expected[0].QuizId, expected[0].AuthorId, expected[0].NumSessions, expected[0].AvgRate).
				AddRow(expected[1].QuizId, expected[1].AuthorId, expected[1].NumSessions, expected[1].AvgRate))

		result, err := repo.ListQuizzes(ctx, api.ListQuizzesOption_AVG_RATE)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("invalid option", func(t *testing.T) {
		_, err := repo.ListQuizzes(ctx, 999)
		assert.Error(t, err)
	})
}

// Аналогичные тесты для:
// - GetPlayerStat
// - ListPlayers
// - GetAuthorStat
// - ListAuthors

func TestCloseConn(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)

	repo := &Repository{pg: mock}
	mock.ExpectClose()

	repo.CloseConn()
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Мок для postgres.New
var postgresNew = func(ctx context.Context, config config.Config) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, config.Postgres.Host)
}
