package db

import (
	"fmt"
	"github.com/adraismawur/mibig-submission/config"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
)

func Connect() *gorm.DB {
	dbLog := log.New(os.Stdout, "[db] ", log.LstdFlags)
	dbLog.Println("Opening database connection")

	var err error
	var db *gorm.DB

	dialect := config.Envs["DB_DIALECT"]
	dbLog.Println("Dialect: ", dialect)

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
		dbLog.Panicf("Unsupported database dialect: %s", dialect)
	}

	if err != nil {
		dbLog.Panicf("Failed to connect to database: %v", err)
	}

	dbLog.Println("Database connection established")

	return db
}
