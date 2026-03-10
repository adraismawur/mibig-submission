package entry

import (
	"fmt"
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/util/constants"
	"github.com/beevik/guid"
	"gorm.io/gorm"
	"log/slog"
	"strconv"
	"time"
)

type SubmissionState string

const (
	DraftSubmission    SubmissionState = "draft"
	PendingReview                      = "pending_review"
	AcceptedSubmission                 = "accepted"
)

type SubmissionType string

const (
	NewSubmission  SubmissionType = "new"
	SubmissionEdit SubmissionType = "edit"
)

type UserSubmission struct {
	ID      uint64          `json:"db_id"`
	EntryID uint64          `json:"submission_id"`
	UserID  uint64          `json:"user_id"`
	Type    SubmissionType  `json:"type"`
	State   SubmissionState `json:"state"`
}

func init() {
	models.Models = append(models.Models, UserSubmission{})
}

func CreateNewUserSubmission(db *gorm.DB, minimalEntry MinimalEntry, user models.User) (*Entry, error) {
	var newEntry Entry

	// first thing we do is add the one locus someone can submit
	newEntry.Loci = append(newEntry.Loci, minimalEntry.Locus)

	// then we copy over the compounds
	// TODO: validate these compounds
	newEntry.Compounds = minimalEntry.Compounds

	// generate a new changelog
	var currentDate = time.Now().Format(time.DateOnly)

	newEntry.Changelog = Changelog{
		Releases: []Release{
			{
				Version: "1",
				Date:    currentDate,
				Entries: []ReleaseEntry{
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

	newAccession, err := GeneratePlaceholderAccession(db)

	newEntry.Accession = newAccession

	if err != nil {
		slog.Error("[endpoints] [submission] Failed to generate new entry accession", "error", err)
		return nil, err
	}

	db.Create(&newEntry)

	var userSubmission UserSubmission

	userSubmission.UserID = user.ID
	userSubmission.EntryID = newEntry.ID
	userSubmission.State = DraftSubmission

	err = db.Create(&userSubmission).Error

	return &newEntry, err
}

func CreateAntismashWorkerTask(db *gorm.DB, newEntry Entry) (*models.AntismashRun, error) {
	antismashTask := models.AntismashRun{
		Accession:   newEntry.Loci[0].Accession,
		BGCID:       newEntry.Accession,
		GUID:        guid.NewString(),
		SubmittedAt: time.Now(),
	}

	err := db.Create(antismashTask).Error

	if err != nil {
		return nil, err
	}

	return &antismashTask, nil
}

const PlaceholderNumLen = 7

func GeneratePlaceholderAccession(db *gorm.DB) (string, error) {
	newPart := "new"

	// get count for if there are no submission yet
	var count int64

	err := db.Table("entries").
		Select("accession").
		Where("accession LIKE 'new%'").
		Count(&count).
		Error

	if err != nil {
		return "", err
	}

	// if so just return new0000001
	if count == 0 {
		numPart := fmt.Sprintf("%0*d", PlaceholderNumLen, 1)

		return newPart + numPart, nil
	}

	// otherwise get the latest placeholder number
	var accession string

	err = db.Model(&Entry{}).
		Select("accession").
		Where("accession LIKE 'new%'").
		Last(&accession).
		Error

	lastNum, err := strconv.Atoi(accession[3:])

	if err != nil {
		return "", err
	}

	numPart := fmt.Sprintf("%0*d", PlaceholderNumLen, lastNum+1)

	return newPart + numPart, nil
}
