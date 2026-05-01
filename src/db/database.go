package db

import (
	"database/sql"
	"fmt"
	"log"
	"speechToText/src/config"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func NewDB(cfg *config.DatabaseConfig) (*sql.DB, error) {
	connURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	conn, err := sql.Open("postgres", connURL)
	if err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(5 * time.Minute)

	if err := conn.Ping(); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("db ping: %w", err)
	}

	m, err := migrate.New(cfg.MigrationPath, connURL)
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("migrate init: %w", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			log.Printf("migrate source close: %v", srcErr)
		}
		if dbErr != nil {
			log.Printf("migrate db close: %v", dbErr)
		}
	}()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		_ = conn.Close()
		return nil, fmt.Errorf("migrate up: %w", err)
	}
	log.Println("Database migrations completed successfully")
	return conn, nil
}
