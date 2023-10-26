package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"analytics_collector/internal/config"
	"analytics_collector/internal/storage"
)

type Storage struct {
	db *sqlx.DB
}

const (
	healthcheckCount = 5
)

func New(ctx context.Context, cfg config.DBConfig) (*Storage, error) {
	const op = "storage.postgres.New"

	// open connection
	db, err := sqlx.Open("pgx", getURLFromConfig(cfg))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// close connection
	go func() {
		<-ctx.Done()

		err := db.Close()
		if err != nil {
			log.Printf("DB connection can't be closed: %s", err)
		}
		log.Printf("DB connection was closed")
	}()

	// healthcheck
	err = tryPingConnection(ctx, db, healthcheckCount)
	if err != nil {
		return nil, fmt.Errorf("%s: Ping database error: %w. Database URL: %s", op, err, getURLFromConfig(cfg))
	}

	return &Storage{db: db}, nil
}

func getURLFromConfig(cfg config.DBConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.Username,
		cfg.Password,
		cfg.Address,
		cfg.DBName,
		cfg.SSLMode,
	)
}

func tryPingConnection(ctx context.Context, db *sqlx.DB, count int) error {
	var err error

	for count > 0 {
		err = db.PingContext(ctx)
		if err != nil {
			count--
			time.Sleep(1 * time.Second)
		} else {
			return nil
		}
	}

	return err
}

func (s *Storage) Save(ctx context.Context, info storage.UserActionInfo) error {
	const op = "storage.postgres.Save"

	stmt, err := s.db.PrepareContext(ctx, `INSERT INTO metrics(query_time, user_id, query_data) VALUES ($1, $2, $3)`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, info.Time, info.UserID, info.Data)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
