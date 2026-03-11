package endpoints

import (
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
	RegisterEndpointGenerator(BiosynthEndpoint)
}

func BiosynthEndpoint(db *gorm.DB) Endpoint {
	return Endpoint{
		Routes: []Route{
			{
				Method: http.MethodGet,
				Path:   "/entry/:accession/biosynth",
				Handler: func(c *gin.Context) {
					getEntryBiosynthesis(db, c)
				},
			},
			{
				Method: http.MethodGet,
				Path:   "/entry/:accession/biosynth/class/:id",
				Handler: func(c *gin.Context) {
					getEntryBiosynthesisClass(db, c)
				},
			},
			{
				Method: http.MethodPost,
				Path:   "/entry/:accession/biosynth/class",
				Handler: func(c *gin.Context) {
					createEntryBiosynthesisClass(db, c)
				},
			},
			{
				Method: http.MethodPost,
				Path:   "/entry/:accession/biosynth/class/:id",
				Handler: func(c *gin.Context) {
					updateEntryBiosynthesisClass(db, c)
				},
			},
			{
				Method: http.MethodDelete,
				Path:   "/entry/:accession/biosynth/class/:id",
				Handler: func(c *gin.Context) {
					deleteEntryBiosynthesisClass(db, c)
				},
			},
			{
				Method: http.MethodGet,
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
				Path:   "/entry/:accession/biosynth/module_reorder",
				Handler: func(c *gin.Context) {
					reorderBiosynthModules(db, c)
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
				Method: http.MethodDelete,
				Path:   "/entry/:accession/biosynth/module/:name",
				Handler: func(c *gin.Context) {
					deleteEntryBiosynthesisModule(db, c)
				},
			},
		},
	}
}

func getEntryBiosynthesisClass(db *gorm.DB, c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad parameter: id"})
		return
	}

	class_id, err := strconv.Atoi(id)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not parse parameter: id"})
		return
	}

	class, err := biosynthesis.GetEntryBiosynthesisClass(db, class_id)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, class)
}

func createEntryBiosynthesisClass(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")

	var class biosynthesis.BiosyntheticClass
	err := c.ShouldBindJSON(&class)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	exists, err := entry.GetEntryExists(db, accession)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "entry not found"})
		return
	}

	biosynth, err := biosynthesis.GetEntryBiosynthesis(db, accession)

	if err != nil {
		slog.Error("[endpoints] [biosynth] could not find entry", "accession", accession)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	err = biosynthesis.CreateBiosynthesisClass(db, biosynth.ID, class)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.Status(http.StatusOK)
}

func updateEntryBiosynthesisClass(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")
	id := c.Param("id")

	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad parameter: id"})
		return
	}

	classId, err := strconv.Atoi(id)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not parse parameter: id"})
		return
	}

	var class biosynthesis.BiosyntheticClass
	err = c.ShouldBindJSON(&class)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		slog.Error("[endpoints] [biosynth] Failed to marshal existing class", "error", err.Error())
		return
	}

	if uint64(classId) != class.ID {
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

	err = biosynthesis.UpdateEntryBiosynthesisClass(db, classId, &class)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		slog.Error("[endpoints] [biosynth] Failed to update biosynthesis class", "error", err.Error())
		return
	}

	c.Status(http.StatusOK)

}

func deleteEntryBiosynthesisClass(db *gorm.DB, c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad parameter: id"})
		return
	}

	class_id, err := strconv.Atoi(id)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not parse parameter: id"})
		return
	}

	err = biosynthesis.DeleteEntryBiosynthesisClass(db, class_id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.Status(http.StatusOK)
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

func reorderBiosynthModules(db *gorm.DB, c *gin.Context) {
	type ReorderRequest struct {
		IDFrom uint64 `json:"id_from"`
		IDTo   uint64 `json:"id_to"`
	}

	var request ReorderRequest

	err := c.ShouldBindJSON(&request)

	if err != nil {
		slog.Error("[endpoints] [biosynth] Failed to unmarshal reorder request", "error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = biosynthesis.ReorderEntryBiosynthesisModules(db, request.IDFrom, request.IDTo)

	if err != nil {
		slog.Error("[endpoints] [biosynth] Could not reorder modules", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
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
		slog.Error("[endpoints] [biosynth] Failed to marshal existing module", "error", err.Error())
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
		slog.Error("[endpoints] [biosynth] Failed to update biosynthesis module", "error", err.Error())
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
		slog.Error("[endpoints] [biosynth] Failed to marshal existing entry", "error", err.Error())
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
