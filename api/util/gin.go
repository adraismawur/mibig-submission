package util

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http/httptest"
	"strings"
)

// CreateTestGinGetRequest creates a new gin context with a GET request for testing purposes
func CreateTestGinGetRequest(path string) (*gin.Context, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest("GET", path, nil)

	return c, recorder
}

// CreateTestGinJsonRequest creates a new gin context with a JSON POST request for testing purposes
func CreateTestGinJsonRequest(json string) *gin.Context {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/login", nil)

	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = ioutil.NopCloser(strings.NewReader(json))

	return c
}

// CreateTestGinJsonRequestWithRecorder creates a new gin context with a JSON POST request for testing purposes
// this also supplies a recorder so that the response can be checked
func CreateTestGinJsonRequestWithRecorder(json string) (*gin.Context, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest("POST", "/login", nil)

	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = ioutil.NopCloser(strings.NewReader(json))

	return c, recorder
}

// AddTokenToHeader adds a token to the request header for a given gin context
func AddTokenToHeader(c *gin.Context, token string) {
	c.Request.Header.Set("Authorization", "Bearer "+token)
}
