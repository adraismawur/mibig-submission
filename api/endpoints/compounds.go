package endpoints

import (
	"github.com/adraismawur/mibig-submission/models/entry/compound"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
				Method: "GET",
				Path:   "/entry/:accession/compounds/:idx",
				Handler: func(c *gin.Context) {
					getEntryCompound(db, c)
				},
			},
		},
	}
}

func getEntryCompounds(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")

	var compounds *[]compound.Compound

	err := db.Table("compounds").
		Where("entry_id = (select id from entries where accession = ?)", accession).
		Find(&compounds).
		Error

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"compounds": compounds})
}

func getEntryCompound(db *gorm.DB, c *gin.Context) {

}
