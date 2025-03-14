package endpoints

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetUserEndpoint(db *gorm.DB) Endpoint {
	return Endpoint{
		Routes: []Route{
			{
				Method: "CREATE",
				Path:   "/user",
				Handler: func(c *gin.Context) {
					createUser(db, c)
				},
			},
			{
				Method: "GET",
				Path:   "/user",
				Handler: func(c *gin.Context) {
					getUser(db, c)
				},
			},
			{
				Method: "UPDATE",
				Path:   "/user",
				Handler: func(c *gin.Context) {
					updateUser(db, c)
				},
			},
			{
				Method: "DELETE",
				Path:   "/user",
				Handler: func(c *gin.Context) {
					deleteUser(db, c)
				},
			},
		},
	}
}

func createUser(db *gorm.DB, c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "create user",
	})
}

func getUser(db *gorm.DB, c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "get user",
	})
}

func updateUser(db *gorm.DB, c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "update user",
	})
}

func deleteUser(db *gorm.DB, c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "delete user",
	})
}
