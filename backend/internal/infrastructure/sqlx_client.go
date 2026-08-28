// Package infrastructure
package infrastructure

import (
	"log/slog"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type SQLxClient struct {
	db *sqlx.DB
}

func NewSQLxClient(connStr string) SQLxClient {
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		slog.Error("sqlx connect failed", "error", err)
		os.Exit(1)
	}

	if err := db.Ping(); err != nil {
		slog.Error("sqlx ping failed", "error", err)
		os.Exit(1)
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)

	return SQLxClient{
		db: db,
	}
}

func (c SQLxClient) GetDB() *sqlx.DB {
	return c.db
}

func (c SQLxClient) Shutdown() error {
	if err := c.db.Close(); err != nil {
		return err
	}

	return nil
}
