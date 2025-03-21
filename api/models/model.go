// Package models contains the models for the application in use with gorm.
package models

import (
	"gorm.io/gorm"
	"log/slog"
)

func Migrate(db *gorm.DB) {
	slog.Info("[db] Migrating models")

	db.AutoMigrate(&User{}, &UserInfo{})

	slog.Info("[db] Done migrating models")
}
