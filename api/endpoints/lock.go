package endpoints

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/models/entry"
	"github.com/adraismawur/mibig-submission/models/lock"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func init() {
	RegisterEndpointGenerator(LockEndpoint)
}

func LockEndpoint(db *gorm.DB) Endpoint {
	return Endpoint{
		Routes: []Route{
			{
				Method: http.MethodGet,
				Path:   "/lock/list/:accession",
				Handler: func(c *gin.Context) {
					getActiveEntryLocks(db, c)
				},
			},
			{
				Method: http.MethodPost,
				Path:   "/lock/check",
				Handler: func(c *gin.Context) {
					checkEntryLocks(db, c)
				},
			},
			{
				Method: http.MethodPost,
				Path:   "/lock/request",
				Handler: func(c *gin.Context) {
					requestEntryLock(db, c)
				},
			},
			{
				Method: http.MethodPost,
				Path:   "/lock/release",
				Handler: func(c *gin.Context) {
					releaseEntryLock(db, c)
				},
			},
		},
	}
}

type LockRequest struct {
	Accession string               `json:"accession"`
	Category  lock.LockingCategory `json:"category"`
}

func getActiveEntryLocks(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")

	if accession == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing accession parameter"})
		return
	}

	activeEntry, err := entry.GetEntryFromAccession(db, accession)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	locks, err := lock.GetEntryLocks(db, int(activeEntry.ID))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	user, err := models.GetUserFromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	type LockInfo struct {
		Category  lock.LockingCategory `json:"category" gorm:"uniqueIndex:compositeLockIndex"`
		UnlocksAt time.Time            `json:"unlocks_at"`
		IsOwner   bool                 `json:"is_owner"`
	}

	var response = make([]LockInfo, len(*locks))
	for _, lockEntry := range *locks {
		response = append(response, LockInfo{
			Category:  lockEntry.Category,
			UnlocksAt: lockEntry.UnlocksAt,
			IsOwner:   lockEntry.LockOwnerID == user.ID,
		})
	}

	c.JSON(http.StatusOK, response)
}

func checkEntryLocks(db *gorm.DB, c *gin.Context) {
	var request LockRequest

	err := c.ShouldBindJSON(&request)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	activeEntry, err := entry.GetEntryFromAccession(db, request.Accession)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	user, err := models.GetUserFromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	activeLock, err := lock.GetEntryLock(db, int(activeEntry.ID), request.Category)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if activeLock.UnlocksAt.Before(time.Now()) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "lock timed out"})
		return
	}

	if activeLock.LockOwnerID != user.ID {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "lock does not belong to user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"unlocks_at": activeLock.UnlocksAt})
}

func requestEntryLock(db *gorm.DB, c *gin.Context) {
	var request LockRequest

	err := c.ShouldBindJSON(&request)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	activeEntry, err := entry.GetEntryFromAccession(db, request.Accession)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	user, err := models.GetUserFromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	activeLock, err := lock.CreateOrGetLock(db, int(activeEntry.ID), request.Category, *user)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"unlocks_at": activeLock.UnlocksAt})
}

func releaseEntryLock(db *gorm.DB, c *gin.Context) {

	var request LockRequest

	err := c.ShouldBindJSON(&request)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	activeEntry, err := entry.GetEntryFromAccession(db, request.Accession)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	user, err := models.GetUserFromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	err = lock.ReleaseLock(db, int(activeEntry.ID), request.Category, *user)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.Status(http.StatusOK)
}
