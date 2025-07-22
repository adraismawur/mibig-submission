package endpoints

import (
	"github.com/adraismawur/mibig-submission/middleware"
	"github.com/adraismawur/mibig-submission/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
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

const DefaultGetLimit = 10

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
				Path:   UserPath + "/:id",
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("[endpoints] [user] Failed to bind request", "error", err)
		return
	}

	// validate user
	// needs to have email and password
	if request.Email == "" || request.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}

	// validate roles
	for _, userRole := range request.Roles {
		if userRole.Role < models.Submitter || userRole.Role > models.Admin {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
			return
		}
	}

	// ensure a minimum role is set
	if len(request.Roles) == 0 {
		request.Roles = append(request.Roles, models.UserRole{
			Role: models.Submitter,
		})
	}

	// check if user exists
	exists, err := models.GetUserExistsByEmail(db, request.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("[endpoints] [user] Error when checking for existing user", "error", err)
		return
	}

	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with this email already exists"})
		return
	}

	err = models.CreateUser(db, request.Email, request.Password, request.Roles)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("[endpoints] [user] Error when creating user", "error", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created"})
}

func getUsers(db *gorm.DB, c *gin.Context) {
	// check if requesting user is an admin
	token := models.Token{}

	valid := middleware.ValidateAuthHeader(c, &token)

	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	if !models.HasRole(token.User, models.Admin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		c.Abort()
		return
	}

	// get optional query parameters
	limit, err := strconv.ParseInt(c.Query("limit"), 10, 32)
	if err != nil {
		limit = DefaultGetLimit
	}

	offset, err := strconv.ParseInt(c.Query("offset"), 10, 32)
	if err != nil {
		offset = 0
	}

	users, err := models.GetUsers(db, int(offset), int(limit))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("[endpoints] [user] Error when retrieving all users", "error", err)
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

	// check if user is admin or the user itself
	token := models.Token{}
	valid := middleware.ValidateAuthHeader(c, &token)
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	// if user is not the requested user and is not an admin, return forbidden
	if token.User.ID != uint(id) && !models.HasRole(token.User, models.Admin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		c.Abort()
		return
	}

	user, err := models.GetUser(db, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("[endpoints] [user] Error when getting user by Role", "error", err)
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

	// check if oldUser exists
	oldUser, err := models.GetUser(db, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("[endpoints] [oldUser] Error getting existing oldUser", "error", err)
		return
	}

	if oldUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User does not exist"})
		return
	}

	var newUser models.User

	// bind request
	err = c.ShouldBindJSON(&newUser)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = models.UpdateUser(db, id, oldUser, &newUser)

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
