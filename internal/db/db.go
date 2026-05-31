package db

import (
	"database/sql"
	"log/slog"

	_ "modernc.org/sqlite"

	"github.com/pressly/goose/v3"
)

func OpenAndMigrate(path string) *sql.DB {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		slog.Error("failed to open database", "error", err)
		panic(err)
	}

	if err := goose.SetDialect("sqlite3"); err != nil {
		slog.Error("goose set dialect error", "error", err)
		panic(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		slog.Error("goose migration error", "error", err)
		panic(err)
	}

	return db
}
