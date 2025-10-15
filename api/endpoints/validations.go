package endpoints

import (
	"github.com/adraismawur/mibig-submission/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func init() {
	RegisterEndpointGenerator(GenbankValidationEndpoint)
}

// GenbankValidationEndpoint returns the auth endpoint, which is responsible for specifically handling authentication.
// This means acquiring a token (logging in) and refreshing a token.
func GenbankValidationEndpoint(db *gorm.DB) Endpoint {
	return Endpoint{
		Routes: []Route{
			{
				Method: "GET",
				Path:   "/validations/genbank/:accession",
				Handler: func(c *gin.Context) {
					validateGenBank(db, c)
				},
			},
		},
	}
}

func validateGenBank(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")

	result, err := util.GetGenbankAccessionSummary(accession)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	if result == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Accession not found"})
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}
