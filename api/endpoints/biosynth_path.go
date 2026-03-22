package endpoints

import (
	"github.com/adraismawur/mibig-submission/models/entry"
	"github.com/adraismawur/mibig-submission/models/entry/biosynthesis"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"strconv"
)

func init() {
	RegisterEndpointGenerator(BiosynthPathEndpoint)
}

func BiosynthPathEndpoint(db *gorm.DB) Endpoint {
	return Endpoint{
		Routes: []Route{
			{
				Method: http.MethodGet,
				Path:   "/entry/:accession/biosynth/path/:id",
				Handler: func(c *gin.Context) {
					getEntryBiosynthesisPath(db, c)
				},
			},
			{
				Method: http.MethodPost,
				Path:   "/entry/:accession/biosynth/path",
				Handler: func(c *gin.Context) {
					createEntryBiosynthesisPath(db, c)
				},
			},
			{
				Method: http.MethodPost,
				Path:   "/entry/:accession/biosynth/path/:id",
				Handler: func(c *gin.Context) {
					updateEntryBiosynthesisPath(db, c)
				},
			},
			{
				Method: http.MethodDelete,
				Path:   "/entry/:accession/biosynth/path/:id",
				Handler: func(c *gin.Context) {
					deleteEntryBiosynthesisPath(db, c)
				},
			},
		},
	}
}

func getEntryBiosynthesisPath(db *gorm.DB, c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad parameter: id"})
		return
	}

	path_id, err := strconv.Atoi(id)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not parse parameter: id"})
		return
	}

	path, err := biosynthesis.GetBiosynthesisPath(db, path_id)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, path)
}

func createEntryBiosynthesisPath(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")

	var path biosynthesis.BiosyntheticPathway
	err := c.ShouldBindJSON(&path)

	if err != nil {
		slog.Error("[endpoints] [biosynth_path] Could not bind JSON payload")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err := entry.GetEntryExists(db, accession)

	if err != nil {
		slog.Error("[endpoints] [biosynth_path] Could not retrieve entry")
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "entry not found"})
		return
	}

	err = biosynthesis.CreateBiosynthesisPath(db, path)

	if err != nil {
		slog.Error("[endpoints] [biosynth_path] Could not create biosynth path")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func updateEntryBiosynthesisPath(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")
	id := c.Param("id")

	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad parameter: id"})
		return
	}

	pathId, err := strconv.Atoi(id)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not parse parameter: id"})
		return
	}

	var path biosynthesis.BiosyntheticPathway
	err = c.ShouldBindJSON(&path)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("[endpoints] [biosynth] Failed to marshal existing path", "error", err.Error())
		return
	}

	if uint64(pathId) != path.ID {
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

	err = biosynthesis.UpdateBiosynthesisPath(db, path)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("[endpoints] [biosynth] Failed to update biosynthesis path", "error", err.Error())
		return
	}

	c.Status(http.StatusOK)

}

func deleteEntryBiosynthesisPath(db *gorm.DB, c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad parameter: id"})
		return
	}

	path_id, err := strconv.Atoi(id)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not parse parameter: id"})
		return
	}

	err = biosynthesis.DeleteBiosynthesisPath(db, path_id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
