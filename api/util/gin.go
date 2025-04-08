package util

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http/httptest"
	"strings"
)

// CreateMockGinGetRequest creates a mock gin context with a GET request for testing purposes
func CreateMockGinGetRequest(path string) (*gin.Context, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest("GET", path, nil)

	return c, recorder
}

// CreateMockGinJsonRequest creates a mock gin context with a JSON POST request for testing purposes
func CreateMockGinJsonRequest(json string) *gin.Context {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/login", nil)

	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = ioutil.NopCloser(strings.NewReader(json))

	return c
}

func CreateMockGinJsonRequestWithRecorder(json string) (*gin.Context, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest("POST", "/login", nil)

	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = ioutil.NopCloser(strings.NewReader(json))

	return c, recorder
}
