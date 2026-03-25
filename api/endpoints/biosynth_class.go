package endpoints

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/models/entry"
	"github.com/adraismawur/mibig-submission/models/entry/biosynthesis"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"strconv"
)

func init() {
	RegisterEndpointGenerator(BiosynthClassEndpoint)
}

func BiosynthClassEndpoint(db *gorm.DB) Endpoint {
	return Endpoint{
		Routes: []Route{
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, class)
}

func createEntryBiosynthesisClass(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")

	var class biosynthesis.BiosyntheticClass
	err := c.ShouldBindJSON(&class)

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	entry.AddContributor(db, accession, user.ID)

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

	user, err := models.GetUserFromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var class biosynthesis.BiosyntheticClass
	err = c.ShouldBindJSON(&class)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	err = biosynthesis.UpdateEntryBiosynthesisClass(db, classId, class)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("[endpoints] [biosynth] Failed to update biosynthesis class", "error", err.Error())
		return
	}

	entry.AddContributor(db, accession, user.ID)

	c.Status(http.StatusOK)

}

func deleteEntryBiosynthesisClass(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")
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

	user, err := models.GetUserFromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = biosynthesis.DeleteEntryBiosynthesisClass(db, class_id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	entry.AddContributor(db, accession, user.ID)

	c.Status(http.StatusOK)
}
