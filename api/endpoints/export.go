package endpoints

import (
	"github.com/adraismawur/mibig-submission/models/entry"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
)

func init() {
	RegisterEndpointGenerator(ExportEndpoint)
}

func ExportEndpoint(db *gorm.DB) Endpoint {
	return Endpoint{
		Routes: []Route{
			{
				Method: "GET",
				Path:   "/export/entry/:accession",
				Handler: func(c *gin.Context) {
					getEntryExport(db, c)
				},
			},
		},
	}
}

func getEntryExport(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")
	formatJson := c.Query("pretty") == "true"

	exists, err := entry.GetEntryExists(db, accession)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "entry not found"})
		return
	}

	existingEntry, err := entry.GetEntryExportFromAccession(db, accession)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	if !formatJson {
		c.JSON(http.StatusOK, existingEntry)
		return
	}

	formattedJson, err := json.MarshalIndent(existingEntry, "", "  ")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		slog.Error("[endpoints] [entry] Failed to marshal existing entry", "error", err.Error())
		return
	}

	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, string(formattedJson))
}
