package functions

import (
	"fmt"
	"github.com/adraismawur/mibig-submission/config"
	"github.com/adraismawur/mibig-submission/db"
	"github.com/adraismawur/mibig-submission/endpoints"
	"github.com/adraismawur/mibig-submission/middleware"
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/models/entry"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log/slog"
	path2 "path"
)

func StartApi() {
	// setup logging
	slog.Info("Starting MIBiG entry portal API")

	slog.Info("Setting up database")
	// setup database
	dbConnection, err := db.Connect()

	if err != nil {
		slog.Error("[main] Could not connect to database")
		panic("Panic in main function: Could not connect to database")
	}

	slog.Info("Setting up router")
	// setup router
	router := gin.Default()

	dataPath, err := config.GetConfig(config.EnvDataPath)

	if err != nil {
		slog.Error("[main] Could not get env variable for data path")
		panic("Panic in main function: Could not get env variable for data path")
	}

	outputDir := path2.Join(dataPath, "antismash")
	router.Static("/antismash/result", outputDir)

	slog.Info("Registering middleware")
	router.Use(middleware.AuthMiddleware())

	// of cors, we use cors
	router.Use(cors.Default())

	slog.Info("Registering endpoints")
	endpoints.RegisterEndpointHandlers(router, dbConnection)

	// populate the database if this is the first time we are starting it
	err = models.Populate(dbConnection)

	if err != nil {
		slog.Error("[main] Could not populate database")
		panic("Panic in main function: Could not populate database")
	}

	slog.Info("Downloading MIBiG database")
	err = entry.DownloadMIBiGdatabase()

	if err != nil {
		slog.Error("[main] Could not download mibiG database")
		panic("Panic in main function: Could not download mibiG database")
	}

	slog.Info("Preloading MIBiG entries")
	err = entry.PreloadMibigDatabase(dbConnection)

	if err != nil {
		slog.Error("[main] Could not load mibiG entries into database")
		panic("Panic in main function: Could not load mibiG entries into database")
	}

	slog.Info("Starting AntiSMASH runner goroutine")
	go endpoints.AntismashWorker(dbConnection)

	slog.Info("Starting server")

	serverPort, err := config.GetConfig("SERVER_PORT")

	if err != nil {
		slog.Error("[main] Could not get env variable for server port")
		panic("Panic in main function: could not start server")
	}

	err = router.Run(serverPort)

	if err != nil {
		slog.Error(fmt.Sprintf("[main] Failed to start server: %v", err))
		panic("Panic in main function: could not start server")
	}
}
