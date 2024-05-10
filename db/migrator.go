package db

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

// Needed to add *.go below
//
//go:embed migrations/*.sql *.go
var schemaFS embed.FS

func Migrate(db *sql.DB) error {
	driver := "postgres"

	goose.SetBaseFS(schemaFS)

	if err := goose.SetDialect(driver); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil

}
