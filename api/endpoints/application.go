package endpoints

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"time"
)

func init() {
	RegisterEndpointGenerator(applicationEndpoint)
}

// AuthEndpoint returns the auth endpoint, which is responsible for specifically handling authentication.
// This means acquiring a token (logging in) and refreshing a token.
func applicationEndpoint(db *gorm.DB) Endpoint {
	return Endpoint{
		Routes: []Route{
			{
				Method: "GET",
				Path:   "/application/state",
				Handler: func(c *gin.Context) {
					getState(db, c)
				},
			},
		},
	}
}

func getState(db *gorm.DB, c *gin.Context) {

	var stateResponse struct {
		State         models.ApplicationState          `json:"state"`
		Announcements []models.ApplicationAnnouncement `json:"announcements"`
	}

	err := db.First(&stateResponse.State).Error

	if err != nil {
		slog.Error("Could not get application state")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	now := time.Now().Unix() * 1000

	err = db.Where("expires_at > ?", now).Find(&stateResponse.Announcements).Error

	if err != nil {
		slog.Error("Could not get application announcements")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stateResponse)
}
