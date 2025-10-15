package endpoints

import (
	"github.com/adraismawur/mibig-submission/util"
	"github.com/beevik/guid"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
)

func init() {
	RegisterEndpointGenerator(AntismashEndpoint)
}

// AntismashEndpoint returns the antismash endpoint, used for submitting, checking and stopping antismash runs
func AntismashEndpoint(db *gorm.DB) Endpoint {
	return Endpoint{
		Routes: []Route{
			//{
			//	Method: http.MethodGet,
			//	Path:   "/antismash/:accession",
			//	Handler: func(c *gin.Context) {
			//		getAntismashResults(db, c)
			//	},
			//},
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

func getAntismashResults(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")

	if accession == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no accession given"})
		return
	}

	//outputDir := path2.Join(config.Envs["DATA_PATH"], "antismash", accession)
	//
	//c.FileFromFS()
}

func startAntismashRun(db *gorm.DB, c *gin.Context) {
	request := util.AntismashRun{}

	if err := c.ShouldBindJSON(&request); err != nil {
		slog.Error("[Antismash] Could not bind request json")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request.GUID = guid.NewString()

	slog.Info("[Antismash] Starting Antismash Run", "accession", request.Accession)

	db.Create(&request)
}
