package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"quiz_app/internal/config"
	api "quiz_app/pkg/api/v1"
	"quiz_app/pkg/postgres"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pg *pgxpool.Pool
}

func NewRepository(ctx context.Context, config *config.Config) (*Repository, error) {
	pg, err := postgres.New(ctx, config.Postgres)
	if err != nil {
		return nil, fmt.Errorf("unable to create a repository: %w", err)
	}
	return &Repository{
		pg: pg,
	}, nil
}

func (r *Repository) CreateSession(
	ctx context.Context,
	quiz_id string,
	start_time, end_time string,
	results []*api.Result,
) (string, error) {
	session_id := uuid.New().String()
	query := `INSERT INTO 
	sessions (session_id, quiz_id, start_time, end_time)
	VALUES ($1, $2, $3, $4);`
	_, err := r.pg.Exec(ctx, query,
		session_id,
		quiz_id,
		start_time,
		end_time,
	)
	if err != nil {
		return "", fmt.Errorf("unable to insert new session: %w", err)
	}
	for _, res := range results {
		for _, answer_id := range res.Answers {
			query := `INSERT INTO 
			user_answers (session_id, user_id, answer_id)
			VALUES ($1, $2, $3);`
			_, err := r.pg.Exec(ctx, query,
				session_id,
				res.UserId,
				answer_id,
			)
			if err != nil {
				return "", fmt.Errorf("unable to insert answer %s of user %s: %w", answer_id, res.UserId, err)
			}
		}
	}
	return session_id, nil
}

func (r *Repository) GetSession(ctx context.Context, session_id string) (*api.Session, error) {
	query := `SELECT 
				s.session_id, 
				s.quiz_id, 
				s.start_time, 
				s.end_time 
			FROM sessions s
			WHERE s.session_id=$1;`
	var (
		quiz_id    string
		start_time string
		end_time   string
		results    []*api.Result
	)
	err := r.pg.QueryRow(ctx, query, session_id).Scan(&quiz_id, &start_time, &end_time)
	if err != nil {
		return nil, fmt.Errorf("unable to find session %s: %w", session_id, err)
	}
	query = `SELECT 
				ua.user_id,
				array_agg(ua.answer_id) as answers
			FROM user_answers ua
			WHERE ua.session_id=$1
			GROUP BY ua.user_id;`
	rows, err := r.pg.Query(ctx, query, session_id)
	if err != nil {
		return nil, fmt.Errorf("unable to find users' answers %s: %w", session_id, err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			user_id string
			answers []string
		)
		err = rows.Scan(&user_id, &answers)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		results = append(results,
			&api.Result{
				UserId:  user_id,
				Answers: answers,
			})
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("query failed: %w", rows.Err().Error())
	}

	return &api.Session{
		SessionId: session_id,
		QuizId:    quiz_id,
		StartTime: start_time,
		EndTime:   end_time,
		Results:   results,
	}, nil
}

func (r *Repository) ListSessions(ctx context.Context) ([]*api.Session, error) {
	query := `
        SELECT 
    		s.session_id,
    		s.quiz_id::text,
    		s.start_time,
    		s.end_time,
    	COALESCE(
        	jsonb_agg(
            	jsonb_build_object(
                	'user_id', ua.user_id::text,
                	'answers', ua.answers
            	)
        	), '[]'::jsonb
    	) AS results
		FROM sessions s
		LEFT JOIN (
    		SELECT 
        		session_id,
        		user_id,
        		array_agg(answer_id) AS answers
    		FROM user_answers ua
    		GROUP BY session_id, user_id
		) ua ON s.session_id = ua.session_id
		GROUP BY s.session_id, s.quiz_id, s.start_time, s.end_time
		ORDER BY s.start_time DESC;
	`

	rows, err := r.pg.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var sessions []*api.Session
	for rows.Next() {
		var session api.Session
		var rawResults []byte

		err := rows.Scan(
			&session.SessionId,
			&session.QuizId,
			&session.StartTime,
			&session.EndTime,
			&rawResults,
		)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		err = json.Unmarshal(rawResults, &session.Results)
		if err != nil {
			return nil, fmt.Errorf("json decoding error: %w", err)
		}

		sessions = append(sessions, &session)
	}
	return sessions, nil
}
