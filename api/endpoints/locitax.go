package endpoints

import (
	"github.com/adraismawur/mibig-submission/models/entry"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func init() {
	RegisterEndpointGenerator(LociTaxEndpoint)
}

func LociTaxEndpoint(db *gorm.DB) Endpoint {
	return Endpoint{
		Routes: []Route{
			{
				Method: http.MethodGet,
				Path:   "/entry/:accession/locitax",
				Handler: func(c *gin.Context) {
					getEntryLociTax(db, c)
				},
			},
			{
				Method: http.MethodPost,
				Path:   "/entry/:accession/locitax",
				Handler: func(c *gin.Context) {
					updateEntryLociTax(db, c)
				},
			},
		},
	}
}

func getEntryLociTax(db *gorm.DB, c *gin.Context) {

	var accession = c.Param("accession")

	exists, err := entry.GetEntryExists(db, accession)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "entry not found"})
		return
	}

	result, err := entry.GetLociTax(db, accession)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, result)
}

func updateEntryLociTax(db *gorm.DB, c *gin.Context) {
	var newLociTax entry.LociTax

	err := c.ShouldBindJSON(&newLociTax)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not unmarshal json"})
		return
	}

	var accession = c.Param("accession")

	exists, err := entry.GetEntryExists(db, accession)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "entry not found"})
		return
	}

	oldLociTax, err := entry.GetLociTax(db, accession)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not get locus and taxonomy information"})
		return
	}

	err = entry.UpdateLociTax(db, accession, *oldLociTax, newLociTax)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.Status(http.StatusOK)
}
