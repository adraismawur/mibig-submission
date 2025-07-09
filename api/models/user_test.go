package models

import (
	"github.com/adraismawur/mibig-submission/util"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"testing"
)

var testDb *gorm.DB

const numTestUsers = 10

var testUsers []User

func TestMain(m *testing.M) {
	// setup
	testDb = util.CreateTestDB()
	Migrate(testDb)

	for i := 0; i < numTestUsers; i++ {
		userRoles := []UserRole{{Role: Admin}}

		user := User{
			Email:    util.GenerateRandomEmail(),
			Password: "test",
			Active:   true,
			Roles:    userRoles,
		}
		testDb.Create(&user)

		userRoles[0].UserID = user.ID

		testUsers = append(testUsers, user)
	}

	m.Run()

	// teardown
	testDb.Exec(`DELETE FROM users WHERE true`)
}

func TestCreateUser(t *testing.T) {
	testRoles := []UserRole{{Role: Admin}}
	err := CreateUser(testDb, "test@localhost", "test", testRoles)

	assert.NoError(t, err, "Error should be nil")

	var user User
	testDb.Preload("Roles").Last(&user)

	assert.Equal(t, "test@localhost", user.Email)
	assert.Equal(t, Admin, user.Roles[0].Role)
	assert.Nil(t, bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("test")))

	// cleanup after this test
	testDb.Delete(&user)
}

func TestGetUser(t *testing.T) {
	testRoles := []UserRole{{Role: Admin}}
	user := User{
		Email:    "test@localhost",
		Password: "test",
		Active:   true,
		Roles:    testRoles,
	}
	testDb.Create(&user)

	testDb.Preload("Roles").Last(&user)
	id := user.ID

	user2, err := GetUser(testDb, int(id))

	assert.Nil(t, err)

	assert.Equal(t, user.ID, user2.ID)
	assert.Equal(t, user.Email, user2.Email)
	assert.Equal(t, len(user.Roles), len(user2.Roles))
	assert.Equal(t, user.Roles[0], user2.Roles[0])
	assert.Equal(t, user.Active, user2.Active)
	// important: do not return the password
	assert.Equal(t, "", user2.Password)

	// cleanup
	testDb.Delete(&user)
}

func TestGetUsers(t *testing.T) {
	users, err := GetUsers(testDb, 0, 10)

	assert.Nil(t, err)

	assert.Equal(t, numTestUsers, len(users))

	for idx, user := range users {
		expectedUser := testUsers[idx]

		assert.Equal(t, expectedUser.Email, user.Email)
		assert.Equal(t, expectedUser.Roles[0], user.Roles[0])
		assert.Equal(t, expectedUser.Active, user.Active)
		// important: do not return the password
		assert.Equal(t, "", user.Password)
	}
}

func TestGetUsersWithOffset(t *testing.T) {
	expectedTestUsers := 5

	users, err := GetUsers(testDb, 5, 10)

	assert.Nil(t, err)

	assert.Equal(t, expectedTestUsers, len(users))
}

func TestGetUsersWithLimit(t *testing.T) {
	expectedTestUsers := 7

	users, err := GetUsers(testDb, 0, expectedTestUsers)

	assert.Nil(t, err)

	assert.Equal(t, expectedTestUsers, len(users))
}

func TestUpdateUserFull(t *testing.T) {
	testRoles := []UserRole{{Role: Admin}}
	testInfo := UserInfo{
		Alias:         "testalias",
		Name:          "testname",
		CallName:      "testcallname",
		Organization1: "testorg1",
		Organization2: "testorg2",
		Organization3: "testorg3",
		OrcID:         "testorcid",
		Public:        true,
	}

	user := User{
		Email:    "test@localhost",
		Password: "$2a$10$wOLM7A7gHgQXKKnyZX2J.uWi41KZKd.vfzKqa.w.9hUVFGVk.4LB.",
		Active:   true,
		Roles:    testRoles,
		Info:     testInfo,
	}

	err := testDb.Create(&user).Error

	assert.NoError(t, err)

	expectedRoles := []UserRole{{Role: Submitter}}
	expectedInfos := UserInfo{
		Alias:         user.Info.Alias + "_update",
		Name:          user.Info.Name + "_update",
		CallName:      user.Info.CallName + "_update",
		Organization1: user.Info.Organization1 + "_update",
		Organization2: user.Info.Organization2 + "_update",
		Organization3: user.Info.Organization3 + "_update",
		OrcID:         user.Info.OrcID + "_update",
		Public:        false,
	}

	expectedUser := User{
		ID:     user.ID,
		Email:  user.Email + "_update",
		Active: false,
		Roles:  expectedRoles,
		Info:   expectedInfos,
	}

	err = UpdateUser(testDb, int(user.ID), &expectedUser)

	assert.NoError(t, err)

	// get user
	var actualUser User
	testDb.
		Preload("Roles").
		Preload("Info").
		Omit("Password").
		Last(&actualUser)

	// assert base changes
	assert.Equal(t, expectedUser.Email, actualUser.Email)
	assert.Equal(t, expectedUser.Active, actualUser.Active)

	// assert role changes
	assert.Equal(t, expectedUser.Roles, actualUser.Roles)

	// assert info changes
	assert.Equal(t, expectedUser.Info.Alias, actualUser.Info.Alias)
	assert.Equal(t, expectedUser.Info.Name, actualUser.Info.Name)
	assert.Equal(t, expectedUser.Info.CallName, actualUser.Info.CallName)
	assert.Equal(t, expectedUser.Info.Organization1, actualUser.Info.Organization1)
	assert.Equal(t, expectedUser.Info.Organization2, actualUser.Info.Organization2)
	assert.Equal(t, expectedUser.Info.Organization3, actualUser.Info.Organization3)
	assert.Equal(t, expectedUser.Info.OrcID, actualUser.Info.OrcID)
	assert.Equal(t, expectedUser.Info.Public, actualUser.Info.Public)

	// cleanup after this test
	testDb.Delete(&user)
}
