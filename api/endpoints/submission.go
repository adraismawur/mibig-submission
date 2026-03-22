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
				Method: "GET",
				Path:   "/reviews",
				Handler: func(c *gin.Context) {
					getReviews(db, c)
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
				Path:   "/submission/claim_review/:accession",
				Handler: func(c *gin.Context) {
					claimReview(db, c)
				},
			},
			{
				Method: "POST",
				Path:   "/submission/cancel_review/:accession",
				Handler: func(c *gin.Context) {
					cancelReview(db, c)
				},
			},
			{
				Method: "POST",
				Path:   "/submission/accept/:accession",
				Handler: func(c *gin.Context) {
					acceptSubmission(db, c)
				},
			},
			{
				Method: "POST",
				Path:   "/submission/rfc/:accession",
				Handler: func(c *gin.Context) {
					requestSubmissionChanges(db, c)
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

	q.Where("user_submissions.user_id = $1", userID)

	q.Where("state != $2", entry.Discarded)

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
		q.Where("user_submissions.user_id = $1", userID)
	}

	if state != "" {
		q.Where("user_submissions.state = $1", state)
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

func getReviews(db *gorm.DB, c *gin.Context) {
	var submissions []SubmissionInfo

	user, err := models.GetUserFromContext(c)

	if !models.GetIsUserRole(user, models.Reviewer) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user is not a reviewer"})
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = db.Table("user_submissions").
		Joins("JOIN entries ON entries.accession = user_submissions.entry_accession").
		Joins("JOIN submission_reviewers ON submission_reviewers.accession = user_submissions.entry_accession").
		Where("user_submissions.state = $1", entry.Reviewing).
		Where("submission_reviewers.user_id = $2", user.ID).
		Select("entries.accession, user_submissions.type, user_submissions.source_accession, user_submissions.state").
		Find(&submissions).
		Error

	if err != nil {
		slog.Error("[endpoints] [submission] Could not find reviews", "user_id", user.ID, "error", err)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		Where("user_id = $1", userId).
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
		Where("entries.accession = $1", accession).
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

func claimReview(db *gorm.DB, c *gin.Context) {
	var err error
	accession := c.Param("accession")

	user, err := models.GetUserFromContext(c)

	if !models.GetIsUserRole(user, models.Reviewer) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "requester is not a reviewer"})
		return
	}

	err = db.Session(&gorm.Session{FullSaveAssociations: true}).Transaction(func(tx *gorm.DB) error {

		exists, transactionErr := entry.GetEntryExists(tx, accession)

		if transactionErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": transactionErr.Error()})
			return transactionErr
		}

		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "entry not found"})
			return transactionErr
		}

		var userSubmission entry.UserSubmission

		transactionErr = tx.
			Joins("JOIN entries ON entries.accession = user_submissions.entry_accession").
			Where("entries.accession = $1", accession).
			First(&userSubmission).Error

		if transactionErr != nil {
			slog.Error("[endpoints] [submission] Failed to find user submission", "error", transactionErr.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user submission"})
			return transactionErr
		}

		if userSubmission.State != entry.PendingReview {
			slog.Error("[endpoints] [submission] User Submission is not pending review")
			c.JSON(http.StatusBadRequest, gin.H{"error": "User Submission is not pending review"})
			return transactionErr
		}

		userSubmission.State = entry.Reviewing

		transactionErr = tx.Save(&userSubmission).Error

		if transactionErr != nil {
			slog.Error("[endpoints] [submission] Failed to update user submission", "error", transactionErr.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user submission"})
			return transactionErr
		}

		// finally create the reviewer
		reviewer := entry.SubmissionReviewer{
			User:      *user,
			UserID:    user.ID,
			Accession: accession,
		}

		transactionErr = tx.
			Omit("User").
			Create(&reviewer).
			Error

		if transactionErr != nil {
			slog.Error("[endpoints] [submission] Failed to create submission reviewer", "error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create submission reviewer"})
			return transactionErr
		}

		_, transactionErr = lock.CreateOrGetLock(tx, accession, lock.Full, *user)

		if transactionErr != nil {
			slog.Error("[endpoints] [submission] Failed to create submission locks", "error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create submission locks"})
			return transactionErr
		}

		return nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func cancelReview(db *gorm.DB, c *gin.Context) {
	var err error
	accession := c.Param("accession")

	user, err := models.GetUserFromContext(c)

	if !models.GetIsUserRole(user, models.Reviewer) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "requester is not a reviewer"})
		return
	}

	err = db.Session(&gorm.Session{FullSaveAssociations: true}).Transaction(func(tx *gorm.DB) error {
		exists, transactionErr := entry.GetEntryExists(tx, accession)

		if transactionErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return transactionErr
		}

		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "entry not found"})
			return transactionErr
		}

		var userSubmission entry.UserSubmission

		transactionErr = tx.
			Joins("JOIN entries ON entries.accession = user_submissions.entry_accession").
			Where("entries.accession = $1", accession).
			First(&userSubmission).Error

		if transactionErr != nil {
			slog.Error("[endpoints] [submission] Failed to find user submission", "error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user submission"})
			return transactionErr
		}

		if userSubmission.State != entry.Reviewing {
			slog.Error("[endpoints] [submission] User Submission is not being reviewed")
			c.JSON(http.StatusBadRequest, gin.H{"error": "User Submission is not being reviewed"})
			return transactionErr
		}

		userSubmission.State = entry.PendingReview

		transactionErr = tx.Save(&userSubmission).Error

		if transactionErr != nil {
			slog.Error("[endpoints] [submission] Failed to update user submission", "error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user submission"})
			return transactionErr
		}

		transactionErr = tx.
			Model(&entry.SubmissionReviewer{}).
			Where("accession = $1", accession).
			Delete(&entry.SubmissionReviewer{}).
			Error

		if transactionErr != nil {
			slog.Error("[endpoints] [submission] Failed to delete submission reviewer", "error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete submission reviewer"})
			return transactionErr
		}

		transactionErr = lock.ClearLocks(tx, accession, *user)

		if transactionErr != nil {
			slog.Error("[endpoints] [submission] Failed to clear submission locks", "error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear submission locks"})
			return transactionErr
		}

		return nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func acceptSubmission(db *gorm.DB, c *gin.Context) {
	var err error
	accession := c.Param("accession")

	user, err := models.GetUserFromContext(c)

	if !models.GetIsUserRole(user, models.Reviewer) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "requester is not a reviewer"})
		return
	}

	err = db.Session(&gorm.Session{FullSaveAssociations: true}).Transaction(func(tx *gorm.DB) error {
		exists, transactionErr := entry.GetEntryExists(tx, accession)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return transactionErr
		}

		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "entry not found"})
			return transactionErr
		}

		var userSubmission entry.UserSubmission

		transactionErr = tx.
			Joins("JOIN entries ON entries.accession = user_submissions.entry_accession").
			Where("entries.accession = $1", accession).
			First(&userSubmission).Error

		if transactionErr != nil {
			slog.Error("[endpoints] [submission] Failed to find user submission", "error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user submission"})
			return transactionErr
		}

		if userSubmission.State != entry.Reviewing {
			slog.Error("[endpoints] [submission] User Submission is not being reviewed")
			c.JSON(http.StatusBadRequest, gin.H{"error": "User Submission is not being reviewed"})
			return transactionErr
		}

		userSubmission.State = entry.Accepted

		transactionErr = tx.Save(&userSubmission).Error

		if transactionErr != nil {
			slog.Error("[endpoints] [submission] Failed to update user submission", "error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user submission"})
			return transactionErr
		}

		transactionErr = lock.ClearLocks(tx, accession, *user)

		if transactionErr != nil {
			slog.Error("[endpoints] [submission] Failed to clear submission locks", "error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear submission locks"})
			return transactionErr
		}

		return nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func requestSubmissionChanges(db *gorm.DB, c *gin.Context) {

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
		Where("entries.accession = $1", accession).
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
		Where("entries.accession = $1", accession).
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
