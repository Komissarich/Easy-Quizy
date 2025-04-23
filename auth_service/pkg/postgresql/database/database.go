package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func New(ctx context.Context, username string, password string, host string, port int, name string) (*DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", username, password, host, port, name))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			username VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database schema: %w", err)
	}

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS friends (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			friend_id INTEGER NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users (id),
			FOREIGN KEY (friend_id) REFERENCES users (id)
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database schema: %w", err)

	}

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS favorite_quizzes (
			user_id BIGINT NOT NULL,
			quiz_id VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			PRIMARY KEY (user_id, quiz_id),
			CONSTRAINT fk_favorite_quizzes_user
				FOREIGN KEY (user_id) 
				REFERENCES users(id)
				ON DELETE CASCADE
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database schema: %w", err)
	}

	return &DB{DB: db}, nil
}
func (db *DB) Close(ctx context.Context) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		err = db.DB.Close()
		cancel()
	}()

	<-ctx.Done()
	return
}
