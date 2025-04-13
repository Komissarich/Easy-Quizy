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
