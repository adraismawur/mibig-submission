// Package main contains only the main function, which is the entry point of the application.
package main

import (
	"fmt"
	"github.com/adraismawur/mibig-submission/config"
	"github.com/adraismawur/mibig-submission/db"
	"github.com/adraismawur/mibig-submission/endpoints"
	"github.com/adraismawur/mibig-submission/middleware"
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/models/entry"
	"github.com/gin-gonic/gin"
	"log/slog"
	path2 "path"
)

// main is the entry point of the application
func main() {
	// setup logging
	slog.Info("Starting MIBiG entry portal API")

	slog.Info("Setting up database")
	// setup database
	dbConnection := db.Connect()

	slog.Info("Setting up router")
	// setup router
	router := gin.Default()

	outputDir := path2.Join(config.Envs["DATA_PATH"], "antismash")
	router.Static("/antismash/result", outputDir)

	slog.Info("Registering middleware")
	router.Use(middleware.AuthMiddleware())

	slog.Info("Registering endpoints")
	endpoints.RegisterEndpointHandlers(router, dbConnection)

	// populate the database if this is the first time we are starting it
	models.Populate(dbConnection)

	slog.Info("Preloading MIBiG entries")
	entry.PreloadMibigDatabase(dbConnection)

	slog.Info("Starting AntiSMASH runner goroutine")
	go models.AntismashWorker(dbConnection)

	slog.Info("Starting server")
	err := router.Run(config.Envs["SERVER_PORT"])
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to start server: %v", err))
	}
}
