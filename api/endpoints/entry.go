package endpoints

import (
	"github.com/adraismawur/mibig-submission/middleware"
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/models/entry"
	"github.com/adraismawur/mibig-submission/util/constants"
	"github.com/adraismawur/mibig-submission/util/entry_utils"
	"github.com/beevik/guid"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"strconv"
	"time"
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

// createEntry creates a minimal entry from a request
func createEntry(db *gorm.DB, c *gin.Context) {
	var newEntry entry.Entry

	if err := c.BindJSON(&newEntry); err != nil {
		slog.Error("[endpoints] [submission] Failed to unmarshal submission JSON", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid entry submitted"})
		return
	}

	var currentDate = time.Now().Format(time.DateOnly)

	newEntry.Changelog = entry.Changelog{
		Releases: []entry.Release{
			{
				Version: "1",
				Date:    currentDate,
				Entries: []entry.ReleaseEntry{
					{
						Contributors: []string{
							constants.AnonymousUserId,
						},
						Reviewers: nil,
						Date:      currentDate,
						Comment:   constants.NewEntryComment,
					},
				},
			},
		},
	}

	bearerToken, err := middleware.GetAuthHeaderToken(c)

	if err != nil {
		return
	}

	user, err := models.GetUserFromToken(bearerToken)

	if err != nil {
		slog.Error("[endpoints] [submission] Failed to generate new entry accession", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate new entry accession"})
		return
	}

	newEntry.Accession = entry_utils.GeneratePlaceholderAccession(*user)

	db.Create(&newEntry)

	antismashTask := models.AntismashRun{
		Accession: newEntry.Loci[0].Accession,
		BGCID:     newEntry.Accession,
		GUID:      guid.NewString(),
	}

	err = db.Create(antismashTask).Error

	if err != nil {
		slog.Error("[endpoints] [entry] Failed to create antismash task", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create antismash task"})
	}

	c.JSON(http.StatusOK, gin.H{"status": antismashTask})
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
