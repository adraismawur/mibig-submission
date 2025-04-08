package endpoints

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func init() {
	RegisterEndpointGenerator(ReviewEndpoint)
}

// ReviewEndpoint returns the review endpoint. This endpoint will implement adding, updating, and deleting reviews.
func ReviewEndpoint(db *gorm.DB) Endpoint {
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
	c.JSON(http.StatusOK, gin.H{
		"message": "get review",
	})
}

func postReview(db *gorm.DB, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "post review",
	})
}
