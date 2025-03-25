// Package main contains only the main function, which is the entry point of the application.
package main

import (
	"fmt"
	"github.com/adraismawur/mibig-submission/db"
	"github.com/adraismawur/mibig-submission/endpoints"
	"github.com/adraismawur/mibig-submission/middleware"
	"github.com/gin-gonic/gin"
	"log/slog"
)

// main is the entry point of the application
func main() {
	// setup logging
	slog.Info("Starting MIBiG submission portal API")

	slog.Info("Setting up database")
	// setup database
	dbConnection := db.Connect()

	slog.Info("Setting up router")
	// setup router
	router := gin.Default()

	slog.Info("Registering middleware")
	router.Use(middleware.AuthMiddleware())

	slog.Info("Registering endpoints")
	endpoints.RegisterEndpointHandlers(router, dbConnection)

	slog.Info("Starting server")
	err := router.Run(":8080")
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to start server: %v", err))
	}
}
