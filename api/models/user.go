package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log/slog"
)

// Role type that represents a user's role
type Role int

const (
	Submitter Role = iota
	Reviewer
	Admin
)

// User model that represents a singular user
type User struct {
	gorm.Model `json:"-"`
	Email      string   `json:"email"`
	Password   string   `json:"password"`
	Active     bool     `json:"active"`
	Role       Role     `json:"role"`
	Info       UserInfo `json:"info"`
}

// HasRole returns true if a user has a specific role
func HasRole(db *gorm.DB, user User, role Role) bool {
	tx := db.First(&UserInfo{}, "user_id = ? AND role = ?", user.ID, role)
	return tx.RowsAffected > 0
}

func GetUserExistsByEmail(db *gorm.DB, email string) (bool, error) {

	tx := db.Exec(`SELECT EXISTS (SELECT 1 FROM users WHERE email = ?)`, email)

	if tx.Error != nil {
		return false, tx.Error
	}

	return tx.RowsAffected > 0, nil
}

func GetUserExistsByID(db *gorm.DB, id int) (bool, error) {
	var user User

	tx := db.Limit(1).Find(&user).Where("id = ?", id)

	if tx.Error != nil {
		return false, tx.Error
	}

	return tx.RowsAffected > 0, nil
}

func CreateUser(db *gorm.DB, email string, password string, role Role) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		slog.Error("[user] Could not hash password of new user")
		return err
	}

	user := User{
		Email:    email,
		Password: string(hashedPassword),
		Active:   true,
		Role:     role,
	}

	if err := db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func GetUser(db *gorm.DB, id int) (*User, error) {
	user := User{}

	tx := db.
		Where("id = ?", id).
		Find(&user)

	if tx.Error != nil {
		// other error
		return nil, tx.Error
	}

	// not found
	if tx.RowsAffected == 0 {
		return nil, nil
	}

	return &user, nil
}

func GetUsers(db *gorm.DB, offset int, limit int) ([]User, error) {
	var users []User

	tx := db.Find(&users).Offset(offset).Limit(limit)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return users, nil
		}

		return nil, tx.Error
	}

	return users, nil
}

func UpdateUser(db *gorm.DB, id int, user *User) error {
	tx := db.
		Model(&user).
		Where("id = ?", id).
		Omit("password").
		Updates(user)

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func UpdateUserPassword(db *gorm.DB, user *User, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		slog.Error("[user] Could not hash updated password")
		return err
	}

	user.Password = string(hashedPassword)

	if err := db.Save(user).Error; err != nil {
		return err
	}
	return nil
}

// LoginRequest type that represents a user request given by a client through a POST request
type LoginRequest struct {
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
