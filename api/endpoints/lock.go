package endpoints

import (
	"fmt"
	"github.com/adraismawur/mibig-submission/models"
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
			{
				Method: http.MethodPost,
				Path:   "/lock/clear/:accession",
				Handler: func(c *gin.Context) {
					clearEntryLock(db, c)
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

	locks, err := lock.GetEntryLocks(db, accession)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := models.GetUserFromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := models.GetUserFromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	activeLock, err := lock.GetEntryLock(db, request.Accession, request.Category)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if activeLock.ID == 0 {
		errorString := fmt.Sprintf("user does not have lock for entry %s category %s", request.Accession, request.Category)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errorString})
		return
	}

	if activeLock.UnlocksAt.Before(time.Now()) {
		errorString := fmt.Sprintf("lock timed out on entry %s category %s", request.Accession, request.Category)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errorString})
		return
	}

	if activeLock.LockOwnerID != user.ID {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "lock does not belong to user"})
		return
	}

	var response struct {
		UnlocksAt time.Time            `json:"unlocks_at"`
		Category  lock.LockingCategory `json:"category"`
		Full      bool                 `json:"full"`
	}

	response.UnlocksAt = activeLock.UnlocksAt
	response.Category = activeLock.Category
	response.Full = activeLock.Category == lock.Full

	c.JSON(http.StatusOK, response)
}

func requestEntryLock(db *gorm.DB, c *gin.Context) {
	var request LockRequest

	err := c.ShouldBindJSON(&request)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := models.GetUserFromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	activeLock, err := lock.CreateOrGetLock(db, request.Accession, request.Category, *user)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"unlocks_at": activeLock.UnlocksAt})
}

func releaseEntryLock(db *gorm.DB, c *gin.Context) {

	var request LockRequest

	err := c.ShouldBindJSON(&request)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := models.GetUserFromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = lock.ReleaseLock(db, request.Accession, request.Category, *user)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func clearEntryLock(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")

	user, err := models.GetUserFromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = lock.ClearLocks(db, accession, *user)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
