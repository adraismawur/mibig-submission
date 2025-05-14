package endpoints

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/util"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"testing"
	"time"
)

var testDb *gorm.DB

const numTestUsers = 10

var testUsers []models.User

func TestMain(m *testing.M) {
	// setup
	testDb = util.CreateTestDB()
	models.Migrate(testDb)

	// create some test users

	for i := 0; i < numTestUsers; i++ {
		user := models.User{
			Email:    util.GenerateRandomEmail(),
			Password: "test",
			Active:   true,
			Role:     models.Admin,
		}
		testDb.Create(&user)
		testUsers = append(testUsers, user)
	}

	m.Run()

	// teardown
	testDb.Exec(`DELETE FROM users WHERE true`)
}

func TestCreateUserBadData(t *testing.T) {
	db, _ := util.CreateMockDB()

	c := util.CreateTestGinJsonRequest("{}")

	c.Request.Method = "POST"
	c.Request.URL.Path = "/user"

	createUser(db, c)

	assert.Equal(t, http.StatusBadRequest, c.Writer.Status(), "Status code should be 400")
}

func TestCreateUserNoEmail(t *testing.T) {
	db, _ := util.CreateMockDB()

	c, r := util.CreateTestGinJsonRequestWithRecorder("{\"email\": \"\", \"password\": \"test\", \"role\": 2}")

	c.Request.Method = "POST"
	c.Request.URL.Path = "/user"

	createUser(db, c)

	assert.Equal(t, http.StatusBadRequest, c.Writer.Status(), "Status code should be 400")

	response := r.Body.String()
	assert.Contains(t, response, "Email and password are required", "Response should contain 'Email and password are required'")
}

func TestCreateUserNoPassword(t *testing.T) {
	db, _ := util.CreateMockDB()

	c, r := util.CreateTestGinJsonRequestWithRecorder("{\"email\": \"test@localhost\", \"password\": \"\", \"role\": 2}")

	c.Request.Method = "POST"
	c.Request.URL.Path = "/user"

	createUser(db, c)

	assert.Equal(t, http.StatusBadRequest, c.Writer.Status(), "Status code should be 400")

	response := r.Body.String()
	assert.Contains(t, response, "Email and password are required", "Response should contain 'Email and password are required'")
}

func TestCreateUserInvalidRole(t *testing.T) {
	db, _ := util.CreateMockDB()

	c, r := util.CreateTestGinJsonRequestWithRecorder("{\"email\": \"test@localhost\", \"password\": \"test\", \"role\": -1}")

	c.Request.Method = "POST"
	c.Request.URL.Path = "/user"

	createUser(db, c)

	assert.Equal(t, http.StatusBadRequest, c.Writer.Status(), "Status code should be 400")

	response := r.Body.String()
	assert.Contains(t, response, "Invalid role", "Response should contain 'Invalid role'")
}

func TestCreateUserAlreadyExists(t *testing.T) {
	db, mock := util.CreateMockDB()

	expectedResult := sqlmock.NewResult(1, 1)

	mock.ExpectExec("SELECT.*").
		WithArgs("test@localhost").
		WillReturnResult(expectedResult)

	c, r := util.CreateTestGinJsonRequestWithRecorder("{\"email\": \"test@localhost\", \"password\": \"test\", \"role\": 2}")

	c.Request.Method = "POST"
	c.Request.URL.Path = "/user"

	createUser(db, c)

	assert.Equal(t, http.StatusBadRequest, c.Writer.Status(), "Status code should be 400")

	response := r.Body.String()
	assert.Contains(t, response, "User with this email already exists", "Response should contain 'User already exists'")
}

func TestCreateUser(t *testing.T) {
	db, mock := util.CreateMockDB()

	// creating first checks if the user exists. this should return false
	expectedResult := sqlmock.NewResult(1, 0)

	mock.ExpectExec("SELECT.*").
		WithArgs("test@localhost").
		WillReturnResult(expectedResult)

	expectedRows := mock.NewRows([]string{"email", "role"}).
		AddRow("test@localhost", 2)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users"`).
		WithArgs("test@localhost", util.AnyString{}, true, 2).
		WillReturnRows(expectedRows)
	mock.ExpectCommit()

	c, r := util.CreateTestGinJsonRequestWithRecorder("{\"email\": \"test@localhost\", \"password\": \"test\", \"role\": 2}")

	c.Request.Method = "POST"
	c.Request.URL.Path = "/user"

	createUser(db, c)

	assert.Equal(t, http.StatusOK, c.Writer.Status(), "Status code should be 200")
	assert.Contains(t, r.Body.String(), "User created", "Response should contain 'User created'")
}

func TestGetUsers(t *testing.T) {
	db, mock := util.CreateMockDB()

	expectedRows := mock.NewRows([]string{"email", "password", "role"})

	testRowCount := 5

	for i := 0; i < testRowCount; i++ {
		randomEmail := util.GenerateRandomEmail()
		expectedRows.AddRow(randomEmail, "test", 2)
	}

	mock.ExpectQuery(`SELECT .* FROM "users"`).
		WillReturnRows(expectedRows)

	c, r := util.CreateTestGinGetRequest("/user")

	user := models.User{
		ID:    1,
		Email: util.GenerateRandomEmail(),
		Role:  models.Admin,
	}

	testToken, _ := models.GenerateToken(user)

	util.AddTokenToHeader(c, testToken)

	getUsers(db, c)

	assert.Equal(t, http.StatusOK, c.Writer.Status(), "Status code should be 200")

	var users []models.User
	err := json.Unmarshal(r.Body.Bytes(), &users)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}
	assert.Equal(t, testRowCount, len(users), "Number of users should be %d", testRowCount)
}

func TestGetUsersForbidden(t *testing.T) {
	db, mock := util.CreateMockDB()

	expectedRows := mock.NewRows([]string{"email", "password", "role"})

	testRowCount := 5

	for i := 0; i < testRowCount; i++ {
		randomEmail := util.GenerateRandomEmail()
		expectedRows.AddRow(randomEmail, "test", 2)
	}

	mock.ExpectQuery(`SELECT .* FROM "users"`).
		WillReturnRows(expectedRows)

	c, _ := util.CreateTestGinGetRequest("/user")

	user := models.User{
		ID:    1,
		Email: util.GenerateRandomEmail(),
		Role:  models.Submitter,
	}

	testToken, _ := models.GenerateToken(user)

	util.AddTokenToHeader(c, testToken)

	getUsers(db, c)

	assert.Equal(t, http.StatusForbidden, c.Writer.Status(), "Status code should be 403")
}

func TestGetUserWithIdAdmin(t *testing.T) {
	db, mock := util.CreateMockDB()

	expectedRows := mock.NewRows([]string{"email", "password", "role"}).
		AddRow("test@localhost", "test", 2)

	mock.ExpectQuery(`SELECT .* FROM "users"`).
		WithArgs(1).
		WillReturnRows(expectedRows)

	c, r := util.CreateTestGinGetRequest("/user/1")
	c.Params = []gin.Param{
		{
			Key:   "id",
			Value: "1",
		},
	}

	testEmail := util.GenerateRandomEmail()

	user := models.User{
		ID:    1,
		Email: testEmail,
		Role:  models.Admin,
	}

	testToken, _ := models.GenerateToken(user)

	util.AddTokenToHeader(c, testToken)

	getUserWithId(db, c)

	assert.Equal(t, http.StatusOK, c.Writer.Status(), "Status code should be 200")

	var responseUser models.User
	err := json.Unmarshal(r.Body.Bytes(), &responseUser)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}
	assert.Equal(t, testEmail, user.Email, "Email should be 'test@localhost'")
}

func TestGetUserWithIdSelf(t *testing.T) {
	db, mock := util.CreateMockDB()

	expectedRows := mock.NewRows([]string{"testEmail", "password", "role"}).
		AddRow("test@localhost", "test", 2)

	mock.ExpectQuery(`SELECT .* FROM "users"`).
		WithArgs(1).
		WillReturnRows(expectedRows)

	c, r := util.CreateTestGinGetRequest("/user/1")
	c.Params = []gin.Param{
		{
			Key:   "id",
			Value: "1",
		},
	}

	testEmail := util.GenerateRandomEmail()

	user := models.User{
		ID:    1,
		Email: testEmail,
		Role:  models.Submitter,
	}

	testToken, _ := models.GenerateToken(user)

	util.AddTokenToHeader(c, testToken)

	getUserWithId(db, c)

	assert.Equal(t, http.StatusOK, c.Writer.Status(), "Status code should be 200")

	var responseUser models.User
	err := json.Unmarshal(r.Body.Bytes(), &responseUser)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}
	assert.Equal(t, testEmail, user.Email, "Email should be 'test@localhost'")
}

func TestGetUserWithIdForbidden(t *testing.T) {
	db, mock := util.CreateMockDB()

	expectedRows := mock.NewRows([]string{"email", "password", "role"}).
		AddRow("test@localhost", "test", 2)

	mock.ExpectQuery(`SELECT .* FROM "users"`).
		WithArgs(1).
		WillReturnRows(expectedRows)

	c, _ := util.CreateTestGinGetRequest("/user/1")
	c.Params = []gin.Param{
		{
			Key:   "id",
			Value: "1",
		},
	}

	// token has user id 2, should return 403
	user := models.User{
		ID:    2,
		Email: util.GenerateRandomEmail(),
		Role:  models.Submitter,
	}

	testToken, _ := models.GenerateToken(user)

	util.AddTokenToHeader(c, testToken)

	getUserWithId(db, c)

	assert.Equal(t, http.StatusForbidden, c.Writer.Status(), "Status code should be 403")
}

func TestGetUserWithIdNotFound(t *testing.T) {
	db, mock := util.CreateMockDB()

	// should return no rows
	mock.ExpectQuery(`SELECT .* FROM "users"`).
		WithArgs(2).
		WillReturnRows(mock.NewRows([]string{"email", "password", "role"}))

	c, r := util.CreateTestGinGetRequest("/user/2")
	c.Params = []gin.Param{
		{
			Key:   "id",
			Value: "2",
		},
	}

	user := models.User{
		ID:    1,
		Email: util.GenerateRandomEmail(),
		Role:  models.Admin,
	}

	testToken, _ := models.GenerateToken(user)

	util.AddTokenToHeader(c, testToken)

	getUserWithId(db, c)

	assert.Equal(t, http.StatusNotFound, c.Writer.Status(), "Status code should be 404")

	response := r.Body.String()
	assert.Contains(t, response, "User not found", "Response should contain 'User not found'")
}

func TestUpdateUser(t *testing.T) {
	db, mock := util.CreateMockDB()

	// creating first gets the existing user. existing user is an admin
	expectedRows := sqlmock.NewRows([]string{"createdat", "updatedat", "id", "email", "password", "active", "role"})
	expectedRows.AddRow(time.Now(), time.Now(), 1, "test@localhost", "test", true, models.Admin)

	mock.ExpectQuery(`SELECT .* FROM "users"`).WithArgs(1).WillReturnRows(expectedRows)

	// updating the user. this sets the user to be a submitter
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"users\" SET").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// setting the user to be submitter
	c, r := util.CreateTestGinJsonRequestWithRecorder("{\"email\": \"test@localhost\", \"role\": 0}")

	c.Request.Method = "PUT"
	c.Request.URL.Path = "/user/1"

	c.Params = []gin.Param{
		{
			Key:   "id",
			Value: "1",
		},
	}

	updateUser(db, c)

	// response 200
	assert.Equal(t, http.StatusOK, c.Writer.Status(), "Status code should be 200. Response")
	assert.Contains(t, r.Body.String(), "User updated", "Response should contain 'User updated'")
}

func TestUpdateUserDoesNotExist(t *testing.T) {
	db, mock := util.CreateMockDB()

	mock.ExpectQuery(`SELECT .* FROM "users"`).WithArgs(1).
		WillReturnRows(mock.NewRows([]string{"email", "password", "role"}))

	c, r := util.CreateTestGinJsonRequestWithRecorder("{\"email\": \"test@localhost\", \"role\": 0}")

	c.Request.Method = "PUT"
	c.Request.URL.Path = "/user/1"

	c.Params = []gin.Param{
		{
			Key:   "id",
			Value: "1",
		},
	}

	updateUser(db, c)

	// response 404
	assert.Equal(t, http.StatusNotFound, c.Writer.Status(), "Status code should be 404")
	assert.Contains(t, r.Body.String(), "User does not exist", "Response should contain 'User does not exist'")
}
