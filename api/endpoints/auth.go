package endpoints

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAuthEndpoint(db *gorm.DB) Endpoint {
	return Endpoint{
		Routes: []Route{
			{
				Method: "POST",
				Path:   "/login",
				Handler: func(c *gin.Context) {
					login(db, c)
				},
			},
		},
	}
}

func login(db *gorm.DB, c *gin.Context) {
	var userRequest models.UserRequest
	err := c.BindJSON(&userRequest)

	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	user := models.User{}

	tx := db.First(&user, "email = ?", userRequest.Email)

	if tx.RowsAffected == 0 {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	if !models.CheckPassword(userRequest.Password, user.Password) {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(200, gin.H{"token": "token"})
}
