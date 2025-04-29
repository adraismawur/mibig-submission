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

// LoginRequest type that represents a user request given by a client through a POST request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserInfo model that represents additional information about a user
type UserInfo struct {
	gorm.Model `json:"-"`
	UserID     uint `json:"user_id"`
}

// CreateUser creates a new user in the database with the given email, password and role
// the user is automatically set to active
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

// GetUser returns a user from the database with the given ID
// this function returns an error if the user does not exist
func GetUser(db *gorm.DB, id int) (*User, error) {
	user := User{}

	tx := db.
		Where("id = ?", id).
		Omit("password").
		Find(&user)

	if tx.Error != nil {
		// other error
		slog.Error("[user] Error getting user", "id", id, "error", tx.Error)
		return nil, tx.Error
	}

	// not found
	if tx.RowsAffected == 0 {
		return nil, nil
	}

	return &user, nil
}

// GetUsers returns all users from the database with the given offset and limit
func GetUsers(db *gorm.DB, offset int, limit int) ([]User, error) {
	var users []User

	tx := db.Omit("password").Offset(offset).Limit(limit).Find(&users)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			slog.Error("[user] Error getting all users", "offset", offset, "limit", limit, "error", tx.Error)
			return users, nil
		}

		return nil, tx.Error
	}

	return users, nil
}

// GetUserExistsByEmail returns true if a user with a given email exists
func GetUserExistsByEmail(db *gorm.DB, email string) (bool, error) {
	tx := db.Exec(`SELECT EXISTS (SELECT 1 FROM users WHERE email = ?)`, email)

	if tx.Error != nil {
		slog.Error("[user] Error getting user existence by email", "email", email, "error", tx.Error)
		return false, tx.Error
	}

	return tx.RowsAffected > 0, nil
}

// GetUserExistsByID returns true if a user with a given ID exists
func GetUserExistsByID(db *gorm.DB, id int) (bool, error) {
	tx := db.Exec(`SELECT EXISTS (SELECT 1 FROM users WHERE id = ?)`, id)

	if tx.Error != nil {
		slog.Error("[user] Error getting user existence by id", "id", id, "error", tx.Error)
		return false, tx.Error
	}

	return tx.RowsAffected > 0, nil
}

// UpdateUser updates a user in the database with the given ID
// this will overwrite everything in the given user struct
// so if anything is set to default values, it will be overwritten.
// this function will **not** update the password
func UpdateUser(db *gorm.DB, id int, user *User) error {
	tx := db.
		Model(&user).
		Where("id = ?", id).
		Omit("password").
		Updates(user)

	if tx.Error != nil {
		slog.Error("[user] Error updating user", "id", id, "error", tx.Error)
		return tx.Error
	}
	return nil
}

// DeleteUser deletes a user from the database with the given ID
// performs a soft delete
func DeleteUser(db *gorm.DB, id int) error {
	tx := db.Delete(User{}, "id = ?", id)

	if tx.Error != nil {
		slog.Error("[user] Error deleting user", "id", id, "error", tx.Error)
		return tx.Error
	}

	return nil
}

// UpdateUserPassword updates the password of a user in the database with the given ID
// this function will hash the password before updating it as happens in the CreateUser function
// this function will **not** update any other field
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

// HasRole returns true if a user from a given ID has a specific role
func HasRole(db *gorm.DB, userId int, role Role) bool {
	var user User
	tx := db.First(&user, "user_id = ?", userId)

	if tx.Error != nil {
		slog.Error("[user] Error getting user role", "userId", userId, "error", tx.Error)
		return false
	}

	return user.Role == role
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
