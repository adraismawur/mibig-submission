package util

import (
	"database/sql/driver"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

var TestPassword = "$2a$10$wOLM7A7gHgQXKKnyZX2J.uWi41KZKd.vfzKqa.w.9hUVFGVk.4LB." // this resolves to "test"

// CreateTestDB creates a mock database connection for testing purposes
func CreateTestDB() *gorm.DB {
	// Create a new SQLite database in memory
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	return db
}

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type AnyString struct{}

func (a AnyString) Match(v driver.Value) bool {
	_, ok := v.(string)
	return ok
}
