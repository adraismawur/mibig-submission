package endpoints

import "github.com/gin-gonic/gin"

func GetReviewRoutes() []Route {
	return []Route{
		{
			Method:  "GET",
			Path:    "/review",
			Handler: getReview,
		},
		{
			Method:  "POST",
			Path:    "/review",
			Handler: postReview,
		},
	}
}

func getReview(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "get review",
	})
}

func postReview(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "post review",
	})
}
