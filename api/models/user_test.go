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

	// create some test users

	for i := 0; i < numTestUsers; i++ {
		user := User{
			Email:    util.GenerateRandomEmail(),
			Password: "test",
			Active:   true,
			Role:     Admin,
		}
		testDb.Create(&user)
		testUsers = append(testUsers, user)
	}

	m.Run()

	// teardown
	testDb.Exec(`DELETE FROM users WHERE true`)
}

func TestCreateUser(t *testing.T) {
	err := CreateUser(testDb, "test@localhost", "test", Admin)

	assert.NoError(t, err, "Error should be nil")

	var user User
	testDb.Last(&user)

	assert.Equal(t, "test@localhost", user.Email)
	assert.Equal(t, Admin, user.Role)
	assert.Nil(t, bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("test")))

	// cleanup after this test
	testDb.Delete(&user)
}

func TestGetUser(t *testing.T) {
	user := User{
		Email:    "test@localhost",
		Password: "test",
		Active:   true,
		Role:     Admin,
	}
	testDb.Create(&user)

	testDb.Last(&user)
	id := user.ID

	user2, err := GetUser(testDb, int(id))

	assert.Nil(t, err)

	assert.Equal(t, user.ID, user2.ID)
	assert.Equal(t, user.Email, user2.Email)
	assert.Equal(t, user.Role, user2.Role)
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
		assert.Equal(t, expectedUser.Role, user.Role)
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
