package endpoints

import (
	"github.com/adraismawur/mibig-submission/models/entry/compound"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
)

func init() {
	RegisterEndpointGenerator(CompoundsEndpoint)
}

func CompoundsEndpoint(db *gorm.DB) Endpoint {
	return Endpoint{
		Routes: []Route{
			{
				Method: "GET",
				Path:   "/entry/:accession/compounds",
				Handler: func(c *gin.Context) {
					getEntryCompounds(db, c)
				},
			},
			{
				Method: "POST",
				Path:   "/entry/:accession/compounds",
				Handler: func(c *gin.Context) {
					createEntryCompound(db, c)
				},
			},
			{
				Method: "POST",
				Path:   "/entry/:accession/compounds/:compoundId",
				Handler: func(c *gin.Context) {
					updateEntryCompound(db, c)
				},
			},
		},
	}
}

func createEntryCompound(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")

	var newCompound compound.Compound

	err := c.ShouldBindJSON(&newCompound)

	if err != nil {
		slog.Error("[Endpoints] [Compound] Could not bind compound json", "error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var entryId uint64

	err = db.Table("entries").
		Where("accession = ?", accession).
		Select("id").
		Find(&entryId).
		Error

	if err != nil {
		slog.Error("[Endpoints] [Compound] Could not find entry", "error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newCompound.EntryID = entryId

	err = db.Create(&newCompound).Error

	if err != nil {
		slog.Error("[Endpoints] [Compound] Could not update compound", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"compound": newCompound})
}

func getEntryCompounds(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")

	id := c.Query("id")

	var compounds *[]compound.Compound

	q := db.Table("compounds").
		Preload("Evidence").
		Preload("BioActivities").
		Where("entry_id = (select id from entries where accession = ?)", accession)

	if id != "" {
		q = q.Where("id = ?", id)
	}

	err := q.Find(&compounds).
		Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, compounds)
}

func updateEntryCompound(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")

	var newCompound compound.Compound

	err := c.ShouldBindJSON(&newCompound)

	if err != nil {
		slog.Error("[Endpoints] [Compound] Could not bind compound json", "error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var entryId uint64

	err = db.Table("entries").
		Where("accession = ?", accession).
		Select("id").
		Find(&entryId).
		Error

	if err != nil {
		slog.Error("[Endpoints] [Compound] Could not find entry", "error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newCompound.EntryID = entryId

	err = db.Table("compounds").
		Where("entry_id = ?", entryId).
		Save(&newCompound).
		Error

	if err != nil {
		slog.Error("[Endpoints] [Compound] Could not update compound", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"compound": newCompound})
}
