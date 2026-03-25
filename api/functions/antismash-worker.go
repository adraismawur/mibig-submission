package functions

import (
	"github.com/adraismawur/mibig-submission/db"
	"github.com/adraismawur/mibig-submission/endpoints"
	"log/slog"
)

func StartAntismashWorker() {
	dbConnection, err := db.Connect()

	if err != nil {
		slog.Error("[main] Could not connect to database")
		panic("Panic in main function: Could not connect to database")
	}

	slog.Info("Starting AntiSMASH runner goroutine")
	endpoints.AntismashWorker(dbConnection)
}
