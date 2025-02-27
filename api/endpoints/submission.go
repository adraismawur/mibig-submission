package endpoints

import "github.com/gin-gonic/gin"

func GetSubmissionRoutes() []Route {
	return []Route{
		{
			Method:  "CREATE",
			Path:    "/submission",
			Handler: createSubmission,
		},
		{
			Method:  "GET",
			Path:    "/submission",
			Handler: getSubmission,
		},
		{
			Method:  "UPDATE",
			Path:    "/submission",
			Handler: updateSubmission,
		},
		{
			Method:  "DELETE",
			Path:    "/submission",
			Handler: deleteSubmission,
		},
	}
}

func createSubmission(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "create submission",
	})
}

func getSubmission(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "get submission",
	})
}

func updateSubmission(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "update submission",
	})
}

func deleteSubmission(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "delete submission",
	})
}
