// Package models contains the models for the application in use with gorm.
package models

import (
	"gorm.io/gorm"
	"log/slog"
	"time"
)

var Models []interface{}

type InitDataEntry struct {
	Table string
	Model interface{}
}

var InitData []InitDataEntry

func Migrate(db *gorm.DB) {
	slog.Info("[db]", "Number of models", len(Models))

	// create tables and relations
	for _, model := range Models {
		err := db.AutoMigrate(model)
		if err != nil {
			panic(err)
		}
	}

	slog.Info("[db] Done migrating models")
}

func Populate(db *gorm.DB) error {
	// populate database with starting data if the relevant tables are empty
	// this is used to create a user when this is the first time we are starting the DB

	var count int64

	err := db.Table("application_states").Count(&count).Error

	if err != nil {
		slog.Error("[db] Error getting database state count")
		return err
	}

	var state ApplicationState

	if count == 0 {
		state = ApplicationState{
			FirstTimeSetupDone:       false,
			LastGroundTruthDownload:  time.Time{},
			NewSubmissionsEnabled:    true,
			SubmissionReviewsEnabled: true,
		}
	} else {
		err = db.Table("application_states").First(&state).Error

		if err != nil {
			slog.Error("[db] Populate error: ", "error", err)
			return err
		}
	}

	if state.FirstTimeSetupDone {
		return nil
	}

	for _, initDataEntry := range InitData {
		slog.Info("[db] Creating first time start data", "table", initDataEntry.Table)

		err = db.Create(initDataEntry.Model).Error

		if err != nil {
			slog.Error("[db] Populate error: ", "error", err)
			return err
		}
	}

	state.FirstTimeSetupDone = true

	err = db.Create(&state).Error

	if err != nil {
		slog.Error("[db] Error saving metadata: ", "error", err)
		return err
	}

	return nil
}
