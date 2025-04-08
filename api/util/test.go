package util

import (
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
)

// CreateMockDB creates a mock database connection for testing purposes
func CreateMockDB() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()

	if err != nil {
		slog.Error("[test] Could not create mock database connection")
		panic(err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	if err != nil {
		slog.Error("[test] Could not open mock database connection")
		panic(err)
	}

	return gormDB, mock
}
