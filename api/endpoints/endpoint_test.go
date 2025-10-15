package endpoints

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/util/test_utils"
	"gorm.io/gorm"
	"testing"
)

var testDb *gorm.DB

const numTestUsers = 10

var adminTestUser = models.User{
	Email:    "testadmin@localhost",
	Password: test_utils.TestPassword,
	Active:   true,
	Roles: []models.UserRole{
		{
			Role: models.Admin,
		},
	},
	Info: models.UserInfo{},
}
var submitterTestUser = models.User{
	Email:    "testsub@localhost",
	Password: "$2a$10$wOLM7A7gHgQXKKnyZX2J.uWi41KZKd.vfzKqa.w.9hUVFGVk.4LB.",
	Active:   true,
	Roles: []models.UserRole{
		{
			Role: models.Submitter,
		},
	},
	Info: models.UserInfo{},
}

var randomTestUsers []models.User

func TestMain(m *testing.M) {
	// setup
	testDb = test_utils.CreateTestDB()
	models.Migrate(testDb)

	// fixed test users
	testDb.Create(&adminTestUser)
	testDb.Create(&submitterTestUser)

	// create some random test users
	for i := 0; i < numTestUsers; i++ {
		user := models.User{
			Email:    models.GenerateRandomEmail(),
			Password: "test",
			Active:   true,
			Roles: []models.UserRole{
				{
					Role: models.Admin,
				},
			},
		}
		testDb.Create(&user)
		randomTestUsers = append(randomTestUsers, user)
	}

	m.Run()

	// teardown
	testDb.Exec(`DELETE FROM users WHERE true`)
}
