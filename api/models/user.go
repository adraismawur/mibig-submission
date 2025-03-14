package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log/slog"
)

type User struct {
	gorm.Model `json:"-"`
	Email      string   `json:"email"`
	Password   string   `json:"-"`
	Active     bool     `json:"active"`
	Role       UserRole `json:"role"`
	Info       UserInfo `json:"info"`
}

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CheckPassword(in string, against string) bool {
	hash, err := bcrypt.GenerateFromPassword([]byte(in), bcrypt.DefaultCost)

	if err != nil {
		slog.Error("[auth] Could not hash password")
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(against), hash)

	return true
}

type Role int

const (
	Admin Role = iota
	Submitter
	Reviewer
)

type UserRole struct {
	gorm.Model `json:"-"`
	UserID     uint `json:"user_id"`
	Role       Role `json:"role"`
}

func HasRole(db *gorm.DB, user User, role Role) bool {
	tx := db.First(&UserRole{}, "user_id = ? AND role = ?", user.ID, role)
	return tx.RowsAffected > 0
}

type UserInfo struct {
	gorm.Model `json:"-"`
	UserID     uint `json:"user_id"`
}
