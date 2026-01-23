package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log/slog"
	"math/rand"
	"strconv"
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
	ID       uint       `json:"id,omitempty"`
	Email    string     `json:"email"`
	Password string     `json:"password,omitempty"`
	Active   bool       `json:"active,omitempty"`
	Roles    []UserRole `json:"roles,omitempty" gorm:"foreignKey:UserID"`
	Info     UserInfo   `json:"info,omitempty"  gorm:"foreignKey:UserID"`
}

// LoginRequest type that represents a user request given by a client through a POST request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserRole model that represents a list of roles
type UserRole struct {
	ID     uint `json:"-" gorm:"primaryKey"`
	UserID uint `json:"user_id"`
	Role   Role `json:"role"`
}

// UserInfo model that represents additional information about a user
type UserInfo struct {
	ID            uint   `json:"-" gorm:"primaryKey"`
	UserID        uint   `json:"user_id"`
	Alias         string `json:"alias" gorm:"default:"`
	Name          string `json:"name" gorm:"default:"`
	CallName      string `json:"call_name" gorm:"default:"`
	Organization1 string `json:"organization1"`
	Organization2 string `json:"organization2"`
	Organization3 string `json:"organization3"`
	OrcID         string `json:"orc_id"`
	Public        bool   `json:"public" gorm:"default:false"`
}

func init() {
	Models = append(Models, &User{})
	Models = append(Models, &UserRole{})
	Models = append(Models, &UserInfo{})

	// password is 'changeme'
	defaultPassword, err := bcrypt.GenerateFromPassword([]byte("changeme"), bcrypt.DefaultCost)

	if err != nil {
		slog.Error("Could not generate default password", "error", err)
		return
	}

	InitData = append(InitData, InitDataEntry{
		Table: "users",
		Model: &User{
			ID:       1,
			Email:    "admin@localhost",
			Password: string(defaultPassword),
			Active:   true,
			Roles: []UserRole{
				{
					UserID: 1,
					Role:   Admin,
				},
			},
			Info: UserInfo{
				UserID:        1,
				Alias:         "test",
				Name:          "test",
				CallName:      "test",
				Organization1: "test",
				Organization2: "test",
				Organization3: "test",
				OrcID:         "",
				Public:        false,
			},
		},
	})
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
		Info:     UserInfo{},
	}

	if err := db.Create(&user).Error; err != nil {
		return err
	}

	if err = db.Save(&user).Error; err != nil {
		return err
	}

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
	var count int64

	tx := db.Table("users").
		Where("email = ?", email).
		Count(&count)

	if tx.Error != nil {
		slog.Error("[user] Error getting user existence by email", "email", email, "error", tx.Error)
		return false, tx.Error
	}

	return count > 0, nil
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
func UpdateUser(db *gorm.DB, id int, oldUser *User, newUser *User) error {
	tx := db.
		Model(oldUser).
		Where("id = ?", id).
		Omit("Password").
		Save(newUser)

	err := db.Model(&oldUser).
		Association("Info").
		Replace(&newUser.Info)

	if err != nil {
		slog.Error("[user] Error replacing user roles", "error", err)
		return err
	}

	// remove old roles
	err = db.Model(&oldUser).
		Association("Roles").
		Clear()

	if err != nil {
		slog.Error("[user] Error deleting user roles", "error", err)
		return err
	}

	// add new ones
	err = db.Model(&oldUser).
		Association("Roles").
		Append(&newUser.Roles)

	if err != nil {
		slog.Error("[user] Error replacing user roles", "error", err)
		return err
	}

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
	err := bcrypt.CompareHashAndPassword([]byte(against), []byte(in))

	if err != nil {
		return false
	}

	return true
}

func GenerateRandomEmail() string {
	return "test" + strconv.Itoa(rand.Intn(100000)) + "@localhost"
}
