package endpoints

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/models/entry"
	"github.com/adraismawur/mibig-submission/models/entry/biosynthesis"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"strconv"
)

func init() {
	RegisterEndpointGenerator(BiosynthModuleEndpoint)
}

func BiosynthModuleEndpoint(db *gorm.DB) Endpoint {
	return Endpoint{
		Routes: []Route{
			{
				Method: http.MethodGet,
				Path:   "/entry/:accession/biosynth/module/:id",
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
				Path:   "/entry/:accession/biosynth/module_reorder",
				Handler: func(c *gin.Context) {
					reorderBiosynthModules(db, c)
				},
			},
			{
				Method: http.MethodPost,
				Path:   "/entry/:accession/biosynth/module/:id",
				Handler: func(c *gin.Context) {
					updateEntryBiosynthesisModule(db, c)
				},
			},
			{
				Method: http.MethodDelete,
				Path:   "/entry/:accession/biosynth/module/:id",
				Handler: func(c *gin.Context) {
					deleteEntryBiosynthesisModule(db, c)
				},
			},
		},
	}
}

func createEntryBiosynthesisModule(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")

	var module biosynthesis.BiosyntheticModule
	err := c.ShouldBindJSON(&module)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := models.GetUserFromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err := entry.GetEntryExists(db, accession)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "entry not found"})
		return
	}

	err = biosynthesis.CreateEntryBiosynthesisModule(db, accession, module)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	entry.AddContributor(db, accession, user.ID)

	c.Status(http.StatusOK)
}

func reorderBiosynthModules(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")
	type ReorderRequest struct {
		IDFrom uint64 `json:"id_from"`
		IDTo   uint64 `json:"id_to"`
	}

	var request ReorderRequest

	user, err := models.GetUserFromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.ShouldBindJSON(&request)

	if err != nil {
		slog.Error("[endpoints] [biosynth] Failed to unmarshal reorder request", "error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = biosynthesis.ReorderEntryBiosynthesisModules(db, request.IDFrom, request.IDTo)

	if err != nil {
		slog.Error("[endpoints] [biosynth] Could not reorder modules", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	entry.AddContributor(db, accession, user.ID)

	c.Status(http.StatusOK)
}

func updateEntryBiosynthesisModule(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")
	moduleId := c.Param("id")

	iModuleId, err := strconv.Atoi(moduleId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := models.GetUserFromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var module biosynthesis.BiosyntheticModule
	err = c.ShouldBindJSON(&module)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("[endpoints] [biosynth] Failed to marshal existing module", "error", err.Error())
		return
	}

	if uint64(iModuleId) != module.ID {
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

	err = biosynthesis.UpdateEntryBiosynthesisModule(db, module)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("[endpoints] [biosynth] Failed to update biosynthesis module", "error", err.Error())
		return
	}

	entry.AddContributor(db, accession, user.ID)

	c.Status(http.StatusOK)
}

func getEntryBiosynthesisModule(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")
	moduleId := c.Param("id")

	iModuleId, err := strconv.Atoi(moduleId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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

	entryBioSynth, err := biosynthesis.GetEntryBiosynthesisModule(db, iModuleId)

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
		slog.Error("[endpoints] [biosynth] Failed to marshal existing entry", "error", err.Error())
		return
	}

	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, string(formattedJson))
}

func deleteEntryBiosynthesisModule(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")
	moduleId := c.Param("id")

	iModuleId, err := strconv.Atoi(moduleId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := models.GetUserFromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err := entry.GetEntryExists(db, accession)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "entry not found"})
	}

	err = biosynthesis.DeleteEntryBiosynthesisModule(db, iModuleId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	entry.AddContributor(db, accession, user.ID)

	c.Status(http.StatusOK)
}
