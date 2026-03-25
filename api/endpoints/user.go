package endpoints

import (
	"errors"
	"github.com/adraismawur/mibig-submission/middleware"
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/util"
	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-password/password"
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
const DefaultPasswordLength = 8

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
				Method: http.MethodPost,
				Path:   UserPath + "/password/:id",
				Handler: func(c *gin.Context) {
					updateUserPassword(db, c)
				},
			},
			{
				Method: http.MethodPost,
				Path:   UserPath + "/password/reset",
				Handler: func(c *gin.Context) {
					passwordResetRequest(db, c)
				},
			},
			{
				Method: http.MethodPost,
				Path:   UserPath + "/password/challenge",
				Handler: func(c *gin.Context) {
					passwordResetChallenge(db, c)
				},
			},
			{
				Method: http.MethodPost,
				Path:   UserPath + "/register",
				Handler: func(c *gin.Context) {
					registerFirstTimeUser(db, c)
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
	if request.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	if request.Password == "" {
		pass, err := password.Generate(DefaultPasswordLength, DefaultPasswordLength/2, 0, true, false)

		if err != nil {
			slog.Error("[endpoints] [user] Failed to generate password", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		request.Password = pass
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

	search := c.Query("search")

	users, err := models.GetUsers(db, int(offset), int(limit), search)

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
	if token.User.ID != uint64(id) && !models.HasRole(token.User, models.Admin) {
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
		slog.Error("[endpoints] [oldUser] Error getting existing oldUser", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	err = models.UpdateUser(db, id, newUser)

	if err != nil {
		slog.Error("[user] Error when updating user", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated"})
}

func updateUserPassword(db *gorm.DB, c *gin.Context) {
	// get id from url
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	type updatePasswordRequest struct {
		NewPassword string `json:"new_password"`
	}

	var request *updatePasswordRequest

	err = c.ShouldBindJSON(&request)

	if err != nil {
		slog.Error("[endpoints] [user] Cannot bind request json")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := models.GetUser(db, id)

	if err != nil {
		slog.Error("[endpoints] [user] Could not find user", "id", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = models.UpdateUserPassword(db, int(user.ID), request.NewPassword)

	if err != nil {
		slog.Error("[endpoints] [user] Could not update user password")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update user password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated"})
}

func passwordResetRequest(db *gorm.DB, c *gin.Context) {
	var request struct {
		Email string `json:"email"`
	}

	err := c.ShouldBindJSON(&request)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not bind password reset request challenge json: " + err.Error()})
		return
	}

	err = db.Transaction(func(tx *gorm.DB) error {

		exists, transactionErr := models.GetUserExistsByEmail(tx, request.Email)

		// keep success vague for security reasons
		if transactionErr != nil {
			c.Status(http.StatusInternalServerError)
			return transactionErr
		}

		if !exists {
			c.Status(http.StatusUnauthorized)
			return nil
		}

		// from here on we can create the challenge
		randomString := util.RandomString(10)
		newChallenge := models.PasswordChallenge{
			Email:     request.Email,
			Challenge: randomString,
		}

		transactionErr = tx.
			Create(&newChallenge).
			Error

		if transactionErr != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate password challenge"})
			return transactionErr
		}

		c.JSON(http.StatusOK, newChallenge)

		return nil
	})

	if err != nil {
		slog.Error("[endpoints] [user] handling password reset request threw an error: " + err.Error())
	}

	return
}

func passwordResetChallenge(db *gorm.DB, c *gin.Context) {
	var request struct {
		models.PasswordChallenge
		NewPassword string `json:"new_password"`
	}

	err := c.ShouldBindJSON(&request)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not bind password request challenge json: " + err.Error()})
		return
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		var matchingChallenge models.PasswordChallenge

		transactionErr := tx.
			Model(&models.PasswordChallenge{}).
			Where("email = $1 AND challenge = $2", request.Email, request.Challenge).
			Find(&matchingChallenge).
			Error

		if matchingChallenge.ID == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Challenge failed"})
			return errors.New("challenge failed")
		}

		var userId int

		transactionErr = tx.
			Model(&models.User{}).
			Select("id").
			Where("email = $1", request.Email).
			Find(&userId).
			Error

		if transactionErr != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not get user associated with email"})
			return transactionErr
		}

		transactionErr = models.UpdateUserPassword(tx, userId, request.NewPassword)

		if transactionErr != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not update user password"})
			return transactionErr
		}

		return nil
	})

	if err != nil {
		return
	}

	c.Status(http.StatusOK)
}

func registerFirstTimeUser(db *gorm.DB, c *gin.Context) {
	var userInfo models.UserInfo

	err := db.Transaction(func(tx *gorm.DB) error {
		transactionErr := c.ShouldBindJSON(&userInfo)

		if transactionErr != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not bind user info json: " + transactionErr.Error()})
			return transactionErr
		}

		user, transactionErr := models.GetUserFromContext(c)

		if transactionErr != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "could not get user from request: " + transactionErr.Error()})
			return transactionErr
		}

		transactionErr = tx.
			Model(&user).
			Select("active").
			Where("id = $1", user.ID).
			Update("active", true).
			Error

		if transactionErr != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not update user: " + transactionErr.Error()})
			return transactionErr
		}

		userInfo.ID = user.Info.ID
		userInfo.UserID = user.Info.UserID

		transactionErr = tx.
			Model(&user).
			Association("Info").
			Replace(&userInfo)

		if transactionErr != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not update user info: " + transactionErr.Error()})
			return transactionErr
		}

		return nil
	})

	if err != nil {
		return
	}

	c.Status(http.StatusOK)
}

func deleteUser(db *gorm.DB, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "delete user",
	})
}

func deleteUserById(db *gorm.DB, c *gin.Context) {

}
