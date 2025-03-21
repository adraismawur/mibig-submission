package endpoints

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetReviewEndpoint returns the review endpoint. This endpoint will implement adding, updating, and deleting reviews.
func GetReviewEndpoint(db *gorm.DB) Endpoint {
	return Endpoint{
		Routes: []Route{
			{
				Method: "GET",
				Path:   "/review",
				Handler: func(c *gin.Context) {
					getReview(db, c)
				},
			},
			{
				Method: "POST",
				Path:   "/review",
				Handler: func(c *gin.Context) {
					postReview(db, c)
				},
			},
		},
	}
}

func getReview(db *gorm.DB, c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "get review",
	})
}

func postReview(db *gorm.DB, c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "post review",
	})
}
