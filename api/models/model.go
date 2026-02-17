// Package models contains the models for the application in use with gorm.
package models

import (
	"gorm.io/gorm"
	"log/slog"
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

func Populate(db *gorm.DB) {
	// populate database with starting data if the relevant tables are empty
	// this is used to create a user when this is the first time we are starting the DB

	var count int64

	err := db.Table("database_meta").Count(&count).Error

	if err != nil {
		slog.Error("[db] Error getting database meta count")
		panic("Error getting database meta count")
	}

	var meta DatabaseMeta

	if count == 0 {
		meta = DatabaseMeta{}
	} else {
		err = db.Table("database_meta").First(&meta).Error

		if err != nil {
			slog.Error("[db] Populate error: ", "error", err)
			panic(err)
		}
	}

	if meta.FirstTimeSetupDone {
		return
	}

	for _, initDataEntry := range InitData {
		slog.Info("[db] Creating first time start data", "table", initDataEntry.Table)

		db.Create(initDataEntry.Model)
	}

	meta.FirstTimeSetupDone = true
}
