package endpoints

import (
	"github.com/adraismawur/mibig-submission/middleware"
	"github.com/adraismawur/mibig-submission/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

const UserPath = "/user"

func init() {
	RegisterEndpointGenerator(UserEndpoint)
	middleware.AddProtectedRoute(http.MethodPut, UserPath, models.Admin)
	middleware.AddProtectedRoute(http.MethodPatch, UserPath, models.Admin)
	middleware.AddProtectedRoute(http.MethodDelete, UserPath, models.Admin)
}

const DEFAULT_GET_LIMIT = 10

// UserEndpoint returns the user endpoint. This endpoint will implement creating, updating, and deleting users.
func UserEndpoint(db *gorm.DB) Endpoint {
	return Endpoint{
		Routes: []Route{
			{
				Method: http.MethodPut,
				Path:   UserPath,
				Handler: func(c *gin.Context) {
					createUser(db, c)
				},
			},
			{
				Method: http.MethodGet,
				Path:   UserPath,
				Handler: func(c *gin.Context) {
					getUsers(db, c)
				},
			},
			{
				Method: http.MethodGet,
				Path:   UserPath + "/:id",
				Handler: func(c *gin.Context) {
					getUserWithId(db, c)
				},
			},
			{
				Method: http.MethodPatch,
				Path:   UserPath,
				Handler: func(c *gin.Context) {
					updateUser(db, c)
				},
			},
			{
				Method: http.MethodDelete,
				Path:   UserPath,
				Handler: func(c *gin.Context) {
					deleteUser(db, c)
				},
			},
		},
	}
}

func createUser(db *gorm.DB, c *gin.Context) {
	// bind json
	var request models.User
	if err := c.ShouldBindJSON(&request); err != nil {
		panic(err.Error())
		return
	}

	// validate user
	// needs to have email and password
	if request.Email == "" || request.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}

	// validate role
	if request.Role < models.Submitter || request.Role > models.Admin {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}

	// check if user exists
	exists, err := models.GetUserExistsByEmail(db, request.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with this email already exists"})
		return
	}

	err = models.CreateUser(db, request.Email, request.Password, request.Role)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created"})
}

func getUsers(db *gorm.DB, c *gin.Context) {
	// get optional query parameters
	limit, err := strconv.ParseInt(c.Query("limit"), 10, 32)
	if err != nil {
		limit = DEFAULT_GET_LIMIT
	}

	offset, err := strconv.ParseInt(c.Query("offset"), 10, 32)
	if err != nil {
		offset = 0
	}

	users, err := models.GetUsers(db, int(offset), int(limit))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func getUserWithId(db *gorm.DB, c *gin.Context) {
	// get id from url
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	user, err := models.GetUser(db, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func updateUser(db *gorm.DB, c *gin.Context) {
	// get id from url
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	// check if user exists
	user, err := models.GetUser(db, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User does not exist"})
		return
	}

	// bind request
	err = c.BindJSON(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = models.UpdateUser(db, id, user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated"})
}

func deleteUser(db *gorm.DB, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "delete user",
	})
}

func deleteUserById(db *gorm.DB, c *gin.Context) {

}
