package endpoints

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/models/entry"
	"github.com/adraismawur/mibig-submission/models/lock"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"strconv"
)

func init() {
	RegisterEndpointGenerator(SubmissionEndpoint)
}

type SubmissionInfo struct {
	Accession       string                `json:"accession"`
	Type            entry.SubmissionType  `json:"type"`
	SourceAccession string                `json:"source_accession"`
	State           entry.SubmissionState `json:"state"`
	Category        lock.LockingCategory  `json:"category"`
}

// SubmissionEndpoint returns the submission endpoint.
// This endpoint will implement creating and updating submissions, as well as perform some
// specific checks on submissions.
func SubmissionEndpoint(db *gorm.DB) Endpoint {
	return Endpoint{
		Routes: []Route{
			{
				Method: "GET",
				Path:   "/submission/:userId",
				Handler: func(c *gin.Context) {
					getUserSubmissions(db, c)
				},
			},
			{
				Method: "GET",
				Path:   "/submission/",
				Handler: func(c *gin.Context) {
					getSubmissions(db, c)
				},
			},
			{
				Method: "POST",
				Path:   "/submission",
				Handler: func(c *gin.Context) {
					createNewSubmission(db, c)
				},
			},
			{
				Method: "POST",
				Path:   "/mutation",
				Handler: func(c *gin.Context) {
					createNewMutation(db, c)
				},
			},
			{
				Method: "POST",
				Path:   "/submission/promote/:accession",
				Handler: func(c *gin.Context) {
					promoteSubmission(db, c)
				},
			},
			{
				Method: "POST",
				Path:   "/submission/redraft/:accession",
				Handler: func(c *gin.Context) {
					redraftSubmission(db, c)
				},
			},
			{
				Method: "POST",
				Path:   "/submission/discard/:accession",
				Handler: func(c *gin.Context) {
					discardSubmission(db, c)
				},
			},
		},
	}
}

func getUserSubmissions(db *gorm.DB, c *gin.Context) {
	var submissions []SubmissionInfo

	userID := c.Param("userId")

	q := db.Table("user_submissions").
		Joins("JOIN entries ON entries.accession = user_submissions.entry_accession")

	// optional clause
	if userID != "" {
		q.Where("user_submissions.user_id = ?", userID)
	}

	q.Where("state != ?", entry.Discarded)

	err := q.Select("entries.accession, user_submissions.type, user_submissions.source_accession, user_submissions.state").
		Find(&submissions).Error

	if err != nil {
		slog.Error("[endpoints] [submission] Could not find submissions for user", "user_id", userID, "error", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, submissions)
}

func getSubmissions(db *gorm.DB, c *gin.Context) {
	var submissions []SubmissionInfo

	userID := c.Query("id")
	state := c.Query("state")

	q := db.Table("user_submissions").
		Joins("JOIN entries ON entries.accession = user_submissions.entry_accession")

	// optional clause
	if userID != "" {
		q.Where("user_submissions.user_id = ?", userID)
	}

	if state != "" {
		q.Where("user_submissions.state = ?", state)
	}

	err := q.Select("entries.accession, user_submissions.type, user_submissions.source_accession, user_submissions.state").
		Find(&submissions).Error

	if err != nil {
		slog.Error("[endpoints] [submission] Could not find submissions", "user_id", userID, "error", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, submissions)
}

// createNewSubmission creates a minimal draft submission from a request
func createNewSubmission(db *gorm.DB, c *gin.Context) {
	var minimalEntry entry.MinimalEntry

	if err := c.BindJSON(&minimalEntry); err != nil {
		slog.Error("[endpoints] [submission] Failed to unmarshal submission JSON", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid entry submitted"})
		return
	}

	user, err := models.GetUserFromContext(c)

	if err != nil {
		slog.Error("[endpoints] [submission] Failed to retrieve user from context", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve user from context"})
		return
	}

	newEntry, err := entry.CreateNewUserSubmission(db, minimalEntry, *user)

	if err != nil {
		slog.Error("[endpoints] [submission] Failed to create user submission record", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user submission record"})
		return
	}

	_, err = lock.CreateOrGetLock(db, newEntry.Accession, "full", *user)

	if err != nil {
		slog.Error("[endpoints] [submission] Failed to create lock for user submission", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user submission lock"})
		return
	}

	antismashTask, err := entry.CreateAntismashWorkerTask(db, *newEntry)

	if err != nil {
		slog.Error("[endpoints] [entry] Failed to create antismash task", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create antismash task"})
	}

	c.JSON(http.StatusOK, gin.H{"status": antismashTask})
}

func createNewMutation(db *gorm.DB, c *gin.Context) {
	type NewMutationRequest struct {
		FromAccession string `json:"from_accession"`
	}

	var request NewMutationRequest

	if err := c.BindJSON(&request); err != nil {
		slog.Error("[endpoints] [submission] Failed to unmarshal mutation request JSON", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	user, err := models.GetUserFromContext(c)

	if err != nil {
		slog.Error("[endpoints] [submission] Failed to retrieve user from context", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve user from context"})
		return
	}

	newEntry, err := entry.CreateNewUserMutation(db, request.FromAccession, *user)

	if err != nil {
		slog.Error("[endpoints] [submission] Failed to create user submission record", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user submission record"})
		return
	}

	_, err = lock.CreateOrGetLock(db, newEntry.Accession, "full", *user)

	if err != nil {
		slog.Error("[endpoints] [submission] Failed to create lock for user submission", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user submission lock"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accession": newEntry.Accession})
}

func getUserEntries(db *gorm.DB, c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("userId"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	exists, err := models.GetUserExistsByID(db, userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
	}

	accessions, err := GetUserSubmissions(db, userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	c.JSON(http.StatusOK, accessions)
}

func GetUserSubmissions(db *gorm.DB, userId int) ([]string, error) {
	var accessions []string

	err := db.
		Table("user_submissions").
		Select("accession").
		Joins("JOIN entries ON entries.accession = user_submissions.entry_accession").
		Where("user_id = ?", userId).
		Find(&accessions).
		Error

	if err != nil {
		return nil, err
	}

	return accessions, nil
}

func promoteSubmission(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")

	exists, err := entry.GetEntryExists(db, accession)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "entry not found"})
		return
	}

	var userSubmission entry.UserSubmission

	err = db.
		Joins("JOIN entries ON entries.accession = user_submissions.entry_accession").
		Where("entries.accession = ?", accession).
		First(&userSubmission).Error

	if err != nil {
		slog.Error("[endpoints] [submission] Failed to find user submission", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user submission"})
		return
	}

	if userSubmission.State != entry.Draft {
		slog.Error("[endpoints] [submission] User Submission is not a draft")
		c.JSON(http.StatusBadRequest, gin.H{"error": "User Submission is not a draft"})
		return
	}

	if err != nil {
		slog.Error("[endpoints] [submission] Failed to update entry changelog", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update entry changelog"})
		return
	}

	userSubmission.State = entry.PendingReview

	err = db.Save(&userSubmission).Error

	if err != nil {
		slog.Error("[endpoints] [submission] Failed to update user submission", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user submission"})
		return
	}

	c.Status(http.StatusOK)
}

func redraftSubmission(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")

	exists, err := entry.GetEntryExists(db, accession)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "entry not found"})
		return
	}

	var userSubmission entry.UserSubmission

	err = db.
		Joins("JOIN entries ON entries.accession = user_submissions.entry_accession").
		Where("entries.accession = ?", accession).
		First(&userSubmission).Error

	if err != nil {
		slog.Error("[endpoints] [submission] Failed to find user submission", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user submission"})
		return
	}

	if userSubmission.State != entry.PendingReview {
		slog.Error("[endpoints] [submission] User Submission is not pending review")
		c.JSON(http.StatusBadRequest, gin.H{"error": "User Submission is not pending review"})
		return
	}

	userSubmission.State = entry.Draft

	err = db.Save(&userSubmission).Error

	if err != nil {
		slog.Error("[endpoints] [submission] Failed to update user submission", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user submission"})
		return
	}

	c.Status(http.StatusOK)
}

func discardSubmission(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")

	exists, err := entry.GetEntryExists(db, accession)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "entry not found"})
		return
	}

	var userSubmission entry.UserSubmission

	err = db.
		Joins("JOIN entries ON entries.accession = user_submissions.entry_accession").
		Where("entries.accession = ?", accession).
		First(&userSubmission).Error

	if err != nil {
		slog.Error("[endpoints] [submission] Failed to find user submission", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user submission"})
		return
	}

	if userSubmission.State != entry.Draft {
		slog.Error("[endpoints] [submission] User Submission is not a draft")
		c.JSON(http.StatusBadRequest, gin.H{"error": "User Submission is not a draft"})
		return
	}

	userSubmission.State = entry.Discarded

	err = db.Save(&userSubmission).Error

	if err != nil {
		slog.Error("[endpoints] [submission] Failed to update user submission", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user submission"})
		return
	}

	c.Status(http.StatusOK)
}
