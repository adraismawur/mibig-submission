package main

import (
	"github.com/adraismawur/mibig-submission/endpoints"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func registerRoute(router *gin.Engine, route endpoints.Route) {
	router.Handle(route.Method, route.Path, route.Handler)
}

func registerRoutes(router *gin.Engine, routes []endpoints.Route) {
	for _, route := range routes {
		registerRoute(router, route)
	}
}

func main() {
	mainLog := log.New(os.Stdout, "[main] ", log.LstdFlags)
	mainLog.Println("Starting MIBiG submission portal API")

	router := gin.Default()

	registerRoutes(router, endpoints.GetAuthRoutes())
	registerRoutes(router, endpoints.GetUserRoutes())
	registerRoutes(router, endpoints.GetSubmissionRoutes())
	registerRoutes(router, endpoints.GetReviewRoutes())

	err := router.Run(":8080")
	if err != nil {
		mainLog.Fatalf("Failed to start server: %v", err)
	}
}
