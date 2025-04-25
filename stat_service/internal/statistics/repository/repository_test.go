package repository

import (
	"context"
	"errors"
	"testing"

	"quiz_app/internal/config"
	api "quiz_app/pkg/api/v1"
	"quiz_app/pkg/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Мок для postgres.New
var postgresNew = func(ctx context.Context, config config.Config) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, config.Postgres.Host)
}

func TestUpdateStats(t *testing.T) {
	ctx, _ := logger.New(context.Background())
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := &Repository{pg: mock}

	quizID := "quiz1"
	authorID := "author1"
	player_id := "player1"
	player_score := float32(3.4)
	rate := float32(4.5)

	t.Run("successful update", func(t *testing.T) {
		// Ожидаем вызовы для каждого запроса
		mock.ExpectExec("INSERT INTO stats.quizzes").
			WithArgs(rate, quizID, authorID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		mock.ExpectExec("WITH author_stats AS").
			WithArgs(authorID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		mock.ExpectExec("INSERT INTO stats.players").
			WithArgs(player_score, player_id).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := repo.UpdateStats(ctx, quizID, authorID, player_id, player_score, rate)
		assert.NoError(t, err)
	})

	t.Run("wrong format update error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO stats.quizzes").
			WithArgs(quizID, authorID).
			WillReturnError(errors.New("update failed"))

		err := repo.UpdateStats(ctx, "quiz1", "author1", "", 0, 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "wrong request format")
	})

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

	t.Run("by AvgRate", func(t *testing.T) {
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

	t.Run("by NumSessions", func(t *testing.T) {
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

		result, err := repo.ListQuizzes(ctx, api.ListQuizzesOption_NUM_SESSIONS)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("invalid option", func(t *testing.T) {
		_, err := repo.ListQuizzes(ctx, 999)
		assert.Error(t, err)
	})
}

func TestGetPlayerStat(t *testing.T) {
	ctx, _ := logger.New(context.Background())
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := &Repository{pg: mock}

	t.Run("successful get", func(t *testing.T) {
		userID := "player1"
		expected := &api.PlayerStat{
			UserId:      userID,
			TotalScore:  100,
			BestScore:   100,
			AvgScore:    100,
			NumSessions: 1,
		}

		mock.ExpectQuery("SELECT").
			WithArgs(userID).
			WillReturnRows(pgxmock.NewRows([]string{"total_score", "best_score", "avg_score", "num_sessions"}).
				AddRow(expected.TotalScore, expected.BestScore, expected.AvgScore, expected.NumSessions))

		result, err := repo.GetPlayerStat(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("not found", func(t *testing.T) {
		nonUser := "nonexistent"
		expected := &api.PlayerStat{
			UserId:      nonUser,
			TotalScore:  0,
			BestScore:   0,
			AvgScore:    0,
			NumSessions: 0,
		}
		mock.ExpectQuery("SELECT").
			WithArgs(nonUser).
			WillReturnRows(pgxmock.NewRows([]string{"total_score", "best_score", "avg_score", "num_sessions"}).
				AddRow(expected.TotalScore, expected.BestScore, expected.AvgScore, expected.NumSessions))

		result, err := repo.GetPlayerStat(ctx, nonUser)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestListPlayers(t *testing.T) {
	ctx, _ := logger.New(context.Background())
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := &Repository{pg: mock}

	t.Run("by TotalScore", func(t *testing.T) {
		expected := []*api.PlayerStat{
			{
				UserId:      "player1",
				TotalScore:  100,
				BestScore:   100,
				AvgScore:    100,
				NumSessions: 1,
			},
			{
				UserId:      "player2",
				TotalScore:  200,
				BestScore:   200,
				AvgScore:    200,
				NumSessions: 1,
			},
		}

		mock.ExpectQuery("SELECT").
			WillReturnRows(pgxmock.NewRows([]string{"user_id", "total_score", "best_score", "avg_score", "num_sessions"}).
				AddRow(expected[0].UserId, expected[0].TotalScore, expected[0].BestScore, expected[0].AvgScore, expected[0].NumSessions).
				AddRow(expected[1].UserId, expected[1].TotalScore, expected[1].BestScore, expected[1].AvgScore, expected[1].NumSessions))

		result, err := repo.ListPlayers(ctx, api.ListPlayersOption_TOTAL_SCORE)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("by BestScore", func(t *testing.T) {
		expected := []*api.PlayerStat{
			{
				UserId:      "player1",
				TotalScore:  100,
				BestScore:   100,
				AvgScore:    100,
				NumSessions: 1,
			},
			{
				UserId:      "player2",
				TotalScore:  200,
				BestScore:   200,
				AvgScore:    200,
				NumSessions: 2,
			},
		}

		mock.ExpectQuery("SELECT").
			WillReturnRows(pgxmock.NewRows([]string{"user_id", "total_score", "best_score", "avg_score", "num_sessions"}).
				AddRow(expected[0].UserId, expected[0].TotalScore, expected[0].BestScore, expected[0].AvgScore, expected[0].NumSessions).
				AddRow(expected[1].UserId, expected[1].TotalScore, expected[1].BestScore, expected[1].AvgScore, expected[1].NumSessions))

		result, err := repo.ListPlayers(ctx, api.ListPlayersOption_BEST_SCORE)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("by AvgScore", func(t *testing.T) {
		expected := []*api.PlayerStat{
			{
				UserId:      "player1",
				TotalScore:  100,
				BestScore:   100,
				AvgScore:    100,
				NumSessions: 1,
			},
			{
				UserId:      "player2",
				TotalScore:  200,
				BestScore:   200,
				AvgScore:    200,
				NumSessions: 2,
			},
		}

		mock.ExpectQuery("SELECT").
			WillReturnRows(pgxmock.NewRows([]string{"user_id", "total_score", "best_score", "avg_score", "num_sessions"}).
				AddRow(expected[0].UserId, expected[0].TotalScore, expected[0].BestScore, expected[0].AvgScore, expected[0].NumSessions).
				AddRow(expected[1].UserId, expected[1].TotalScore, expected[1].BestScore, expected[1].AvgScore, expected[1].NumSessions))

		result, err := repo.ListPlayers(ctx, api.ListPlayersOption_AVG_SCORE)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("invalid option", func(t *testing.T) {
		_, err := repo.ListPlayers(ctx, 999)
		assert.Error(t, err)
	})
}

func TestGetAuthorStat(t *testing.T) {
	ctx, _ := logger.New(context.Background())
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := &Repository{pg: mock}

	t.Run("successful get", func(t *testing.T) {
		userID := "author1"
		expected := &api.AuthorStat{
			UserId:       userID,
			NumQuizzes:   5,
			AvgQuizRate:  float32(4.5),
			BestQuizRate: float32(4.9),
		}

		mock.ExpectQuery("SELECT").
			WithArgs(userID).
			WillReturnRows(pgxmock.NewRows([]string{"num_quizzes", "avg_quiz_rate", "best_quiz_rate"}).
				AddRow(expected.NumQuizzes, expected.AvgQuizRate, expected.BestQuizRate))

		result, err := repo.GetAuthorStat(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("not found", func(t *testing.T) {
		nonUser := "nonexistent"
		expected := &api.AuthorStat{
			UserId:       nonUser,
			NumQuizzes:   0,
			AvgQuizRate:  0,
			BestQuizRate: 0,
		}
		mock.ExpectQuery("SELECT").
			WithArgs(nonUser).
			WillReturnRows(pgxmock.NewRows([]string{"num_quizzes", "avg_quiz_rate", "best_quiz_rate"}).
				AddRow(expected.NumQuizzes, expected.AvgQuizRate, expected.BestQuizRate))

		result, err := repo.GetAuthorStat(ctx, nonUser)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestListAuthors(t *testing.T) {
	ctx, _ := logger.New(context.Background())
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := &Repository{pg: mock}

	t.Run("by AvgQuizRate", func(t *testing.T) {
		expected := []*api.AuthorStat{
			{
				UserId:       "author1",
				NumQuizzes:   5,
				AvgQuizRate:  float32(3.5),
				BestQuizRate: float32(3.9),
			},
			{
				UserId:       "author2",
				NumQuizzes:   10,
				AvgQuizRate:  float32(4.5),
				BestQuizRate: float32(4.9),
			},
		}

		mock.ExpectQuery("SELECT").
			WillReturnRows(pgxmock.NewRows([]string{"user_id", "num_quizzes", "avg_quiz_rate", "best_quiz_rate"}).
				AddRow(expected[0].UserId, expected[0].NumQuizzes, expected[0].AvgQuizRate, expected[0].BestQuizRate).
				AddRow(expected[1].UserId, expected[1].NumQuizzes, expected[1].AvgQuizRate, expected[1].BestQuizRate))

		result, err := repo.ListAuthors(ctx, api.ListAuthorsOption_AVG_QUIZ_RATE)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("by BestQuizRate", func(t *testing.T) {
		expected := []*api.AuthorStat{
			{
				UserId:       "author1",
				NumQuizzes:   5,
				AvgQuizRate:  float32(3.5),
				BestQuizRate: float32(3.9),
			},
			{
				UserId:       "author2",
				NumQuizzes:   10,
				AvgQuizRate:  float32(4.5),
				BestQuizRate: float32(4.9),
			},
		}

		mock.ExpectQuery("SELECT").
			WillReturnRows(pgxmock.NewRows([]string{"user_id", "num_quizzes", "avg_quiz_rate", "best_quiz_rate"}).
				AddRow(expected[0].UserId, expected[0].NumQuizzes, expected[0].AvgQuizRate, expected[0].BestQuizRate).
				AddRow(expected[1].UserId, expected[1].NumQuizzes, expected[1].AvgQuizRate, expected[1].BestQuizRate))

		result, err := repo.ListAuthors(ctx, api.ListAuthorsOption_BEST_QUIZ_RATE)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("invalid option", func(t *testing.T) {
		_, err := repo.ListAuthors(ctx, 999)
		assert.Error(t, err)
	})
}

// Аналогичные тесты для:
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
