package main

import (
	"golang-family-tree/application"
	datasource "golang-family-tree/infra/db/datasource"
	migration "golang-family-tree/infra/db/migrations"

	log "github.com/sirupsen/logrus"
)

const (
	migrationsSourceUrl = "file://migrations"
	port                = ":8080"
)

func main() {
	runDBMigrations()
	starApplication()
}

func starApplication() {
	log.Info("Starting web application")
	app := application.NewApplication()
	err := app.Run(port)
	if err != nil {
		log.Fatalf("Unable to start web application. Reason: %v", err)
	}
}

func runDBMigrations() {
	log.Info("Starting database migration")
	db, err := datasource.GetDb()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	migration := migration.NewDBMigration(db)
	if err = migration.Run(migrationsSourceUrl); err != nil {
		log.Fatalf("Unable to run database migrations: %s", err)
	}
	log.Info("Database migration completed successfully")
}
