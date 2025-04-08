package models

import (
	"github.com/adraismawur/mibig-submission/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckPassword(t *testing.T) {
	hashedTestPassword := "$2a$10$wOLM7A7gHgQXKKnyZX2J.uWi41KZKd.vfzKqa.w.9hUVFGVk.4LB."
	plainTestPassword := "test"

	assert.True(t, CheckPassword(plainTestPassword, hashedTestPassword))
}

func TestHasRole(t *testing.T) {
	db, mock := util.CreateMockDB()

	expectedRows := mock.NewRows([]string{"role"}).
		AddRow(Admin)

	mock.ExpectQuery(`SELECT(.*)`).
		WillReturnRows(expectedRows)

	user := User{
		Email:    "test@localhost",
		Role:     Admin,
		Password: "test123",
	}

	assert.True(t, HasRole(db, user, Admin))
}
