package endpoints

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/models/entry"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"strconv"
)

func init() {
	RegisterEndpointGenerator(EntryEndpoint)
}

// EntryEndpoint returns the entry endpoint.
// This endpoint will implement creating and updating entries.
// This is distinct from creating, updating and deleting submissions, which are new proposed entries
func EntryEndpoint(db *gorm.DB) Endpoint {
	return Endpoint{
		Routes: []Route{
			{
				Method: "GET",
				Path:   "/entry",
				Handler: func(c *gin.Context) {
					listEntries(db, c)
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
				Path:   "/entry/:accession/biosynth",
				Handler: func(c *gin.Context) {
					getEntryBiosynthesis(db, c)
				},
			},
			{
				Method: "GET",
				Path:   "/entry/:accession/biosynth/module/:name",
				Handler: func(c *gin.Context) {
					getEntryBiosynthesisModule(db, c)
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

func listEntries(db *gorm.DB, c *gin.Context) {
	// listAll will include user submissions
	listAll := c.Query("list_all") == "true"

	var accessionList []string

	q := db.Table("entries")
	// select only accessions
	q.Select("entries.accession")

	if !listAll {
		q = q.Where("entries.id NOT IN (select entry_id from user_submissions)")
	}
	err := q.Find(&accessionList).Error

	if err != nil {
		slog.Error("[endpoints] [entry] Could not list entries", "error", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, accessionList)
}

func getEntry(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")
	formatJson := c.Query("pretty") == "true"

	exists, err := entry.GetEntryExists(db, accession)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "entry not found"})
		return
	}

	existingEntry, err := entry.GetEntryFromAccession(db, accession)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	if !formatJson {
		c.JSON(http.StatusOK, existingEntry)
		return
	}

	formattedJson, err := json.MarshalIndent(existingEntry, "", "  ")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		slog.Error("[endpoints] [entry] Failed to marshal existing entry", "error", err.Error())
		return
	}

	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, string(formattedJson))
}

func getEntryBiosynthesis(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")

	exists, err := entry.GetEntryExists(db, accession)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "entry not found"})
	}

	entryBioSynth, err := entry.GetEntryBiosynthesis(db, accession)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	c.JSON(http.StatusOK, entryBioSynth)
}

func getEntryBiosynthesisModule(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")
	name := c.Param("name")

	exists, err := entry.GetEntryExists(db, accession)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "entry not found"})
	}

	entryBioSynth, err := entry.GetEntryBiosynthesisModule(db, accession, name)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	c.JSON(http.StatusOK, entryBioSynth)
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
