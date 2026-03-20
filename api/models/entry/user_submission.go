package entry

import (
	"fmt"
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/models/entry/export"
	"github.com/adraismawur/mibig-submission/util/constants"
	"github.com/beevik/guid"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
	"log/slog"
	"strconv"
	"time"
)

type SubmissionState string

const (
	Draft         SubmissionState = "draft"
	PendingReview                 = "pending review"
	Reviewing                     = "being reviewed"
	RFC                           = "requested changes"
	Accepted                      = "accepted"
	Discarded                     = "discarded"
)

type SubmissionType string

const (
	NewSubmission SubmissionType = "new submission"
	Mutation      SubmissionType = "entry mutation"
)

type UserSubmission struct {
	ID              uint64          `json:"db_id"`
	EntryAccession  string          `json:"db_submission_accession"`
	SourceAccession string          `json:"source_accession"`
	UserID          uint64          `json:"user_id"`
	Type            SubmissionType  `json:"type"`
	State           SubmissionState `json:"state"`
}

type SubmissionReviewer struct {
	ID        uint64      `json:"db_id"`
	Accession string      `json:"accession"`
	UserID    uint64      `json:"db_reviewer_id"`
	User      models.User `json:"reviewer"`
}

func init() {
	models.Models = append(models.Models, UserSubmission{})
	models.Models = append(models.Models, SubmissionReviewer{})
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

	newAccession, err := GeneratePlaceholderAccession(db, "new")

	newEntry.Accession = newAccession

	if err != nil {
		slog.Error("[endpoints] [submission] Failed to generate new entry accession", "error", err)
		return nil, err
	}

	db.Create(&newEntry)

	var userSubmission UserSubmission

	userSubmission.UserID = user.ID
	userSubmission.EntryAccession = newEntry.Accession
	userSubmission.SourceAccession = "new"
	userSubmission.Type = NewSubmission
	userSubmission.State = Draft

	err = db.Create(&userSubmission).Error

	return &newEntry, err
}

func CreateNewUserMutation(db *gorm.DB, accession string, user models.User) (*Entry, error) {

	var entryExport export.Entry

	existingEntry, err := GetEntryFromAccession(db, accession)

	if err != nil {
		return nil, err
	}

	var newEntry Entry

	mapstructure.Decode(existingEntry, &entryExport)
	mapstructure.Decode(&entryExport, &newEntry)

	mutAccession, err := GeneratePlaceholderAccession(db, "mut")

	newEntry.Accession = mutAccession

	err = db.Create(&newEntry).Error

	var userSubmission UserSubmission

	userSubmission.UserID = user.ID
	userSubmission.EntryAccession = newEntry.Accession
	userSubmission.SourceAccession = accession
	userSubmission.Type = Mutation
	userSubmission.State = Draft

	err = db.Create(&userSubmission).Error

	return &newEntry, err
}

func CreateAntismashWorkerTask(db *gorm.DB, newEntry Entry) (*models.AntismashRun, error) {
	antismashTask := models.AntismashRun{
		LocusAccession: newEntry.Loci[0].Accession,
		EntryAccession: newEntry.Accession,
		GUID:           guid.NewString(),
		SubmittedAt:    time.Now(),
	}

	err := db.Create(antismashTask).Error

	if err != nil {
		return nil, err
	}

	return &antismashTask, nil
}

const PlaceholderNumLen = 7

func GeneratePlaceholderAccession(db *gorm.DB, prefix string) (string, error) {
	// get count for if there are no submission yet
	var count int64

	// I am not happy about this either
	clause := fmt.Sprintf("accession like '%s%%'", prefix)

	err := db.Table("entries").
		Select("accession").
		Where(clause).
		Count(&count).
		Error

	if err != nil {
		return "", err
	}

	// if so just return new0000001
	if count == 0 {
		numPart := fmt.Sprintf("%0*d", PlaceholderNumLen, 1)

		return prefix + numPart, nil
	}

	// otherwise get the latest placeholder number
	var accession string

	err = db.Model(&Entry{}).
		Select("accession").
		Where(clause).
		Last(&accession).
		Error

	lastNum, err := strconv.Atoi(accession[3:])

	if err != nil {
		return "", err
	}

	numPart := fmt.Sprintf("%0*d", PlaceholderNumLen, lastNum+1)

	return prefix + numPart, nil
}
