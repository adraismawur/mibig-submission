package endpoints

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/models/entry"
	"github.com/adraismawur/mibig-submission/util/constants"
	"github.com/adraismawur/mibig-submission/util/entry_utils"
	"github.com/beevik/guid"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

func init() {
	RegisterEndpointGenerator(SubmissionEndpoint)
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
					createSubmission(db, c)
				},
			},
			{
				Method: "POST",
				Path:   "/submission/submit/:accession",
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
		},
	}
}

func getUserSubmissions(db *gorm.DB, c *gin.Context) {
	var submissions []struct {
		Accession string                `json:"accession"`
		State     entry.SubmissionState `json:"state"`
	}

	userID := c.Param("userId")

	q := db.Table("user_submissions").
		Joins("JOIN entries ON entries.id = user_submissions.entry_id")

	// optional clause
	if userID != "" {
		q.Where("user_submissions.user_id = ?", userID)
	}

	err := q.Select("entries.accession, user_submissions.state").
		Find(&submissions).Error

	if err != nil {
		slog.Error("[endpoints] [submission] Could not find submissions for user", "user_id", userID, "error", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, submissions)
}

func getSubmissions(db *gorm.DB, c *gin.Context) {
	var submissions []struct {
		Accession string                `json:"accession"`
		State     entry.SubmissionState `json:"state"`
	}

	userID := c.Query("id")
	state := c.Query("state")

	q := db.Table("user_submissions").
		Joins("JOIN entries ON entries.id = user_submissions.entry_id")

	// optional clause
	if userID != "" {
		q.Where("user_submissions.user_id = ?", userID)
	}

	if state != "" {
		q.Where("user_submissions.state = ?", state)
	}

	err := q.Select("entries.accession, user_submissions.state").
		Find(&submissions).Error

	if err != nil {
		slog.Error("[endpoints] [submission] Could not find submissions", "user_id", userID, "error", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, submissions)
}

// createSubmission creates a minimal draft submission from a request
func createSubmission(db *gorm.DB, c *gin.Context) {
	var newEntry entry.Entry

	if err := c.BindJSON(&newEntry); err != nil {
		slog.Error("[endpoints] [submission] Failed to unmarshal submission JSON", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid entry submitted"})
		return
	}

	var currentDate = time.Now().Format(time.DateOnly)

	newEntry.Changelog = entry.Changelog{
		Releases: []entry.Release{
			{
				Version: "1",
				Date:    currentDate,
				Entries: []entry.ReleaseEntry{
					{
						Contributors: []string{
							constants.AnonymousUserId,
						},
						Reviewers: nil,
						Date:      currentDate,
						Comment:   constants.NewEntryComment,
					},
				},
			},
		},
	}

	user, err := models.GetUserFromContext(c)

	if err != nil {
		slog.Error("[endpoints] [submission] Failed to generate new entry accession", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate new entry accession"})
		return
	}

	newEntry.Accession = entry_utils.GeneratePlaceholderAccession(*user)

	db.Create(&newEntry)

	var userSubmission entry.UserSubmission

	userSubmission.UserID = user.ID
	userSubmission.EntryID = newEntry.ID
	userSubmission.State = entry.DraftSubmission

	err = db.Create(&userSubmission).Error

	if err != nil {
		slog.Error("[endpoints] [submission] Failed to create user submission record", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user submission record"})
		return
	}

	antismashTask := models.AntismashRun{
		Accession:   newEntry.Loci[0].Accession,
		BGCID:       newEntry.Accession,
		GUID:        guid.NewString(),
		SubmittedAt: time.Now(),
	}

	err = db.Create(antismashTask).Error

	if err != nil {
		slog.Error("[endpoints] [entry] Failed to create antismash task", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create antismash task"})
	}

	c.JSON(http.StatusOK, gin.H{"status": antismashTask})
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
		Joins("JOIN entries ON entries.id = user_submissions.entry_id").
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

	targetEntry, err := entry.GetEntryFromAccession(db, accession)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var finalDetails entry.FinalDetailsRequest

	err = c.ShouldBindJSON(&finalDetails)

	if err != nil {
		slog.Error("[endpoints] [submission] Failed to unmarshal entry JSON", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entry submitted"})
		return
	}

	var userSubmission entry.UserSubmission

	err = db.
		Joins("JOIN entries ON entries.id = user_submissions.entry_id").
		Where("entries.accession = ?", accession).
		First(&userSubmission).Error

	if err != nil {
		slog.Error("[endpoints] [submission] Failed to find user submission", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user submission"})
		return
	}

	if userSubmission.State != entry.DraftSubmission {
		slog.Error("[endpoints] [submission] User Submission is not a draft")
		c.JSON(http.StatusBadRequest, gin.H{"error": "User Submission is not a draft"})
		return
	}

	// after this we're past all the checks and we can do stuff
	// update the changelog with the comment
	// we know this is a submission, so we can edit the one changelog entry that
	// exists
	targetEntry.Changelog.Releases[0].Entries[0].Comment = finalDetails.Comment

	err = db.
		Model(&targetEntry).
		Association("Changelog").
		Replace(&targetEntry.Changelog)

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
		Joins("JOIN entries ON entries.id = user_submissions.entry_id").
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

	userSubmission.State = entry.DraftSubmission

	err = db.Save(&userSubmission).Error

	if err != nil {
		slog.Error("[endpoints] [submission] Failed to update user submission", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user submission"})
		return
	}

	c.Status(http.StatusOK)
}
