// Package main contains only the main function, which is the entry point of the application.
package main

import (
	"fmt"
	"github.com/adraismawur/mibig-submission/config"
	"github.com/adraismawur/mibig-submission/db"
	"github.com/adraismawur/mibig-submission/endpoints"
	"github.com/gin-gonic/gin"
	"log/slog"
)

// main is the entry point of the application
func main() {
	// setup logging
	slog.Info("Starting MIBiG submission portal API")

	// setup environment
	slog.Info("Setting up environment")
	config.Init()

	slog.Info("Setting up database")
	// setup database
	dbConnection := db.Connect()

	slog.Info("Setting up router")
	// setup router
	router := gin.Default()

	slog.Info("Registering endpoints")
	endpoints.RegisterEndpoints(router, dbConnection)

	slog.Info("Starting server")
	err := router.Run(":8080")
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to start server: %v", err))
	}
}
