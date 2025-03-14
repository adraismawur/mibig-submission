package endpoints

import (
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
	c.JSON(200, gin.H{
		"message": "login",
	})
}
