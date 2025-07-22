package endpoints

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/util"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCreateUserBadData(t *testing.T) {
	c := util.CreateTestGinJsonRequest("{}")

	c.Request.Method = "POST"
	c.Request.URL.Path = "/user"

	createUser(testDb, c)

	assert.Equal(t, http.StatusBadRequest, c.Writer.Status(), "Status code should be 400")
}

func TestCreateUserNoEmail(t *testing.T) {
	c, r := util.CreateTestGinJsonRequestWithRecorder("{\"email\": \"\", \"password\": \"test\", \"role\": 2}")

	c.Request.Method = "POST"
	c.Request.URL.Path = "/user"

	createUser(testDb, c)

	assert.Equal(t, http.StatusBadRequest, c.Writer.Status(), "Status code should be 400")

	response := r.Body.String()
	assert.Contains(t, response, "Email and password are required", "Response should contain 'Email and password are required'")
}

func TestCreateUserNoPassword(t *testing.T) {
	c, r := util.CreateTestGinJsonRequestWithRecorder("{\"email\": \"test@localhost\", \"password\": \"\", \"role\": 2}")

	c.Request.Method = "POST"
	c.Request.URL.Path = "/user"

	createUser(testDb, c)

	assert.Equal(t, http.StatusBadRequest, c.Writer.Status(), "Status code should be 400")

	response := r.Body.String()
	assert.Contains(t, response, "Email and password are required", "Response should contain 'Email and password are required'")
}

func TestCreateUserInvalidRole(t *testing.T) {
	testUserEmail := "test@localhost"
	testUserPassword := "test"
	c, r := util.CreateTestGinJsonRequestWithRecorder("{\"email\": \"" + testUserEmail + "\", \"password\": \"" + testUserPassword + "\", \"roles\": [{\"role\": 11037}]}")

	c.Request.Method = "POST"
	c.Request.URL.Path = "/user"

	createUser(testDb, c)

	assert.Equal(t, http.StatusBadRequest, c.Writer.Status(), "Status code should be 400")

	response := r.Body.String()
	assert.Contains(t, response, "Invalid role", "Response should contain 'Invalid role'")
}

func TestCreateUserAlreadyExists(t *testing.T) {
	c, r := util.CreateTestGinJsonRequestWithRecorder("{\"email\": \"testadmin@localhost\", \"password\": \"test\", \"roles\": [{\"role\": 2}]}")

	c.Request.Method = "POST"
	c.Request.URL.Path = "/user"

	createUser(testDb, c)

	assert.Equal(t, http.StatusBadRequest, c.Writer.Status(), "Status code should be 400")

	response := r.Body.String()
	assert.Contains(t, response, "User with this email already exists", "Response should contain 'User already exists'")
}

func TestCreateUser(t *testing.T) {
	testEmail := util.GenerateRandomEmail()
	c, r := util.CreateTestGinJsonRequestWithRecorder("{\"email\": \"" + testEmail + "\", \"password\": \"test\", \"role\": 2}")

	c.Request.Method = "POST"
	c.Request.URL.Path = "/user"

	createUser(testDb, c)

	assert.Equal(t, http.StatusOK, c.Writer.Status(), "Status code should be 200")
	assert.Contains(t, r.Body.String(), "User created", "Response should contain 'User created'")
}

func TestGetUsers(t *testing.T) {
	c, r := util.CreateTestGinGetRequest("/user")

	user := models.User{
		ID:    1,
		Email: util.GenerateRandomEmail(),
		Roles: []models.UserRole{
			{
				Role: models.Admin,
			},
		},
	}

	testToken, _ := models.GenerateToken(user)

	util.AddTokenToHeader(c, testToken)

	getUsers(testDb, c)

	assert.Equal(t, http.StatusOK, c.Writer.Status(), "Status code should be 200")

	var users []models.User
	err := json.Unmarshal(r.Body.Bytes(), &users)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}
}

func TestGetUsersForbidden(t *testing.T) {
	c, _ := util.CreateTestGinGetRequest("/user")

	user := models.User{
		ID:    1,
		Email: util.GenerateRandomEmail(),
		Roles: []models.UserRole{
			{
				Role: models.Submitter,
			},
		},
	}

	testToken, _ := models.GenerateToken(user)

	util.AddTokenToHeader(c, testToken)

	getUsers(testDb, c)

	assert.Equal(t, http.StatusForbidden, c.Writer.Status(), "Status code should be 403")
}

func TestGetUserWithIdAdmin(t *testing.T) {
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
		Roles: []models.UserRole{
			{
				Role: models.Admin,
			},
		},
	}

	testToken, _ := models.GenerateToken(user)

	util.AddTokenToHeader(c, testToken)

	getUserWithId(testDb, c)

	assert.Equal(t, http.StatusOK, c.Writer.Status(), "Status code should be 200")

	var responseUser models.User
	err := json.Unmarshal(r.Body.Bytes(), &responseUser)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}
	assert.Equal(t, testEmail, user.Email, "Email should be 'test@localhost'")
}

func TestGetUserWithIdSelf(t *testing.T) {
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
		Roles: []models.UserRole{
			{
				Role: models.Submitter,
			},
		},
	}

	testToken, _ := models.GenerateToken(user)

	util.AddTokenToHeader(c, testToken)

	getUserWithId(testDb, c)

	assert.Equal(t, http.StatusOK, c.Writer.Status(), "Status code should be 200")

	var responseUser models.User
	err := json.Unmarshal(r.Body.Bytes(), &responseUser)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}
	assert.Equal(t, testEmail, user.Email, "Email should be 'test@localhost'")
}

func TestGetUserWithIdForbidden(t *testing.T) {
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
		Roles: []models.UserRole{
			{
				Role: models.Submitter,
			},
		},
	}

	testToken, _ := models.GenerateToken(user)

	util.AddTokenToHeader(c, testToken)

	getUserWithId(testDb, c)

	assert.Equal(t, http.StatusForbidden, c.Writer.Status(), "Status code should be 403")
}

func TestGetUserWithIdNotFound(t *testing.T) {
	// we're generating around 10 users in these tests. if this ever fails because of an existing user something is
	// seriously wrong
	c, r := util.CreateTestGinGetRequest("/user/100000000")
	c.Params = []gin.Param{
		{
			Key:   "id",
			Value: "100000000",
		},
	}

	user := models.User{
		ID:    1,
		Email: util.GenerateRandomEmail(),
		Roles: []models.UserRole{
			{
				Role: models.Admin,
			},
		},
	}

	testToken, _ := models.GenerateToken(user)

	util.AddTokenToHeader(c, testToken)

	getUserWithId(testDb, c)

	assert.Equal(t, http.StatusNotFound, c.Writer.Status(), "Status code should be 404")

	response := r.Body.String()
	assert.Contains(t, response, "User not found", "Response should contain 'User not found'")
}

func TestUpdateUserDoesNotExist(t *testing.T) {
	c, r := util.CreateTestGinJsonRequestWithRecorder("{\"email\": \"test@localhost\", \"role\": 0}")

	c.Request.Method = "PUT"
	c.Request.URL.Path = "/user/100000000"

	c.Params = []gin.Param{
		{
			Key:   "id",
			Value: "100000000",
		},
	}

	updateUser(testDb, c)

	// response 404
	assert.Equal(t, http.StatusNotFound, c.Writer.Status(), "Status code should be 404")
	assert.Contains(t, r.Body.String(), "User does not exist", "Response should contain 'User does not exist'")
}
