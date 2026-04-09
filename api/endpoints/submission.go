package endpoints

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/models/entry"
	"github.com/adraismawur/mibig-submission/models/entry/biosynthesis"
	"github.com/adraismawur/mibig-submission/models/entry/taxonomy"
	"github.com/adraismawur/mibig-submission/models/lock"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func init() {
	RegisterEndpointGenerator(SubmissionEndpoint)
}

type ExistingSubmissionSubState string

const (
	Incomplete    ExistingSubmissionSubState = "incomplete"
	Unlocked                                 = "unlocked"
	Locked                                   = "locked"
	PendingReview                            = "pending review"
	BeingReviewed                            = "being reviewed"
	Accepted                                 = "accepted"
)

type SubmissionInfo struct {
	Accession       string               `json:"accession"`
	Type            entry.SubmissionType `json:"type"`
	SourceAccession string               `json:"source_accession"`
	Category        lock.LockingCategory `json:"category"`
}

// SubmissionEndpoint returns the submission endpoint.
// This endpoint will implement creating and updating submissions, as well as perform some
// specific checks on submissions.
func SubmissionEndpoint(db *gorm.DB) Endpoint {
	return Endpoint{
		Routes: []Route{
			{
				Method: "GET",
				Path:   "/submission/",
				Handler: func(c *gin.Context) {
					getSubmissions(db, c)
				},
			},
			{
				Method: "GET",
				Path:   "/reviews/active",
				Handler: func(c *gin.Context) {
					getActiveReviews(db, c)
				},
			},
			{
				Method: "GET",
				Path:   "/reviews/pending",
				Handler: func(c *gin.Context) {
					getPendingReviews(db, c)
				},
			},
			{
				Method: "GET",
				Path:   "/reviews/:accession",
				Handler: func(c *gin.Context) {
					getSubmissionReviews(db, c)
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
				Method: "GET",
				Path:   "/mutation/:accession",
				Handler: func(c *gin.Context) {
					getExistingMutations(db, c)
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
				Path:   "/submission/promote",
				Handler: func(c *gin.Context) {
					promoteSubmission(db, c)
				},
			},
			{
				Method: "POST",
				Path:   "/submission/claim_review",
				Handler: func(c *gin.Context) {
					claimReview(db, c)
				},
			},
			{
				Method: "POST",
				Path:   "/submission/cancel_review",
				Handler: func(c *gin.Context) {
					cancelReview(db, c)
				},
			},
			{
				Method: "POST",
				Path:   "/submission/accept",
				Handler: func(c *gin.Context) {
					acceptSubmission(db, c)
				},
			},
			{
				Method: "POST",
				Path:   "/submission/rfc",
				Handler: func(c *gin.Context) {
					requestSubmissionChanges(db, c)
				},
			},
			{
				Method: "POST",
				Path:   "/submission/redraft",
				Handler: func(c *gin.Context) {
					redraftSubmission(db, c)
				},
			},
			//{
			//	Method: "POST",
			//	Path:   "/submission/discard/:accession",
			//	Handler: func(c *gin.Context) {
			//		discardSubmission(db, c)
			//	},
			//},
		},
	}
}

func getSubmissionReviews(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")

	var reviews []entry.SubmissionReview

	err := db.
		Table("submission_reviews").
		Where("accession = $1", accession).
		Find(&reviews).
		Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reviews)
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
	user, err := models.GetUserFromContext(c)

	if err != nil {
		return
	}

	start, err := strconv.Atoi(c.Query("start"))

	if err != nil {
		start = 0
	}

	limit, err := strconv.Atoi(c.Query("limit"))

	if err != nil {
		limit = 20
	}

	search := c.Query("search")

	type ExistingSubmissionSubStateSummary struct {
		Locitax         ExistingSubmissionSubState `json:"locitax"`
		Biosynth        ExistingSubmissionSubState `json:"biosynth"`
		Compounds       ExistingSubmissionSubState `json:"compounds"`
		GeneInformation ExistingSubmissionSubState `json:"gene_information"`
		Finalize        ExistingSubmissionSubState `json:"finalize"`
	}

	type ExistingSubmissionSummary struct {
		EntryAccession  string               `json:"accession"`
		Type            entry.SubmissionType `json:"type"`
		SourceAccession string               `json:"source_accession"`
		Owner           bool                 `json:"owner"`
	}
	var submissions []ExistingSubmissionSummary

	userID := c.Query("id")
	state := c.Query("state")

	var locitaxState string
	var biosynthState string
	var compoundState string
	var geneInformationState string
	var finalState string

	if locitaxState, err = url.QueryUnescape(c.Query("locitax")); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to unescape query param"})
		return
	}
	if biosynthState, err = url.QueryUnescape(c.Query("biosynth")); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to unescape query param"})
		return
	}
	if compoundState, err = url.QueryUnescape(c.Query("compounds")); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to unescape query param"})
		return
	}
	if geneInformationState, err = url.QueryUnescape(c.Query("gene_information")); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to unescape query param"})
		return
	}
	if finalState, err = url.QueryUnescape(c.Query("finalize")); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to unescape query param"})
		return
	}

	now := time.Now().Unix()

	q := db.Session(&gorm.Session{})

	q = q.Table("user_submissions")

	clauseIdx := 1

	// optional clause
	if userID != "" {
		q.Where(fmt.Sprintf("user_submissions.user_id = $%d", clauseIdx), userID)
		clauseIdx += 1
	}

	if state != "" {
		q.Where(fmt.Sprintf("user_submissions.state = $%d", clauseIdx), state)
		clauseIdx += 1
	}

	if search != "" {
		q.Where(fmt.Sprintf("user_submissions.entry_accession LIKE $%d", clauseIdx), "%"+search+"%")
		clauseIdx += 1
	}

	// this process is annoying, probably a result of flawed design of this entire system. the state of an entry/category
	// should probably just be tracked in one table. inferring it through what locks/reviews are present is not a bad idea
	// but it is much more complex and one ought to keep things simple
	// TODO: keep things simple. rework this
	type StateFilter struct {
		Category entry.Category
		State    ExistingSubmissionSubState
	}

	stateFilters := make([]StateFilter, 0)

	if locitaxState != "" {
		stateFilters = append(stateFilters, StateFilter{
			Category: entry.Locitax,
			State:    ExistingSubmissionSubState(locitaxState),
		})
	}

	if biosynthState != "" {
		stateFilters = append(stateFilters, StateFilter{
			Category: entry.Biosynth,
			State:    ExistingSubmissionSubState(biosynthState),
		})
	}
	if compoundState != "" {
		stateFilters = append(stateFilters, StateFilter{
			Category: entry.Compounds,
			State:    ExistingSubmissionSubState(compoundState),
		})
	}
	if geneInformationState != "" {
		stateFilters = append(stateFilters, StateFilter{
			Category: entry.Genes,
			State:    ExistingSubmissionSubState(geneInformationState),
		})
	}
	if finalState != "" {
		stateFilters = append(stateFilters, StateFilter{
			Category: entry.Final,
			State:    ExistingSubmissionSubState(finalState),
		})
	}

	for _, stateFilter := range stateFilters {
		switch stateFilter.State {
		case Unlocked:
			q.Where(
				fmt.Sprintf(
					"user_submissions.entry_accession NOT IN ("+
						"SELECT entry_accession FROM locks WHERE "+
						"(category = $%d OR category = 'full') "+
						"OR unlocks_at <= $%d"+
						") AND user_submissions.entry_accession NOT IN ("+
						"SELECT accession FROM submission_reviews WHERE "+
						"category = $%d)",
					clauseIdx,
					clauseIdx+1,
					clauseIdx+2,
				),
				stateFilter.Category,
				now,
				stateFilter.Category,
			)
			clauseIdx += 3
			break
		case Locked:
			q.Where(
				fmt.Sprintf(
					"user_submissions.entry_accession IN ("+
						"SELECT entry_accession FROM locks WHERE "+
						"(category = $%d OR category = 'full') "+
						"AND unlocks_at > $%d"+
						") AND user_submissions.entry_accession NOT IN ("+
						"SELECT accession FROM submission_reviews WHERE "+
						"state = 'being reviewed' OR state = 'accepted'"+
						")",
					clauseIdx,
					clauseIdx+1,
				),
				stateFilter.Category,
				now,
			)
			clauseIdx += 2
			break
		case PendingReview:
			fallthrough
		case BeingReviewed:
			fallthrough
		case Accepted:
			q.Where(
				fmt.Sprintf(
					"user_submissions.entry_accession IN (SELECT accession FROM submission_reviews WHERE category = $%d AND state = $%d)",
					clauseIdx,
					clauseIdx+1,
				),
				stateFilter.Category,
				stateFilter.State,
			)
			clauseIdx += 2
			break
		}
	}

	var recordCount int64
	q.
		Select(fmt.Sprintf("user_submissions.entry_accession, user_submissions.type, user_submissions.source_accession, user_submissions.user_id = $%d as owner", clauseIdx), user.ID).
		Count(&recordCount).
		Offset(start).
		Limit(limit).
		Order("owner desc").
		Find(&submissions)

	err = q.Error

	if err != nil {
		slog.Error("[endpoints] [submission] Could not find submissions", "user_id", userID, "error", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	type ResponseSubmission struct {
		ExistingSubmissionSummary
		Taxonomy  string                            `json:"taxonomy"`
		Class     []string                          `json:"class"`
		SubStates ExistingSubmissionSubStateSummary `json:"sub_states"`
	}

	var response struct {
		Submissions []ResponseSubmission `json:"submissions"`
		RecordCount int64                `json:"record_count"`
	}

	response.RecordCount = recordCount
	response.Submissions = make([]ResponseSubmission, 0)

	var accessions []string
	submissionMap := make(map[string]int)

	for i, submission := range submissions {
		accessions = append(accessions, submission.EntryAccession)
		submissionMap[submission.EntryAccession] = i
		response.Submissions = append(response.Submissions, ResponseSubmission{
			ExistingSubmissionSummary: submission,
			Taxonomy:                  "",
			Class:                     []string{""},
			SubStates: ExistingSubmissionSubStateSummary{
				Locitax:         Unlocked,
				Biosynth:        Unlocked,
				Compounds:       Unlocked,
				GeneInformation: Unlocked,
				Finalize:        Unlocked,
			},
		})
	}

	var taxonomies []taxonomy.Taxonomy
	err = db.Table("taxonomies").
		Where("entry_accession IN ?", accessions).
		Find(&taxonomies).
		Error

	for _, taxonomy := range taxonomies {
		response.Submissions[submissionMap[taxonomy.EntryAccession]].Taxonomy = taxonomy.Name
	}

	var biosyntheses []biosynthesis.Biosynthesis
	err = db.Table("biosyntheses").
		Where("entry_accession IN ?", accessions).
		Preload("Classes").
		Find(&biosyntheses).
		Error

	for _, biosynth := range biosyntheses {
		var classes []string = make([]string, 0)
		for _, class := range biosynth.Classes {
			classes = append(classes, class.Class)
		}
		response.Submissions[submissionMap[biosynth.EntryAccession]].Class = classes
	}

	var locks []lock.Lock
	err = db.Table("locks").
		Where("entry_accession IN ? AND unlocks_at > ?", accessions, now).
		Find(&locks).
		Error

	for _, submissionLock := range locks {
		submission := response.Submissions[submissionMap[submissionLock.EntryAccession]]
		switch submissionLock.Category {
		case lock.Full:
			submission.SubStates.Locitax = Locked
			submission.SubStates.Biosynth = Locked
			submission.SubStates.Compounds = Locked
			submission.SubStates.GeneInformation = Locked
			submission.SubStates.Finalize = Locked
			break
		case lock.Locitax:
			submission.SubStates.Locitax = Locked
			break
		case lock.Biosynth:
			submission.SubStates.Biosynth = Locked
			break
		case lock.Compounds:
			submission.SubStates.Compounds = Locked
			break
		case lock.GeneInformation:
			submission.SubStates.GeneInformation = Locked
			break
		case lock.Final:
			submission.SubStates.Finalize = Locked
		}
		response.Submissions[submissionMap[submissionLock.EntryAccession]] = submission
	}

	var reviews []entry.SubmissionReview
	err = db.Table("submission_reviews").
		Where("accession IN ?", accessions).
		Find(&reviews).
		Error

	for _, review := range reviews {
		submission := response.Submissions[submissionMap[review.Accession]]

		switch review.Category {
		case entry.Locitax:
			submission.SubStates.Locitax = ExistingSubmissionSubState(review.State)
			break
		case entry.Biosynth:
			submission.SubStates.Biosynth = ExistingSubmissionSubState(review.State)
			break
		case entry.Compounds:
			submission.SubStates.Compounds = ExistingSubmissionSubState(review.State)
			break
		case entry.Genes:
			submission.SubStates.GeneInformation = ExistingSubmissionSubState(review.State)
			break
		case entry.Final:
			submission.SubStates.Finalize = ExistingSubmissionSubState(review.State)
		}

		response.Submissions[submissionMap[review.Accession]] = submission
	}

	c.JSON(http.StatusOK, response)
}

func getPendingReviews(db *gorm.DB, c *gin.Context) {

	type ReviewInfo struct {
		entry.SubmissionReview
		Type            string `json:"type"`
		SourceAccession string `json:"source_accession"`
	}

	var reviews []ReviewInfo

	user, err := models.GetUserFromContext(c)

	if !models.GetIsUserRole(user, models.Reviewer) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user is not a reviewer"})
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = db.Table("submission_reviews").
		Select("submission_reviews.*, user_submissions.type, user_submissions.source_accession").
		Joins("JOIN user_submissions ON user_submissions.entry_accession = submission_reviews.accession").
		Where("submission_reviews.state = $1", entry.PendingReview, user.ID).
		Find(&reviews).
		Error

	if err != nil {
		slog.Error("[endpoints] [submission] Could not find reviews", "user_id", user.ID, "error", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, reviews)
}

func getActiveReviews(db *gorm.DB, c *gin.Context) {

	type ReviewInfo struct {
		entry.SubmissionReview
		Type            string `json:"type"`
		SourceAccession string `json:"source_accession"`
	}

	var reviews []ReviewInfo

	user, err := models.GetUserFromContext(c)

	if !models.GetIsUserRole(user, models.Reviewer) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user is not a reviewer"})
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = db.Table("submission_reviews").
		Select("submission_reviews.*, user_submissions.type, user_submissions.source_accession").
		Joins("JOIN user_submissions ON user_submissions.entry_accession = submission_reviews.accession").
		Where("submission_reviews.state = $1 AND submission_reviews.user_id = $2", entry.Reviewing, user.ID).
		Find(&reviews).
		Error

	if err != nil {
		slog.Error("[endpoints] [submission] Could not find reviews", "user_id", user.ID, "error", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, reviews)
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

	var minimalEntryStart int
	var minimalEntryEnd int

	if minimalEntry.Locus.Location.Start == nil {
		minimalEntryStart = -1
	} else {
		minimalEntryStart = int(*minimalEntry.Locus.Location.Start)
	}

	if minimalEntry.Locus.Location.End == nil {
		minimalEntryEnd = -1
	} else {
		minimalEntryEnd = int(*minimalEntry.Locus.Location.End)
	}

	// check if there is an existing entry with the accession / range in the database
	existingAccession, err := entry.GetEntryAccessionByLociAccession(
		db,
		minimalEntry.Locus.Accession,
		minimalEntryStart,
		minimalEntryEnd,
	)

	if err != nil {
		slog.Error("[endpoints] [submission] Error getting existing submission from locus", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user submission record"})
		return
	}

	if existingAccession != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "locus matches existing entry", "accession": existingAccession})
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

func getExistingMutations(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")

	existingMutations := make([]entry.UserSubmission, 0)

	if accession == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Missing accession"})
		return
	}

	err := db.Table("user_submissions").
		Where("source_accession = $1", accession).
		Find(&existingMutations).
		Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to find existing user mutations"})
		return
	}

	c.JSON(http.StatusOK, existingMutations)
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
	var promotionRequest struct {
		Accession string
		Category  entry.Category
		Notes     string
	}

	err := c.ShouldBindJSON(&promotionRequest)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not bind json: " + err.Error()})
		return
	}

	exists, err := entry.GetEntryExists(db, promotionRequest.Accession)

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
		Where("user_submissions.entry_accession = $1", promotionRequest.Accession).
		First(&userSubmission).Error

	if err != nil {
		slog.Error("[endpoints] [submission] Failed to find user submission", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user submission"})
		return
	}

	newReview := entry.SubmissionReview{
		Accession:      promotionRequest.Accession,
		Category:       promotionRequest.Category,
		State:          entry.PendingReview,
		SubmitterNotes: promotionRequest.Notes,
	}

	err = db.
		Omit("User").
		Save(&newReview).
		Error

	if err != nil {
		slog.Error("[endpoints] [submission] Failed to update user submission", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user submission: " + err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func claimReview(db *gorm.DB, c *gin.Context) {
	var err error

	var reviewClaimRequest struct {
		Accession string
		Category  entry.Category
	}

	err = c.ShouldBindJSON(&reviewClaimRequest)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not bind review claim request json: " + err.Error()})
		return
	}

	user, err := models.GetUserFromContext(c)

	if !models.GetIsUserRole(user, models.Reviewer) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "requester is not a reviewer"})
		return
	}

	err = db.Session(&gorm.Session{FullSaveAssociations: true}).Transaction(func(tx *gorm.DB) error {

		exists, transactionErr := entry.GetEntryExists(tx, reviewClaimRequest.Accession)

		if transactionErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": transactionErr.Error()})
			return transactionErr
		}

		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "entry not found"})
			return transactionErr
		}

		var submissionReview entry.SubmissionReview

		transactionErr = tx.
			Where("accession = $1 AND category = $2", reviewClaimRequest.Accession, reviewClaimRequest.Category).
			First(&submissionReview).
			Error

		if transactionErr != nil {
			slog.Error("[endpoints] [submission] Failed to find submission review", "error", transactionErr.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find submission review"})
			return transactionErr
		}

		if submissionReview.State != entry.PendingReview {
			slog.Error("[endpoints] [submission] review is not pending")
			c.JSON(http.StatusBadRequest, gin.H{"error": "review is not pending"})
			return transactionErr
		}

		submissionReview.State = entry.Reviewing
		submissionReview.UserID = user.ID

		transactionErr = tx.Save(&submissionReview).Error

		if transactionErr != nil {
			slog.Error("[endpoints] [submission] Failed to update review", "error", transactionErr.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update review"})
			return transactionErr
		}

		_, transactionErr = lock.CreateOrGetLock(tx, reviewClaimRequest.Accession, lock.LockingCategory(submissionReview.Category), *user)

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

	entry.AddReviewer(db, reviewClaimRequest.Accession, user.ID)

	c.Status(http.StatusOK)
}

func cancelReview(db *gorm.DB, c *gin.Context) {
	var err error

	var cancelReviewRequest struct {
		Accession string         `json:"accession"`
		Category  entry.Category `json:"category"`
	}

	err = c.ShouldBindJSON(&cancelReviewRequest)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not bind review cancel request json: " + err.Error()})
		return
	}

	user, err := models.GetUserFromContext(c)

	if !models.GetIsUserRole(user, models.Reviewer) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "requester is not a reviewer"})
		return
	}

	err = db.Session(&gorm.Session{FullSaveAssociations: true}).Transaction(func(tx *gorm.DB) error {
		exists, transactionErr := entry.GetEntryExists(tx, cancelReviewRequest.Accession)

		if transactionErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return transactionErr
		}

		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "entry not found"})
			return transactionErr
		}

		var submissionReview entry.SubmissionReview

		transactionErr = tx.
			Where("accession = $1 and category = $2", cancelReviewRequest.Accession, cancelReviewRequest.Category).
			First(&submissionReview).Error

		if transactionErr != nil {
			slog.Error("[endpoints] [submission] Failed to find review", "error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find review"})
			return transactionErr
		}

		if submissionReview.State != entry.Reviewing {
			slog.Error("[endpoints] [submission] review is not being reviewed")
			c.JSON(http.StatusBadRequest, gin.H{"error": "review is not being reviewed"})
			return transactionErr
		}

		submissionReview.State = entry.PendingReview

		transactionErr = tx.Save(&submissionReview).Error

		if transactionErr != nil {
			slog.Error("[endpoints] [submission] Failed to update review", "error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update review"})
			return transactionErr
		}

		transactionErr = lock.ReleaseLock(tx, cancelReviewRequest.Accession, lock.LockingCategory(cancelReviewRequest.Category), *user)

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

	var acceptRequest struct {
		Accession string
		Category  lock.LockingCategory
	}

	err = c.ShouldBindJSON(&acceptRequest)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not bind acceptation request json: " + err.Error()})
		return
	}

	user, err := models.GetUserFromContext(c)

	if !models.GetIsUserRole(user, models.Reviewer) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "requester is not a reviewer"})
		return
	}

	err = db.Session(&gorm.Session{FullSaveAssociations: true}).Transaction(func(tx *gorm.DB) error {
		exists, transactionErr := entry.GetEntryExists(tx, acceptRequest.Accession)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return transactionErr
		}

		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "entry not found"})
			return transactionErr
		}

		var submissionReview entry.SubmissionReview

		transactionErr = tx.
			Where("accession = $1 AND category = $2", acceptRequest.Accession, acceptRequest.Category).
			First(&submissionReview).Error

		if transactionErr != nil {
			slog.Error("[endpoints] [submission] Failed to find review", "error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find review"})
			return transactionErr
		}

		if submissionReview.State != entry.Reviewing {
			slog.Error("[endpoints] [submission] review is not being reviewed")
			c.JSON(http.StatusBadRequest, gin.H{"error": "review is not being reviewed"})
			return transactionErr
		}

		submissionReview.State = entry.Accepted

		transactionErr = tx.Save(&submissionReview).Error

		if transactionErr != nil {
			slog.Error("[endpoints] [submission] Failed to update review", "error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update review"})
			return transactionErr
		}

		transactionErr = lock.ReleaseLock(tx, acceptRequest.Accession, acceptRequest.Category, *user)

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
	var err error

	var redraftRequest struct {
		Accession string
		Category  lock.LockingCategory
	}

	err = c.ShouldBindJSON(&redraftRequest)

	exists, err := entry.GetEntryExists(db, redraftRequest.Accession)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "entry not found"})
		return
	}

	var submissionReview entry.SubmissionReview

	err = db.
		Where("accession = $1 AND category = $2", redraftRequest.Accession, redraftRequest.Category).
		First(&submissionReview).Error

	if err != nil {
		slog.Error("[endpoints] [submission] Failed to find user submission", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user submission"})
		return
	}

	if submissionReview.State != entry.PendingReview {
		slog.Error("[endpoints] [submission] Review is not pending")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Review is not pending"})
		return
	}

	err = db.Delete(&submissionReview).Error

	if err != nil {
		slog.Error("[endpoints] [submission] Failed to update review", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update review"})
		return
	}

	c.Status(http.StatusOK)
}

//
//func discardSubmission(db *gorm.DB, c *gin.Context) {
//	accession := c.Param("accession")
//
//	exists, err := entry.GetEntryExists(db, accession)
//
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	if !exists {
//		c.JSON(http.StatusNotFound, gin.H{"error": "entry not found"})
//		return
//	}
//
//	var userSubmission entry.UserSubmission
//
//	err = db.
//		Joins("JOIN entries ON entries.accession = user_submissions.entry_accession").
//		Where("entries.accession = $1", accession).
//		First(&userSubmission).Error
//
//	if err != nil {
//		slog.Error("[endpoints] [submission] Failed to find user submission", "error", err.Error())
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user submission"})
//		return
//	}
//
//	if userSubmission.State != entry.Draft {
//		slog.Error("[endpoints] [submission] User Submission is not a draft")
//		c.JSON(http.StatusBadRequest, gin.H{"error": "User Submission is not a draft"})
//		return
//	}
//
//	userSubmission.State = entry.Discarded
//
//	err = db.Save(&userSubmission).Error
//
//	if err != nil {
//		slog.Error("[endpoints] [submission] Failed to update user submission", "error", err.Error())
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user submission"})
//		return
//	}
//
//	c.Status(http.StatusOK)
//}
