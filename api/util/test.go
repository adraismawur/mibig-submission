package util

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"log/slog"
	"net/http/httptest"
	"strings"
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

// CreateMockGinGetRequest creates a mock gin context with a GET request for testing purposes
func CreateMockGinGetRequest(path string) *gin.Context {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", path, nil)

	return c
}

// CreateMockGinJsonRequest creates a mock gin context with a JSON POST request for testing purposes
func CreateMockGinJsonRequest(json string) *gin.Context {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/login", nil)

	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = ioutil.NopCloser(strings.NewReader(json))

	return c
}
