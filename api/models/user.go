package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log/slog"
)

// Role type that represents a user's role
type Role uint

const (
	Submitter Role = iota
	Reviewer
	Admin
)

// User model that represents a singular user
type User struct {
	ID       uint       `json:"id"`
	Email    string     `json:"email"`
	Password string     `json:"password"`
	Active   bool       `json:"active"`
	Roles    []UserRole `json:"roles" gorm:"foreignKey:UserID"`
	Info     UserInfo   `json:"info"`
}

// LoginRequest type that represents a user request given by a client through a POST request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserRole model that represents a list of roles
type UserRole struct {
	Role   Role `json:"role"`
	UserID uint `json:"-"`
}

// UserInfo model that represents additional information about a user
type UserInfo struct {
	gorm.Model    `json:"-"`
	UserID        uint   `json:"user_id"`
	Alias         string `json:"alias"`
	Name          string `json:"name"`
	CallName      string `json:"call_name"`
	Organization1 string `json:"organization1"`
	Organization2 string `json:"organization2"`
	Organization3 string `json:"organization3"`
	OrcID         string `json:"orc_id"`
	Public        bool   `json:"public"`
}

// CreateUser creates a new user in the database with the given email, password and role
// the user is automatically set to active
func CreateUser(db *gorm.DB, email string, password string, roles []UserRole) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		slog.Error("[user] Could not hash password of new user")
		return err
	}

	user := User{
		Email:    email,
		Password: string(hashedPassword),
		Active:   true,
		Roles:    roles,
	}

	if err := db.Create(&user).Error; err != nil {
		return err
	}

	// append roles
	//for _, role := range roles {
	//	err := db.Model(&user).Association("Roles").Append(role)
	//	if err != nil {
	//		slog.Error("[user] Could not append role to user", "user", user.Role, "role", role.Role)
	//		return err
	//	}
	//}

	return nil
}

// GetUser returns a user from the database with the given Role
// this function returns an error if the user does not exist
func GetUser(db *gorm.DB, id int) (*User, error) {
	user := User{}

	tx := db.
		Preload("Roles").
		Preload("Info").
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

	tx := db.
		Preload("Roles").
		Preload("Info").
		Omit("password").
		Offset(offset).
		Limit(limit).
		Find(&users)

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

// GetUserExistsByID returns true if a user with a given Role exists
func GetUserExistsByID(db *gorm.DB, id int) (bool, error) {
	tx := db.Exec(`SELECT EXISTS (SELECT 1 FROM users WHERE id = ?)`, id)

	if tx.Error != nil {
		slog.Error("[user] Error getting user existence by id", "id", id, "error", tx.Error)
		return false, tx.Error
	}

	return tx.RowsAffected > 0, nil
}

// UpdateUser updates a user in the database with the given Role
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

// DeleteUser deletes a user from the database with the given Role
// performs a soft delete
func DeleteUser(db *gorm.DB, id int) error {
	tx := db.Delete(User{}, "id = ?", id)

	if tx.Error != nil {
		slog.Error("[user] Error deleting user", "id", id, "error", tx.Error)
		return tx.Error
	}

	return nil
}

// UpdateUserPassword updates the password of a user in the database with the given Role
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

// HasRole returns true if a user from a given Role has a specific role
func HasRole(user User, role Role) bool {
	for _, r := range user.Roles {
		if r.Role == role {
			return true
		}
	}

	return false
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
