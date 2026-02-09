package endpoints

import (
	"github.com/adraismawur/mibig-submission/models/entry"
	"github.com/adraismawur/mibig-submission/models/entry/biosynthesis"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
)

func init() {
	RegisterEndpointGenerator(BiosynthEndpoint)
}

func BiosynthEndpoint(db *gorm.DB) Endpoint {
	return Endpoint{
		Routes: []Route{
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
				Method: http.MethodPost,
				Path:   "/entry/:accession/biosynth/module",
				Handler: func(c *gin.Context) {
					createEntryBiosynthesisModule(db, c)
				},
			},
			{
				Method: http.MethodPost,
				Path:   "/entry/:accession/biosynth/module/:name",
				Handler: func(c *gin.Context) {
					updateEntryBiosynthesisModule(db, c)
				},
			},
			{
				Method: "DELETE",
				Path:   "/entry/:accession/biosynth/module/:name",
				Handler: func(c *gin.Context) {
					deleteEntryBiosynthesisModule(db, c)
				},
			},
		},
	}
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

	entryBioSynth, err := biosynthesis.GetEntryBiosynthesis(db, accession)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	c.JSON(http.StatusOK, entryBioSynth)
}

func createEntryBiosynthesisModule(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")

	var module biosynthesis.BiosyntheticModule
	c.ShouldBindJSON(&module)

	exists, err := entry.GetEntryExists(db, accession)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "entry not found"})
		return
	}

	err = biosynthesis.CreateEntryBiosynthesisModule(db, accession, module)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.Status(http.StatusOK)
}

func updateEntryBiosynthesisModule(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")
	name := c.Param("name")

	var module biosynthesis.BiosyntheticModule
	err := c.ShouldBindJSON(&module)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		slog.Error("[endpoints] [entry] Failed to marshal existing module", "error", err.Error())
		return
	}

	if name != module.Name {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name mismatch between request URL and data"})
		return
	}

	exists, err := entry.GetEntryExists(db, accession)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		slog.Error("[endpoints] [entry] Error finding entry", "error", err.Error())
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "entry not found"})
		return
	}

	err = biosynthesis.UpdateEntryBiosynthesisModule(db, accession, &module)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		slog.Error("[endpoints] [entry] Failed to update biosynthesis module", "error", err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func getEntryBiosynthesisModule(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")
	name := c.Param("name")

	formatJson := false

	if c.Query("pretty") == "true" {
		formatJson = true
	}

	exists, err := entry.GetEntryExists(db, accession)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "entry not found"})
		return
	}

	entryBioSynth, err := biosynthesis.GetEntryBiosynthesisModule(db, accession, name)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if !formatJson {
		c.JSON(http.StatusOK, entryBioSynth)
		return
	}

	formattedJson, err := json.MarshalIndent(entryBioSynth, "", "  ")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		slog.Error("[endpoints] [entry] Failed to marshal existing entry", "error", err.Error())
		return
	}

	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, string(formattedJson))
}

func deleteEntryBiosynthesisModule(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")
	name := c.Param("name")

	exists, err := entry.GetEntryExists(db, accession)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "entry not found"})
	}

	err = biosynthesis.DeleteEntryBiosynthesisModule(db, accession, name)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.Status(http.StatusOK)
}
