// Package models contains the models for the application in use with gorm.
package models

import (
	"gorm.io/gorm"
	"log/slog"
)

var Models []interface{}

func Migrate(db *gorm.DB) {
	slog.Info("[db]", "Number of models", len(Models))

	for _, model := range Models {
		err := db.AutoMigrate(model)
		if err != nil {
			panic(err)
		}
	}

	slog.Info("[db] Done migrating models")
}
