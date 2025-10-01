package endpoints

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/models/entry"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"strconv"
)

func init() {
	RegisterEndpointGenerator(EntryEndpoint)
}

// EntryEndpoint returns the entry endpoint.
// This endpoint will implement creating and updating submissions, as well as perform some
// specific checks on submissions.
func EntryEndpoint(db *gorm.DB) Endpoint {
	return Endpoint{
		Routes: []Route{
			{
				Method: "POST",
				Path:   "/entry",
				Handler: func(c *gin.Context) {
					createEntry(db, c)
				},
			},
			{
				Method: "GET",
				Path:   "/entry/:accession",
				Handler: func(c *gin.Context) {
					getEntry(db, c)
				},
			},
			{
				Method: "GET",
				Path:   "/entry/user/:userId",
				Handler: func(c *gin.Context) {
					getUserentries(db, c)
				},
			},
			{
				Method: "UPDATE",
				Path:   "/entry",
				Handler: func(c *gin.Context) {
					updateEntry(db, c)
				},
			},
			{
				Method: "DELETE",
				Path:   "/entry",
				Handler: func(c *gin.Context) {
					deleteEntry(db, c)
				},
			},
		},
	}
}

// createEntry creates a minimal entry from a request
func createEntry(db *gorm.DB, c *gin.Context) {

	var entry entry.Entry

	if err := c.BindJSON(&entry); err != nil {
		slog.Error("[endpoints] [submission] Failed to unmarshal submission JSON", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid entry submitted"})
		return
	}

	db.Create(&entry)
}

func getEntry(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")

	exists, err := entry.GetEntryExists(db, accession)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "entry not found"})
	}

	existingEntry, err := entry.GetEntryFromAccession(db, accession)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	c.JSON(http.StatusOK, existingEntry)
}

func getUserentries(db *gorm.DB, c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("userId"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	exists, err := models.GetUserExistsByID(db, userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
	}

	accessions, err := entry.GetUserEntries(db, userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	c.JSON(http.StatusOK, accessions)
}

func updateEntry(db *gorm.DB, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "update entry",
	})
}

func deleteEntry(db *gorm.DB, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "delete entry",
	})
}
