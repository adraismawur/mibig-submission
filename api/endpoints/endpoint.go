package endpoints

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Endpoint struct {
	Routes []Route
}

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

func RegisterEndpoints(router *gin.Engine, db *gorm.DB) {
	registerEndpoint(router, GetAuthEndpoint(db))
	registerEndpoint(router, GetUserEndpoint(db))
	registerEndpoint(router, GetSubmissionEndpoint(db))
	registerEndpoint(router, GetReviewEndpoint(db))
}
