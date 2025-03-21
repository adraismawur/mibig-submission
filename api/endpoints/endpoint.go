// Package endpoints contains all the endpoints of the API
package endpoints

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Endpoint represents a collection of routes for a given endpoint, each with a method, path, and handler
type Endpoint struct {
	Routes []Route
}

// Route represents a single route with a method, path, and handler
type Route struct {
	Method  string
	Path    string
	Handler func(c *gin.Context)
}

func registerRoutes(router *gin.Engine, routes []Route) {
	for _, route := range routes {
		router.Handle(route.Method, route.Path, route.Handler)
	}
}

func registerEndpoint(router *gin.Engine, endpoints ...Endpoint) {
	for _, endpoint := range endpoints {
		registerRoutes(router, endpoint.Routes)
	}
}

// RegisterEndpoints registers all the endpoints of the API. This will grow as more endpoints are added
func RegisterEndpoints(router *gin.Engine, db *gorm.DB) {
	registerEndpoint(router, GetAuthEndpoint(db))
	registerEndpoint(router, GetUserEndpoint(db))
	registerEndpoint(router, GetSubmissionEndpoint(db))
	registerEndpoint(router, GetReviewEndpoint(db))
}
