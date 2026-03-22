package endpoints

import (
	"github.com/adraismawur/mibig-submission/models/entry"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
)

func init() {
	RegisterEndpointGenerator(FinalizeEndpointGenerator)
}

func FinalizeEndpointGenerator(db *gorm.DB) Endpoint {
	return Endpoint{
		Routes: []Route{
			{
				Method: http.MethodPost,
				Path:   "/entry/:accession/finalize",
				Handler: func(c *gin.Context) {
					updateFinalDetails(db, c)
				},
			},
		},
	}
}

func updateFinalDetails(db *gorm.DB, c *gin.Context) {
	var finalDetails entry.FinalDetails

	err := c.ShouldBindJSON(&finalDetails)

	if err != nil {
		slog.Error("[endpoints] [finalize] Could not bind final details json")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = entry.UpdateFinalDetails(db, finalDetails)

	if err != nil {
		slog.Error("[endpoints] [finalize] Could not update final details")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
