package endpoints

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/util"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestLogin(t *testing.T) {
	gormDB, mock := util.CreateMockDB()

	models.Migrate(gormDB)

	expectedRows := mock.NewRows([]string{"email", "password"}).
		AddRow("test@localhost", "$2a$10$wOLM7A7gHgQXKKnyZX2J.uWi41KZKd.vfzKqa.w.9hUVFGVk.4LB.")

	mock.ExpectQuery(`SELECT(.*)`).
		WillReturnRows(expectedRows)

	jsonParam := `{"email": "test@localhost", "password": "test"}`
	c := util.CreateMockGinJsonRequest(jsonParam)

	login(gormDB, c)

	assert.Equal(t, http.StatusOK, c.Writer.Status(), "Status code should be 200")
}
