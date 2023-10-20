package postgres

import (
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"analytics_collector/internal/config"
)

type Storage struct {
	db *sqlx.DB
}

const (
	healthcheckCount = 5
)

func New(cfg config.DBConfig) (*Storage, error) {
	const op = "storage.postgres.New"

	// open connection
	db, err := sqlx.Open("pgx", getURLFromConfig(cfg))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// healthcheck
	err = tryPingConnection(db, healthcheckCount)

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

func tryPingConnection(db *sqlx.DB, count int) error {
	var err error

	for count > 5 {
		err = db.Ping()
		if err != nil {
			count--
			time.Sleep(1 * time.Second)
		} else {
			return nil
		}
	}

	return err
}
