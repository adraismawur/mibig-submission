package endpoints

import (
	"github.com/gin-gonic/gin"
)

func GetAuthRoutes() []Route {
	return []Route{
		{
			Method:  "POST",
			Path:    "/auth/login",
			Handler: login,
		},
	}
}

func login(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "login",
	})
}
