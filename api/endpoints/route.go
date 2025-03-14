package endpoints

import (
	"github.com/gin-gonic/gin"
)

type Endpoint struct {
	Routes []Route
}

type Route struct {
	Method  string
	Path    string
	Handler func(c *gin.Context)
}
