package migrations

import (
	"database/sql"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"
)

type DBMigration interface {
	Run(sourceUrl string) error
}

type dbMigrationImpl struct {
	db *sql.DB
}

// NewDBMigration create DBMigration
func NewDBMigration(db *sql.DB) DBMigration {
	return dbMigrationImpl{
		db: db,
	}
}

// Run execute database migrations
func (t dbMigrationImpl) Run(sourceUrl string) error {
	log.Info("Running migration scripts")
	driver, err := postgres.WithInstance(t.db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		sourceUrl,
		"postgres", driver)

	if err != nil {
		return err
	}

	err = m.Up()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
