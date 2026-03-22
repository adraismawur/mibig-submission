package endpoints

import (
	"github.com/adraismawur/mibig-submission/config"
	"github.com/adraismawur/mibig-submission/models/entry"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	path2 "path"
	"strconv"
)

func init() {
	RegisterEndpointGenerator(EntryEndpoint)
}

// EntryEndpoint returns the entry endpoint.
// This endpoint will implement creating and updating entries.
// This is distinct from creating, updating and deleting submissions, which are new proposed entries
func EntryEndpoint(db *gorm.DB) Endpoint {
	return Endpoint{
		Routes: []Route{
			{
				Method: "GET",
				Path:   "/entry",
				Handler: func(c *gin.Context) {
					listEntries(db, c)
				},
			},
			{
				Method: "GET",
				Path:   "/entry/:accession/genes",
				Handler: func(c *gin.Context) {
					getEntryGeneList(db, c)
				},
			},
			{
				Method: "GET",
				Path:   "/entry/:accession",
				Handler: func(c *gin.Context) {
					getEntry(db, c)
				},
			},
			{
				Method: "GET",
				Path:   "/entry/raw/:accession",
				Handler: func(c *gin.Context) {
					getRawEntry(db, c)
				},
			},
			{
				Method: "DELETE",
				Path:   "/entry",
				Handler: func(c *gin.Context) {
					deleteEntry(db, c)
				},
			},
		},
	}
}

func listEntries(db *gorm.DB, c *gin.Context) {
	// listAll will include user submissions
	listAll := c.Query("list_all") == "true"
	start, err := strconv.Atoi(c.Query("start"))

	if err != nil {
		start = 0
	}

	limit, err := strconv.Atoi(c.Query("limit"))

	if err != nil {
		limit = 20
	}

	search := c.Query("search")

	type ExistingEntrySummary struct {
		Accession    string `json:"accession"`
		Status       string `json:"status"`
		Completeness string `json:"completeness"`
	}

	var existingEntries []ExistingEntrySummary

	q := db.Table("entries")
	if !listAll {
		q = q.Where("entries.accession NOT IN (select entry_accession from user_submissions)")
	}

	if search != "" {
		q = q.Where("entries.accession LIKE $1", "%"+search+"%")
	}

	q = q.Offset(start)
	q = q.Limit(limit)
	err = q.Find(&existingEntries).Error

	if err != nil {
		slog.Error("[endpoints] [entry] Could not list entries", "error", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	slog.Info(q.Statement.SQL.String())

	var recordCount int64

	q = db.Table("entries")

	if search != "" {
		q = q.Where("entries.accession LIKE $1", "%"+search+"%")
	}

	q = q.Count(&recordCount)

	err = q.Error

	if err != nil {
		slog.Error("[endpoints] [entry] Could not get record count", "error", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	var response struct {
		Entries     []ExistingEntrySummary `json:"entries"`
		RecordCount int64                  `json:"record_count"`
	}

	response.Entries = existingEntries
	response.RecordCount = recordCount

	c.JSON(http.StatusOK, response)
}

func getEntry(db *gorm.DB, c *gin.Context) {
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

	existingEntry, err := entry.GetEntryFromAccession(db, accession)

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

// getRawEntry returns the actual JSON from the data storage instead of the reconstructed JSON from the databse
func getRawEntry(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")

	dataPath, err := config.GetConfig(config.EnvDataPath)

	if err != nil {
		slog.Error("[endpoints] [entry] Could not get env variable for data path")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	jsonPath := path2.Join(dataPath, "json", accession+".json")

	c.File(jsonPath)
}

func deleteEntry(db *gorm.DB, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "delete entry",
	})
}

// getEntryGeneList gets a complete list of genes that are valid within this entry/submission
func getEntryGeneList(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")

	if accession == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bad accession"})
		return
	}

	exists, err := entry.GetEntryExists(db, accession)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !exists {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "accession not found"})
		return
	}

	genes, err := entry.GetEntryGenes(db, accession)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, genes)
}
