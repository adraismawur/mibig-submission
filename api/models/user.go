package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log/slog"
)

// Role type that represents a user's role
type Role int

const (
	Admin Role = iota
	Submitter
	Reviewer
)

// User model that represents a singular user
type User struct {
	gorm.Model `json:"-"`
	Email      string   `json:"email"`
	Password   string   `json:"-"`
	Active     bool     `json:"active"`
	Role       Role     `json:"role"`
	Info       UserInfo `json:"info"`
}

// HasRole returns true if a user has a specific role
func HasRole(db *gorm.DB, user User, role Role) bool {
	tx := db.First(&UserInfo{}, "user_id = ? AND role = ?", user.ID, role)
	return tx.RowsAffected > 0
}

// UserRequest type that represents a user request given by a client through a POST request
type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CheckPassword compares a plain text password with a hashed password and returns true if they match
func CheckPassword(in string, against string) bool {
	hash, err := bcrypt.GenerateFromPassword([]byte(in), bcrypt.DefaultCost)

	if err != nil {
		slog.Error("[auth] Could not hash password")
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(against), hash)

	return true
}

// UserInfo model that represents additional information about a user
type UserInfo struct {
	gorm.Model `json:"-"`
	UserID     uint `json:"user_id"`
}
