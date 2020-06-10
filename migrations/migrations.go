package migrations

import (
	"github.com/gobuffalo/packr/v2"
	migrate "github.com/rubenv/sql-migrate"
)

type ReferenceType struct{}

func MySql() migrate.MigrationSource {
	return &migrate.PackrMigrationSource{
		Box: packr.New("mysql", "./mysql"),
	}
}

func Postgres() migrate.MigrationSource {
	return &migrate.PackrMigrationSource{
		Box: packr.New("postgres", "/postgres"),
	}
}
