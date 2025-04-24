package repository

import (
	"context"
	"fmt"
	"quiz_app/internal/config"
	api "quiz_app/pkg/api/v1"
	"quiz_app/pkg/logger"
	"quiz_app/pkg/postgres"

	"github.com/jackc/pgx/v5"
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
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, "connected to postgres")
	logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprint("pinging postgres: ", pg.Ping(ctx)))

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
	author_id string,
	player_id string,
	player_score float32,
	quiz_rate float32,
) error {
	quiz_upd_query := `
	INSERT INTO stats.quizzes (quiz_id, author_id, avg_rate, num_sessions, updated_at)
	VALUES (
    	$2, 
		$3,
   		$1,  
    	1,   
    	CURRENT_TIMESTAMP
	)
	ON CONFLICT (quiz_id) DO UPDATE
	SET
    	avg_rate = (stats.quizzes.avg_rate * stats.quizzes.num_sessions + $1) / (stats.quizzes.num_sessions + 1),
    	num_sessions = stats.quizzes.num_sessions + 1,
    	updated_at = CURRENT_TIMESTAMP
	`
	author_upd_query := `
	WITH author_stats AS (
		SELECT
	    	COUNT(DISTINCT quiz_id) AS num_quizzes,
	    	AVG(avg_rate) AS avg_rate,
	    	MAX(avg_rate) AS best_rate
		FROM stats.quizzes
		WHERE author_id = $1
		)
	INSERT INTO stats.authors (user_id, num_quizzes, avg_quiz_rate, best_quiz_rate, updated_at)
	VALUES (
		$1,
		COALESCE((SELECT num_quizzes FROM author_stats), 0),
		COALESCE((SELECT avg_rate FROM author_stats), 0),
		COALESCE((SELECT best_rate FROM author_stats), 0),
		CURRENT_TIMESTAMP
	)
	ON CONFLICT (user_id) DO UPDATE
	SET
		num_quizzes = EXCLUDED.num_quizzes,
		avg_quiz_rate = EXCLUDED.avg_quiz_rate,
		best_quiz_rate = EXCLUDED.best_quiz_rate,
		updated_at = EXCLUDED.updated_at;
	`
	player_upd_query := `
	INSERT INTO stats.players (user_id, total_score, best_score, avg_score, num_sessions, updated_at)
	VALUES (
    	$2,
    	$1,  
    	$1,  
    	$1,  
    	1,   
    	CURRENT_TIMESTAMP
	)
	ON CONFLICT (user_id) DO UPDATE
	SET
    	total_score = stats.players.total_score + $1,
    	best_score = GREATEST(stats.players.best_score, $1),
    	avg_score = (stats.players.avg_score * stats.players.num_sessions + $1) / (stats.players.num_sessions + 1),
    	num_sessions = stats.players.num_sessions + 1,
    	updated_at = CURRENT_TIMESTAMP;
	`
	_, err := r.pg.Exec(ctx, quiz_upd_query, quiz_rate, quiz_id, author_id)
	if err != nil {
		return fmt.Errorf("unable to update quiz statistics: %w", err)
	}
	_, err = r.pg.Exec(ctx, author_upd_query, author_id)
	if err != nil {
		return fmt.Errorf("unable to update author statistics: %w", err)
	}

	_, err = r.pg.Exec(ctx, player_upd_query, player_score, player_id)
	if err != nil {
		return fmt.Errorf("unable to update player statistics: %w", err)
	}

	return nil
}

func (r *Repository) GetQuizStat(ctx context.Context, quiz_id string) (*api.QuizStat, error) {
	quiz_stat_query := `
	SELECT 
		stats.quizzes.author_id, 
		stats.quizzes.num_sessions, 
		stats.quizzes.avg_rate
	FROM stats.quizzes
	WHERE stats.quizzes.quiz_id = $1;
	`
	var (
		author_id    string
		num_sessions int32
		avg_rate     float32
	)
	err := r.pg.QueryRow(ctx, quiz_stat_query, quiz_id).Scan(&author_id, &num_sessions, &avg_rate)
	if err != nil {
		if err == pgx.ErrNoRows {
<<<<<<< HEAD
			if err.Error() == "no rows in result set" {
				return &api.QuizStat{
					UserId:      user_id,
					TotalScore:  0,
					BestScore:   0,
					AvgScore:    0,
					NumSessions: 0,
				}, nil
			}
=======
			return &api.QuizStat{
				QuizId:      quiz_id,
				AuthorId:    "",
				NumSessions: 0,
				AvgRate:     0,
			}, nil
>>>>>>> a6504fcd1740986d53766c00606b7aa4b9ebc120
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
		stats.quizzes.quiz_id,
		stats.quizzes.author_id,
		stats.quizzes.num_sessions,
		stats.quizzes.avg_rate
	FROM stats.quizzes
	ORDER BY stats.quizzes.%s;
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
		stats.players.total_score,
		stats.players.best_score,
		stats.players.avg_score,
		stats.players.num_sessions
	FROM stats.players
	WHERE stats.players.user_id = $1;
	`
	var (
		total_score  float32
		best_score   float32
		avg_score    float32
		num_sessions int32
	)

	err := r.pg.QueryRow(ctx, player_stat_query, user_id).Scan(&total_score, &best_score, &avg_score, &num_sessions)
	if err != nil {
<<<<<<< HEAD
		if err.Error() == "no rows in result set" {
=======
		if err == pgx.ErrNoRows {
>>>>>>> a6504fcd1740986d53766c00606b7aa4b9ebc120
			return &api.PlayerStat{
				UserId:      user_id,
				TotalScore:  0,
				BestScore:   0,
				AvgScore:    0,
				NumSessions: 0,
			}, nil
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
		stats.players.user_id,
		stats.players.total_score,
		stats.players.best_score,
		stats.players.avg_score,
		stats.players.num_sessions
	FROM stats.players
	ORDER BY stats.players.%s;
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
		stats.authors.num_quizzes,
    	stats.authors.avg_quiz_rate,
		stats.authors.best_quiz_rate
	FROM stats.authors
	WHERE stats.authors.user_id = $1;
	`
	var (
		num_quizzes    int32
		avg_quiz_rate  float32
		best_quiz_rate float32
	)
	err := r.pg.QueryRow(ctx, author_stat_query, user_id).Scan(&num_quizzes, &avg_quiz_rate, &best_quiz_rate)
	if err != nil {
		if err == pgx.ErrNoRows {
			return &api.AuthorStat{
				UserId:       user_id,
				NumQuizzes:   0,
				AvgQuizRate:  0,
				BestQuizRate: 0,
			}, nil
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
		stats.authors.user_id,
		stats.authors.num_quizzes,
    	stats.authors.avg_quiz_rate,
		stats.authors.best_quiz_rate
	FROM stats.authors
	ORDER BY stats.authors.%s;
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
