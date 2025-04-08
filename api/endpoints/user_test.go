package endpoints

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/util"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"regexp"
	"testing"
)

func TestCreateUserBadData(t *testing.T) {
	db, _ := util.CreateMockDB()

	c := util.CreateMockGinJsonRequest("{}")

	c.Request.Method = "POST"
	c.Request.URL.Path = "/user"

	createUser(db, c)

	assert.Equal(t, http.StatusBadRequest, c.Writer.Status(), "Status code should be 400")
}

func TestCreateUserNoEmail(t *testing.T) {
	db, _ := util.CreateMockDB()

	c, r := util.CreateMockGinJsonRequestWithRecorder("{\"email\": \"\", \"password\": \"test\", \"role\": 2}")

	c.Request.Method = "POST"
	c.Request.URL.Path = "/user"

	createUser(db, c)

	assert.Equal(t, http.StatusBadRequest, c.Writer.Status(), "Status code should be 400")

	response := r.Body.String()
	assert.Contains(t, response, "Email and password are required", "Response should contain 'Email and password are required'")
}

func TestCreateUserNoPassword(t *testing.T) {
	db, _ := util.CreateMockDB()

	c, r := util.CreateMockGinJsonRequestWithRecorder("{\"email\": \"test@localhost\", \"password\": \"\", \"role\": 2}")

	c.Request.Method = "POST"
	c.Request.URL.Path = "/user"

	createUser(db, c)

	assert.Equal(t, http.StatusBadRequest, c.Writer.Status(), "Status code should be 400")

	response := r.Body.String()
	assert.Contains(t, response, "Email and password are required", "Response should contain 'Email and password are required'")
}

func TestCreateUserInvalidRole(t *testing.T) {
	db, _ := util.CreateMockDB()

	c, r := util.CreateMockGinJsonRequestWithRecorder("{\"email\": \"test@localhost\", \"password\": \"test\", \"role\": -1}")

	c.Request.Method = "POST"
	c.Request.URL.Path = "/user"

	createUser(db, c)

	assert.Equal(t, http.StatusBadRequest, c.Writer.Status(), "Status code should be 400")

	response := r.Body.String()
	assert.Contains(t, response, "Invalid role", "Response should contain 'Invalid role'")
}

func TestCreateUserAlreadyExists(t *testing.T) {
	db, mock := util.CreateMockDB()

	expectedRows := mock.NewRows([]string{"email", "password", "role"}).
		AddRow("test@localhost", "test", 2)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
		WillReturnRows(expectedRows)

	c, r := util.CreateMockGinJsonRequestWithRecorder("{\"email\": \"test@localhost\", \"password\": \"test\", \"role\": 2}")

	c.Request.Method = "POST"
	c.Request.URL.Path = "/user"

	createUser(db, c)

	assert.Equal(t, http.StatusBadRequest, c.Writer.Status(), "Status code should be 400")

	response := r.Body.String()
	assert.Contains(t, response, "User with this email already exists", "Response should contain 'User already exists'")
}

func TestCreateUser(t *testing.T) {
	db, _ := util.CreateMockDB()

	c := util.CreateMockGinJsonRequest("{\"email\": \"test@localhost\", \"password\": \"test\", \"role\": 2}")

	c.Request.Method = "POST"
	c.Request.URL.Path = "/user"

	createUser(db, c)

	assert.Equal(t, http.StatusOK, c.Writer.Status(), "Status code should be 200")
}

func TestGetUsers(t *testing.T) {
	db, mock := util.CreateMockDB()

	expectedRows := mock.NewRows([]string{"email", "password", "role"})

	testRowCount := 5

	for i := 0; i < testRowCount; i++ {
		randomEmail := util.GenerateRandomEmail()
		expectedRows.AddRow(randomEmail, "test", 2)
	}

	mock.ExpectQuery(`SELECT \* FROM "users"`).
		WillReturnRows(expectedRows)

	c, r := util.CreateMockGinGetRequest("/user")

	getUsers(db, c)

	assert.Equal(t, http.StatusOK, c.Writer.Status(), "Status code should be 200")

	var users []models.User
	err := json.Unmarshal(r.Body.Bytes(), &users)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}
	assert.Equal(t, testRowCount, len(users), "Number of users should be %d", testRowCount)
}

func TestGetUserWithId(t *testing.T) {
	db, mock := util.CreateMockDB()

	expectedRows := mock.NewRows([]string{"email", "password", "role"}).
		AddRow("test@localhost", "test", 2)

	mock.ExpectQuery(`SELECT \* FROM "users"`).
		WithArgs(1, 1).
		WillReturnRows(expectedRows)

	c, r := util.CreateMockGinGetRequest("/user/1")
	c.Params = []gin.Param{
		{
			Key:   "id",
			Value: "1",
		},
	}

	getUserWithId(db, c)

	assert.Equal(t, http.StatusOK, c.Writer.Status(), "Status code should be 200")

	var user models.User
	err := json.Unmarshal(r.Body.Bytes(), &user)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}
	assert.Equal(t, "test@localhost", user.Email, "Email should be 'test@localhost'")
}

func TestGetUserWithIdNotFound(t *testing.T) {
	db, mock := util.CreateMockDB()

	// should return no rows
	mock.ExpectQuery(`SELECT \* FROM "users"`).
		WithArgs(2, 1).
		WillReturnRows(mock.NewRows([]string{"email", "password", "role"}))

	c, r := util.CreateMockGinGetRequest("/user/2")
	c.Params = []gin.Param{
		{
			Key:   "id",
			Value: "2",
		},
	}

	getUserWithId(db, c)

	assert.Equal(t, http.StatusNotFound, c.Writer.Status(), "Status code should be 404")

	response := r.Body.String()
	assert.Contains(t, response, "User not found", "Response should contain 'User not found'")
}
