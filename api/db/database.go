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

var DB *gorm.DB
var dbLog = log.New(os.Stdout, "[db] ", log.LstdFlags)

func Connect() {
	dbLog.Println("Opening database connection")

	var err error

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
		DB, err = gorm.Open(postgres.Open(connectionUrl))
	} else if dialect == "sqlite" {
		DB, err = gorm.Open(sqlite.Open(config.Envs["DB_PATH"]))
	} else {
		dbLog.Panicf("Unsupported database dialect: %s", dialect)
	}

	if err != nil {
		dbLog.Panicf("Failed to connect to database: %v", err)
	}
}
