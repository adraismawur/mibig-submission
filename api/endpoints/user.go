package endpoints

import "github.com/gin-gonic/gin"

func GetUserRoutes() []Route {
	return []Route{
		{
			Method:  "CREATE",
			Path:    "/user",
			Handler: createUser,
		},
		{
			Method:  "GET",
			Path:    "/user",
			Handler: getUser,
		},
		{
			Method:  "UPDATE",
			Path:    "/user",
			Handler: updateUser,
		},
		{
			Method:  "DELETE",
			Path:    "/user",
			Handler: deleteUser,
		},
	}
}

func createUser(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "create user",
	})
}

func getUser(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "get user",
	})
}

func updateUser(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "update user",
	})
}

func deleteUser(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "delete user",
	})
}
