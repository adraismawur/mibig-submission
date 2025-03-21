// Package db provides a method to connect to a database using GORM.
package db

import (
	"fmt"
	"github.com/adraismawur/mibig-submission/config"
	"github.com/adraismawur/mibig-submission/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log/slog"
)

// Connect uses GORM to connect to a database based on the environment variables.
// This function relies on the user to have set the environment variables correctly.
// The function will panic if the dialect is not supported.
func Connect() *gorm.DB {
	slog.Info("[db] Opening database connection")

	var err error
	var db *gorm.DB

	dialect := config.Envs["DB_DIALECT"]
	slog.Info("[db] Dialect: ", dialect)

	// We want to be specific about what we are expecting to use
	// Postgres for production and SQLite for testing
	if dialect == "postgres" {
		connectionUrl := fmt.Sprintf(
			"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
			config.Envs["DB_HOST"],
			config.Envs["DB_PORT"],
			config.Envs["DB_DBNAME"],
			config.Envs["DB_USER"],
			config.Envs["DB_PASS"],
		)
		db, err = gorm.Open(postgres.Open(connectionUrl))
	} else if dialect == "sqlite" {
		db, err = gorm.Open(sqlite.Open(config.Envs["DB_PATH"]))
	} else {
		slog.Error(fmt.Sprintf("Unsupported database dialect: %s", dialect))
		panic(fmt.Sprintf("Unsupported database dialect: %s", dialect))
	}

	if err != nil {
		slog.Error("[db] Failed to connect to database: %v", err)
		panic(err)
	}

	slog.Info("[db] Database connection established")

	slog.Info("[db] Migrating models")
	models.Migrate(db)

	return db
}
