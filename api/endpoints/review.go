package endpoints

import (
	"net/http"

	"github.com/adraismawur/mibig-submission/models/entry"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
			{
				Method: "POST",
				Path:   "/review/check",
				Handler: func(c *gin.Context) {
					checkReview(db, c)
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

type ReviewRequest struct {
	Accession string         `json:"accession"`
	Category  entry.Category `json:"category"`
}

func checkReview(db *gorm.DB, c *gin.Context) {
	var request ReviewRequest

	err := c.ShouldBindJSON(&request)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var activeReview entry.SubmissionReview
	err = db.
		Table("submission_reviews").
		Where("accession = $1 AND category = $2", request.Accession, request.Category).
		Find(&activeReview).
		Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response struct {
		State entry.ReviewState `json:"state"`
		User  uint64            `json:"user"`
	}

	response.State = activeReview.State
	response.User = activeReview.UserID

	c.JSON(http.StatusOK, response)
}
