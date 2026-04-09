package endpoints

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/models/entry"
	"github.com/adraismawur/mibig-submission/models/entry/compound"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
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
				Method: http.MethodGet,
				Path:   "/entry/:accession/compounds",
				Handler: func(c *gin.Context) {
					getEntryCompounds(db, c)
				},
			},
			{
				Method: http.MethodPost,
				Path:   "/entry/:accession/compounds",
				Handler: func(c *gin.Context) {
					createEntryCompound(db, c)
				},
			},
			{
				Method: http.MethodPost,
				Path:   "/entry/:accession/compounds/:compoundId",
				Handler: func(c *gin.Context) {
					updateEntryCompound(db, c)
				},
			},
			{
				Method: http.MethodDelete,
				Path:   "/entry/:accession/compounds/:compoundId",
				Handler: func(c *gin.Context) {
					deleteEntryCompound(db, c)
				},
			},
		},
	}
}

func getEntryCompounds(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")

	id := c.Query("id")
	formatJson := c.Query("pretty") == "true"

	var response struct {
		Compounds []compound.Compound `json:"compounds"`
	}

	q := db.Table("compounds").
		Preload("Evidence").
		Preload("BioActivities.Assays.Measurement").
		Preload("BioActivities.Assays.TestSystem").
		Where("entry_accession = $1", accession)

	if id != "" {
		q = q.Where("compounds.id = $2", id)
	}

	err := q.Find(&response.Compounds).
		Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !formatJson {
		c.JSON(http.StatusOK, response)
		return
	}

	formattedJson, err := json.MarshalIndent(response, "", "  ")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		slog.Error("[endpoints] [compound] Failed to marshal existing compound", "error", err.Error())
		return
	}

	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, string(formattedJson))
}

func createEntryCompound(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")

	var newCompound compound.Compound

	err := c.ShouldBindJSON(&newCompound)

	if err != nil {
		slog.Error("[endpoints] [compound] Could not bind compound json", "error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := models.GetUserFromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newCompound.EntryAccession = accession

	err = db.Create(&newCompound).Error

	if err != nil {
		slog.Error("[endpoints] [compound] Could not update compound", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	entry.AddContributor(db, accession, user.ID)

	c.JSON(http.StatusOK, gin.H{"compound": newCompound})
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

	user, err := models.GetUserFromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newCompound.EntryAccession = accession

	err = db.Session(&gorm.Session{FullSaveAssociations: true}).
		Table("compounds").
		Save(&newCompound).
		Error

	if err != nil {
		slog.Error("[Endpoints] [Compound] Could not update compound", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = db.Model(&newCompound).
		Association("Evidence").
		Replace(&newCompound.Evidence)

	if err != nil {
		slog.Error("[Endpoints] [Compound] Could not update compound evidence", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = db.Model(&newCompound).
		Association("BioActivities").
		Replace(&newCompound.BioActivities)

	if err != nil {
		slog.Error("[endpoints] [compound] Could not update compound bioactivities", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	entry.AddContributor(db, accession, user.ID)

	c.JSON(http.StatusOK, gin.H{"compound": newCompound})
}

func deleteEntryCompound(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")
	compoundId := c.Param("compoundId")

	user, err := models.GetUserFromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var compounds compound.Compound

	err = db.
		Model(&compounds).
		Delete("id = $1", compoundId).
		Error

	//if err != nil {
	//	slog.Error("[endpoints] [compound] Could find compound to delete", "error", err)
	//	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}
	//
	//err = db..Delete(compounds).Error

	if err != nil {
		slog.Error("[endpoints] [compound] Could not delete compound", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	entry.AddContributor(db, accession, user.ID)

	c.Status(http.StatusOK)
}
