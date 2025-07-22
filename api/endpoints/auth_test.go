package endpoints

import (
	"github.com/adraismawur/mibig-submission/util"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestLogin(t *testing.T) {
	jsonParam := `{"email": "testadmin@localhost", "password": "test"}`
	c := util.CreateTestGinJsonRequest(jsonParam)

	login(testDb, c)

	assert.Equal(t, http.StatusOK, c.Writer.Status(), "Status code should be 200")
}
