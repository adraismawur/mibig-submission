package main

import (
	"github.com/adraismawur/mibig-submission/config"
	"github.com/adraismawur/mibig-submission/db"
	"github.com/adraismawur/mibig-submission/endpoints"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func registerRoute(router *gin.Engine, route endpoints.Route) {
	router.Handle(route.Method, route.Path, route.Handler)
}

func registerRoutes(router *gin.Engine, endpoint endpoints.Endpoint) {
	for _, route := range endpoint.Routes {
		registerRoute(router, route)
	}
}

func main() {
	// setup logging
	mainLog := log.New(os.Stdout, "[main] ", log.LstdFlags)
	mainLog.Println("Starting MIBiG submission portal API")

	// setup environment
	mainLog.Println("Setting up environment")
	config.Init()

	mainLog.Println("Setting up database")
	// setup database
	dbConnection := db.Connect()

	mainLog.Println("Setting up router")
	// setup router
	router := gin.Default()

	registerRoutes(router, endpoints.GetAuthEndpoint(dbConnection))
	registerRoutes(router, endpoints.GetUserEndpoint(dbConnection))
	registerRoutes(router, endpoints.GetSubmissionEndpoint(dbConnection))
	registerRoutes(router, endpoints.GetReviewEndpoint(dbConnection))

	mainLog.Println("Starting server")
	err := router.Run(":8080")
	if err != nil {
		mainLog.Fatalf("Failed to start server: %v", err)
	}
}
