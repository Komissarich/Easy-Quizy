package postgres

import (
	"context"
	"fmt"

	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string `yaml:"POSTGRES_HOST" env:"POSTGRES_HOST" env-default:"127.0.0.1"`
	Port     uint16 `yaml:"POSTGRES_PORT" env:"POSTGRES_PORT" env-default:"5432"`
	Username string `yaml:"POSTGRES_USER" env:"POSTGRES_USERT" env-default:"postgres"`
	Password string `yaml:"POSTGRES_PASS" env:"POSTGRES_PASS" env-default:"root"`
	Database string `yaml:"POSTGRES_DB" env:"POSTGRES_DB" env-default:"postgres"`

	MinConns int32 `yaml:"POSTGRES_MIN_CONN" env:"POSTGRES_MIN_CONN" env-default:"1"`
	MaxConns int32 `yaml:"POSTGRES_MAX_CONN" env:"POSTGRES_MAX_CONN" env-default:"10"`
}

func New(ctx context.Context, config Config) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&pool_max_conns=%d&pool_min_conns=%d",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.MaxConns,
		config.MinConns,
	)
	conn, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	// migration, err := migrate.New(
	// 	"file://db/migrations",
	// 	fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
	// 		config.Username,
	// 		config.Password,
	// 		config.Host,
	// 		config.Port,
	// 		config.Database,
	// 	))
	// if err != nil {
	// 	return nil, fmt.Errorf("unable to create migrations: %w", err)
	// }
	// if err := migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
	// 	return nil, fmt.Errorf("unable to run migrations: %w", err)
	// }

	return conn, nil
}
func InitTables(ctx context.Context, pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(ctx)

	if err != nil {
		return fmt.Errorf("unable to get connection from pool")
	}

	query := `CREATE TABLE IF NOT EXISTS quizzes (
    Quiz_ID VARCHAR(255) PRIMARY KEY NOT NULL ,
    Name VARCHAR(255) NOT NULL ,
    Author VARCHAR(255) NOT NULL,
    Image_ID VARCHAR(255),
    Description VARCHAR(255)
);`
	conn.Exec(ctx, query)
	query = `CREATE TABLE IF NOT EXISTS questions(
Question_ID VARCHAR(255) PRIMARY KEY NOT NULL,
Quiz_ID VARCHAR(255) NOT NULL,
FOREIGN KEY (Quiz_ID) references quizzes(Quiz_ID),
Question_text VARCHAR(255),
Image_ID VARCHAR(255)
);`
	conn.Exec(ctx, query)

	query = `CREATE TABLE IF NOT EXISTS answers(
    Answer_ID VARCHAR(255) PRIMARY KEY NOT NULL,
    Question_ID VARCHAR(255) NOT NULL,
    FOREIGN KEY (Question_ID) REFERENCES questions(Question_ID),
    Answer_text VARCHAR(255),
    Is_correct BOOLEAN
);`
	conn.Exec(ctx, query)

	return nil

}
