package repository

import (
	"context"
	"fmt"
	"quiz_app/internal/config"
	api "quiz_app/pkg/api/v1"
	"quiz_app/pkg/logger"
	"quiz_app/pkg/postgres"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Repository struct {
	pg *pgxpool.Pool
}

func NewRepository(ctx context.Context, config *config.Config) (*Repository, error) {
	pg, err := postgres.New(ctx, config.Postgres)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, fmt.Sprint("failed to create repository", zap.Error(err)))
		return nil, err
	} else {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "connected to postgres")
		logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprint("pinging postgres: ", pg.Ping(ctx)))
	}

	return &Repository{
		pg: pg,
	}, nil
}

func (r *Repository) CloseConn() {
	r.pg.Close()
}

func (r *Repository) UpdateStats(
	ctx context.Context,
	quiz_id string,
	players_score map[string]float32,
	quiz_rate float32,
) error {
	quiz_upd_query := `
	UPDATE stats.quizzes q
	SET
		q.avg_rate = (q.avg_rate * q.num_sessions + $1) / (q.num_sessions + 1),
		q.num_sessions = q.num_sessions + 1,
		q.updated_at = CURRENT_TIMESTAMP,
	WHERE
		q.quiz_id = $2;
	`
	author_upd_query := `
	UPDATE stats.authors a
	SET
		a.num_quizzes = (
			SELECT COUNT(DISTINCT q.quiz_id)
			FROM stats.quizzes q
			WHERE q.quiz_id = $1
		),
		a.avg_quiz_rate = (
			SELECT AVG(q.avg_rate)
			FROM stats.quizzes q
			WHERE q.quiz_id = $1
		),
		a.best_quiz_rate = (
			SELECT MAX(q.avg_rate)
			FROM stats.quizzes q
			WHERE q.quiz_id = $1
		),
		a.updated_at = CURRENT_TIMESTAMP
	WHERE a.user_id = (SELECT q.author_id FROM stats.quizzes q WHERE q.quiz_id = $1);
	`
	player_upd_query := `
	UPDATE stats.players p
	SET
		p.total_score = p.total_score + $1,
		p.best_score = MAX(p.best_score, $1),
		p.avg_score = (p.avg_score * p.num_sessions + $1) / (p.num_sessions + 1),
		p.num_sessions = p.num_sessions + 1
		p.updated_at = CURRENT_TIMESTAMP
	WHERE p.user_id = $2;
	`
	_, err := r.pg.Exec(ctx, quiz_upd_query, quiz_rate, quiz_id)
	if err != nil {
		return fmt.Errorf("unable to update quiz statistics: %w", err)
	}
	_, err = r.pg.Exec(ctx, author_upd_query, quiz_id)
	if err != nil {
		return fmt.Errorf("unable to update author statistics: %w", err)
	}
	for user_id, score := range players_score {
		_, err = r.pg.Exec(ctx, player_upd_query, score, user_id)
		if err != nil {
			return fmt.Errorf("unable to update player statistics: %w", err)
		}
	}
	return nil
}

func (r *Repository) GetQuizStat(ctx context.Context, quiz_id string) (*api.QuizStat, error) {
	quiz_stat_query := `
	SELECT 
		q.author_id, 
		q.num_sessions, 
		q.avg_rate
	FROM stats.quizzes q 
	WHERE q.quiz_id = $1;
	`
	var (
		author_id    string
		num_sessions int32
		avg_rate     float32
	)
	err := r.pg.QueryRow(ctx, quiz_stat_query, quiz_id).Scan(&author_id, &num_sessions, &avg_rate)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("quiz not found")
		}
		return nil, fmt.Errorf("unable to get quiz statistics: %w", err)
	}
	return &api.QuizStat{
		QuizId:      quiz_id,
		AuthorId:    author_id,
		NumSessions: num_sessions,
		AvgRate:     avg_rate,
	}, nil
}

func (r *Repository) ListQuizzes(ctx context.Context, option api.ListQuizzesOption) ([]*api.QuizStat, error) {
	var order string
	switch option {
	case api.ListQuizzesOption_AVG_RATE:
		order = "avg_rate"
	case api.ListQuizzesOption_NUM_SESSIONS:
		order = "num_sessions"
	default:
		return nil, fmt.Errorf("no such option: %d", option)
	}
	list_query := fmt.Sprintf(`
	SELECT 
		q.quiz_id,
		q.author_id,
		q.num_sessions,
		q.avg_rate
	FROM statc.quizzes q
	ORDER BY q.%s
	`, order)
	rows, err := r.pg.Query(ctx, list_query)
	if err != nil {
		return nil, fmt.Errorf("unable to list quizzes: %w", err)
	}
	defer rows.Close()
	var results []*api.QuizStat
	for rows.Next() {
		var (
			quiz_id      string
			author_id    string
			num_sessions int32
			avg_rate     float32
		)
		err = rows.Scan(&quiz_id, &author_id, &num_sessions, &avg_rate)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		results = append(results,
			&api.QuizStat{
				QuizId:      quiz_id,
				AuthorId:    author_id,
				NumSessions: num_sessions,
				AvgRate:     avg_rate,
			})
	}
	return results, nil
}

func (r *Repository) GetPlayerStat(ctx context.Context, user_id string) (*api.PlayerStat, error) {
	player_stat_query := `
	SELECT 
		p.total_score,
		p.best_score,
		p.avg_score,
		p.num_sessions
	FROM stats.players p 
	WHERE p.user_id = $1;
	`
	var (
		total_score  float32
		best_score   float32
		avg_score    float32
		num_sessions int32
	)
	err := r.pg.QueryRow(ctx, player_stat_query, user_id).Scan(&total_score, &best_score, &avg_score, &num_sessions)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("player not found")
		}
		return nil, fmt.Errorf("unable to get player statistics: %w", err)
	}
	return &api.PlayerStat{
		UserId:      user_id,
		TotalScore:  total_score,
		BestScore:   best_score,
		AvgScore:    avg_score,
		NumSessions: num_sessions,
	}, nil
}

func (r *Repository) ListPlayers(ctx context.Context, option api.ListPlayersOption) ([]*api.PlayerStat, error) {
	var order string
	switch option {
	case api.ListPlayersOption_TOTAL_SCORE:
		order = "total_score"
	case api.ListPlayersOption_BEST_SCORE:
		order = "best_score"
	case api.ListPlayersOption_AVG_SCORE:
		order = "avg_score"
	default:
		return nil, fmt.Errorf("no such option: %d", option)
	}
	list_query := fmt.Sprintf(`
	SELECT 
		p.user_id
		p.total_score,
		p.best_score,
		p.avg_score,
		p.num_sessions
	FROM stats.players p
	ORDER BY p.%s
	`, order)
	rows, err := r.pg.Query(ctx, list_query)
	if err != nil {
		return nil, fmt.Errorf("unable to list players: %w", err)
	}
	defer rows.Close()
	var results []*api.PlayerStat
	for rows.Next() {
		var (
			user_id      string
			total_score  float32
			best_score   float32
			avg_score    float32
			num_sessions int32
		)
		err = rows.Scan(&user_id, &total_score, &best_score, &avg_score, &num_sessions)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		results = append(results,
			&api.PlayerStat{
				UserId:      user_id,
				TotalScore:  total_score,
				BestScore:   best_score,
				AvgScore:    avg_score,
				NumSessions: num_sessions,
			})
	}
	return results, nil
}

func (r *Repository) GetAuthorStat(ctx context.Context, user_id string) (*api.AuthorStat, error) {
	author_stat_query := `
	SELECT 
		a.num_quizzes,
    	a.avg_quiz_rate,
		a.best_quiz_rate,
	FROM stats.authors a 
	WHERE a.user_id = $1;
	`
	var (
		num_quizzes    int32
		avg_quiz_rate  float32
		best_quiz_rate float32
	)
	err := r.pg.QueryRow(ctx, author_stat_query, user_id).Scan(&num_quizzes, &avg_quiz_rate, &best_quiz_rate)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("author not found")
		}
		return nil, fmt.Errorf("unable to get author statistics: %w", err)
	}
	return &api.AuthorStat{
		UserId:       user_id,
		NumQuizzes:   num_quizzes,
		AvgQuizRate:  avg_quiz_rate,
		BestQuizRate: best_quiz_rate,
	}, nil
}

func (r *Repository) ListAuthors(ctx context.Context, option api.ListAuthorsOption) ([]*api.AuthorStat, error) {
	var order string
	switch option {
	case api.ListAuthorsOption_NUM_QUIZZES:
		order = "num_quizzes"
	case api.ListAuthorsOption_AVG_QUIZ_RATE:
		order = "avg_quiz_rate"
	case api.ListAuthorsOption_BEST_QUIZ_RATE:
		order = "best_quiz_rate"
	default:
		return nil, fmt.Errorf("no such option: %d", option)
	}
	list_query := fmt.Sprintf(`
	SELECT 
		a.user_id,
		a.num_quizzes,
    	a.avg_quiz_rate,
		a.best_quiz_rate,
	FROM stats.authors a
	ORDER BY a.%s
	`, order)
	rows, err := r.pg.Query(ctx, list_query)
	if err != nil {
		return nil, fmt.Errorf("unable to list authors: %w", err)
	}
	defer rows.Close()
	var results []*api.AuthorStat
	for rows.Next() {
		var (
			user_id        string
			num_quizzes    int32
			avg_quiz_rate  float32
			best_quiz_rate float32
		)
		err = rows.Scan(&user_id, &num_quizzes, &avg_quiz_rate, &best_quiz_rate)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		results = append(results,
			&api.AuthorStat{
				UserId:       user_id,
				NumQuizzes:   num_quizzes,
				AvgQuizRate:  avg_quiz_rate,
				BestQuizRate: best_quiz_rate,
			})
	}
	return results, nil
}
