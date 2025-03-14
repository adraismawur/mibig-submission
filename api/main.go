package main

import (
	"github.com/adraismawur/mibig-submission/config"
	"github.com/adraismawur/mibig-submission/db"
	"github.com/adraismawur/mibig-submission/endpoints"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

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

	mainLog.Println("Registering endpoints")
	endpoints.RegisterEndpoints(router, dbConnection)

	mainLog.Println("Starting server")
	err := router.Run(":8080")
	if err != nil {
		mainLog.Fatalf("Failed to start server: %v", err)
	}
}
