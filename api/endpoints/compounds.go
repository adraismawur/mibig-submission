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
					updateEntryCompound(db, c)
				},
			},
		},
	}
}

func getEntryCompounds(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")

	id := c.Query("id")

	var compounds *[]compound.Compound

	q := db.Table("compounds").
		Where("entry_id = (select id from entries where accession = ?)", accession)

	if id != "" {
		q = q.Where("id = ?", id)
	}

	err := q.Find(&compounds).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"compounds": compounds})
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

	err = db.Table("compounds").
		Where("entry_id = (select id from entries where accession = ?)", accession).
		Save(&newCompound).
		Error

	if err != nil {
		slog.Error("[Endpoints] [Compound] Could not update compound", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"compound": newCompound})
}
