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
type Role string

const (
	Submitter Role = "submitter"
	Reviewer  Role = "reviewer"
	Admin     Role = "admin"
)

// User model that represents a singular user
type User struct {
	ID        uint64     `json:"db_id"`
	Anonymous bool       `json:"anonymous"`
	Email     string     `json:"email" gorm:"unique"`
	Password  string     `json:"password,omitempty"`
	Active    bool       `json:"active"`
	Roles     []UserRole `json:"roles" gorm:"foreignKey:UserID"`
	Info      UserInfo   `json:"info,omitempty"  gorm:"foreignKey:UserID"`
}

// LoginRequest type that represents a user request given by a client through a POST request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PasswordChallenge struct {
	ID        uint64 `json:"-"`
	Email     string `json:"email"`
	Challenge string `json:"challenge"`
}

// UserRole model that represents a list of roles
type UserRole struct {
	ID     uint64 `json:"-" gorm:"primaryKey"`
	UserID uint64 `json:"user_id" gorm:"uniqueIndex:role_idx"`
	Role   Role   `json:"role" gorm:"uniqueIndex:role_idx"`
}

// UserInfo model that represents additional information about a user
type UserInfo struct {
	ID            uint64 `json:"db_id" gorm:"primaryKey"`
	UserID        uint64 `json:"db_user_id"`
	Alias         string `json:"alias"`
	Name          string `json:"name"`
	CallName      string `json:"call_name"`
	Organisation1 string `json:"organisation_1"`
	Organisation2 string `json:"organisation_2"`
	Organisation3 string `json:"organisation_3"`
	OrcID         string `json:"orc_id"`
	Public        bool   `json:"public"`
}

func init() {
	Models = append(Models, &User{})
	Models = append(Models, &UserRole{})
	Models = append(Models, &UserInfo{})
	Models = append(Models, &PasswordChallenge{})

	// password is 'changeme'
	defaultPassword, err := bcrypt.GenerateFromPassword([]byte("changeme"), bcrypt.DefaultCost)

	if err != nil {
		slog.Error("Could not generate default password", "error", err)
		return
	}

	InitData = append(InitData, InitDataEntry{
		Table: "users",
		Model: &User{
			Email:    "admin@localhost",
			Password: string(defaultPassword),
			Active:   false,
			Roles: []UserRole{
				{
					Role: Admin,
				},
				{
					Role: Submitter,
				},
				{
					Role: Reviewer,
				},
			},
			Info: UserInfo{
				Alias:         "Admin",
				Name:          "Admin",
				CallName:      "Admin",
				Organisation1: "",
				Organisation2: "",
				Organisation3: "",
				OrcID:         "",
			},
		},
	})

	InitData = append(InitData, InitDataEntry{
		Table: "users",
		Model: &User{
			Email:    "reviewer@localhost",
			Password: string(defaultPassword),
			Active:   false,
			Roles: []UserRole{
				{
					Role: Reviewer,
				},
			},
			Info: UserInfo{
				Alias:         "Reviewer",
				Name:          "Reviewer",
				CallName:      "Reviewer",
				Organisation1: "",
				Organisation2: "",
				Organisation3: "",
				OrcID:         "",
			},
		},
	})

	InitData = append(InitData, InitDataEntry{
		Table: "users",
		Model: &User{
			Email:    "submitter@localhost",
			Password: string(defaultPassword),
			Active:   false,
			Roles: []UserRole{
				{
					Role: Submitter,
				},
			},
			Info: UserInfo{
				Alias:         "Submitter",
				Name:          "Submitter",
				CallName:      "Submitter",
				Organisation1: "",
				Organisation2: "",
				Organisation3: "",
				OrcID:         "",
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
		Active:   false,
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
		Where("id = $1", id).
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
func GetUsers(db *gorm.DB, offset int, limit int, search string) ([]User, error) {
	var users []User

	tx := db.
		Preload("Roles").
		Preload("Info").
		Omit("password").
		Joins("LEFT OUTER JOIN user_infos ON user_infos.user_id = users.id")

	if search != "" {
		tx = tx.Where("user_infos.name LIKE $1", "%"+search+"%")
	}

	tx = tx.Offset(offset).
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

	err := db.Table("users").
		Where("email = $1", email).
		Count(&count).
		Error

	if err != nil {
		slog.Error("[user] Error getting user existence by email", "email", email, "error", err.Error())
		return false, err
	}

	return count > 0, nil
}

// GetUserExistsByID returns true if a user with a given Role exists
func GetUserExistsByID(db *gorm.DB, id int) (bool, error) {
	tx := db.Exec(`SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`, id)

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
func UpdateUser(db *gorm.DB, id int, newUser User) error {
	err := db.Session(&gorm.Session{FullSaveAssociations: true}).Transaction(func(tx *gorm.DB) error {
		if newUser.Roles == nil {
			tx.Omit("Roles")
		}

		err := tx.
			Omit("Password"). // password done separately
			//Where("id = $1", id).
			Select("*").
			Save(&newUser).Error

		if err != nil {
			slog.Error("[user] Error saving user info")
			return err
		}

		if newUser.Password == "" {
			return nil
		}

		err = UpdateUserPassword(tx, id, newUser.Password)

		if err != nil {
			slog.Error("[user] Error changing user password")
			return err
		}

		return nil

	})

	return err
}

// DeleteUser deletes a user from the database with the given Role
// performs a soft delete
func DeleteUser(db *gorm.DB, id int) error {
	tx := db.Delete(User{}, "id = $1", id)

	if tx.Error != nil {
		slog.Error("[user] Error deleting user", "id", id, "error", tx.Error)
		return tx.Error
	}

	return nil
}

// UpdateUserPassword updates the password of a user in the database with the given Role
// this function will hash the password before updating it as happens in the CreateUser function
// this function will **not** update any other field
func UpdateUserPassword(db *gorm.DB, userId int, plainPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)

	if err != nil {
		slog.Error("[user] Could not hash updated password")
		return err
	}

	hashedPasswordString := string(hashedPassword)

	err = db.
		Model(&User{}).
		Select("password").
		Where("id = $1", userId).
		Update("password", &hashedPasswordString).
		Error

	if err != nil {
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
