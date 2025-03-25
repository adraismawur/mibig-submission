package endpoints

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func init() {
	RegisterEndpointGenerator(GenerateSubmissionEndpoint)
}

// GenerateSubmissionEndpoint returns the submission endpoint.
// This endpoint will implement creating and updating submissions, as well as perform some
// specific checks on submissions.
func GenerateSubmissionEndpoint(db *gorm.DB) Endpoint {
	return Endpoint{
		Routes: []Route{
			{
				Method: "CREATE",
				Path:   "/submission",
				Handler: func(c *gin.Context) {
					createSubmission(db, c)
				},
			},
			{
				Method: "GET",
				Path:   "/submission",
				Handler: func(c *gin.Context) {
					getSubmission(db, c)
				},
			},
			{
				Method: "UPDATE",
				Path:   "/submission",
				Handler: func(c *gin.Context) {
					updateSubmission(db, c)
				},
			},
			{
				Method: "DELETE",
				Path:   "/submission",
				Handler: func(c *gin.Context) {
					deleteSubmission(db, c)
				},
			},
		},
	}
}

func createSubmission(db *gorm.DB, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "create submission",
	})
}

func getSubmission(db *gorm.DB, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "get submission",
	})
}

func updateSubmission(db *gorm.DB, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "update submission",
	})
}

func deleteSubmission(db *gorm.DB, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "delete submission",
	})
}
