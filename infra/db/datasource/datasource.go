package timescaledb

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

var (
	dbHost           = os.Getenv("DB_HOST")
	dbPort           = os.Getenv("DB_PORT")
	dbName           = os.Getenv("DB_NAME")
	dbUser           = os.Getenv("DB_USER")
	dbPassword       = os.Getenv("DB_PASSWORD")
	dbMaxConnections = os.Getenv("DB_MAX_CONNECTIONS")
)

// GetDb return new datasource
func GetDb() (*sql.DB, error) {

	log.WithField("host", dbHost).WithField("port", dbPort).WithField("database", dbName).Info("Connecting to database")

	maxConnections, err := strconv.Atoi(dbMaxConnections)

	if err != nil {
		return nil, fmt.Errorf("unable to parse DB_MAX_CONNECTIONS: %w", err)
	}

	url := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	db, err := sql.Open("postgres", url)
	db.SetMaxOpenConns(maxConnections)
	if err != nil {
		return nil, fmt.Errorf("unable to open database connection: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return db, nil

}
