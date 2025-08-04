package endpoints

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func init() {
	RegisterEndpointGenerator(SubmissionEndpoint)
}

// SubmissionEndpoint returns the entry endpoint.
// This endpoint will implement creating and updating submissions, as well as perform some
// specific checks on submissions.
func SubmissionEndpoint(db *gorm.DB) Endpoint {
	return Endpoint{
		Routes: []Route{
			{
				Method: "CREATE",
				Path:   "/entry",
				Handler: func(c *gin.Context) {
					createSubmission(db, c)
				},
			},
			{
				Method: "GET",
				Path:   "/entry",
				Handler: func(c *gin.Context) {
					getSubmission(db, c)
				},
			},
			{
				Method: "UPDATE",
				Path:   "/entry",
				Handler: func(c *gin.Context) {
					updateSubmission(db, c)
				},
			},
			{
				Method: "DELETE",
				Path:   "/entry",
				Handler: func(c *gin.Context) {
					deleteSubmission(db, c)
				},
			},
		},
	}
}

func createSubmission(db *gorm.DB, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "create entry",
	})
}

func getSubmission(db *gorm.DB, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "get entry",
	})
}

func updateSubmission(db *gorm.DB, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "update entry",
	})
}

func deleteSubmission(db *gorm.DB, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "delete entry",
	})
}
