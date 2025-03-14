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

func Connect() *gorm.DB {
	slog.Info("[db] Opening database connection")

	var err error
	var db *gorm.DB

	dialect := config.Envs["DB_DIALECT"]
	slog.Info("[db] Dialect: ", dialect)

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
