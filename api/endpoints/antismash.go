package endpoints

import (
	"github.com/adraismawur/mibig-submission/middleware"
	"github.com/adraismawur/mibig-submission/models"
	"github.com/beevik/guid"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
)

func init() {
	RegisterEndpointGenerator(AntismashEndpoint)
	middleware.AddProtectedRoute(http.MethodPost, "/antismash", models.Admin)
}

// AntismashEndpoint returns the antismash endpoint, used for submitting, checking and stopping antismash runs
func AntismashEndpoint(db *gorm.DB) Endpoint {
	return Endpoint{
		Routes: []Route{
			{
				Method: http.MethodGet,
				Path:   "/antismash/:guid",
				Handler: func(c *gin.Context) {
					getAntismashStatus(db, c)
				},
			},
			{
				Method: http.MethodPost,
				Path:   "/antismash",
				Handler: func(c *gin.Context) {
					startAntismashRun(db, c)
				},
			},
		},
	}
}

func getAntismashStatus(db *gorm.DB, c *gin.Context) {
	taskGuid := c.Param("guid")

	if taskGuid == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no guid given"})
		return
	}

	status := models.AntismashRun{}

	err := db.Where("guid = ?", taskGuid).First(&status).Error

	if err != nil {
		slog.Error("[Antismash] Get Antismash Status Error", "err", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Antismash Status Error"})
		return
	}

	c.JSON(http.StatusOK, status)
}

func startAntismashRun(db *gorm.DB, c *gin.Context) {
	request := models.AntismashRun{}

	if err := c.ShouldBindJSON(&request); err != nil {
		slog.Error("[Antismash] Could not bind request json")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request.GUID = guid.NewString()

	slog.Info("[Antismash] Starting Antismash Run", "accession", request.Accession)

	db.Create(&request)
}
