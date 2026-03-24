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

type ReviewState string

const (
	Draft         ReviewState = "draft"
	PendingReview             = "pending review"
	Reviewing                 = "being reviewed"
	RFC                       = "requested changes"
	Accepted                  = "accepted"
	Discarded                 = "discarded"
)

type Category string

const (
	Locitax   Category = "locitax"
	Biosynth  Category = "biosynth"
	Compounds Category = "compounds"
	Genes     Category = "gene_information"
	Final     Category = "finalize"
)

type SubmissionType string

const (
	NewSubmission SubmissionType = "new submission"
	Mutation      SubmissionType = "entry mutation"
)

type UserSubmission struct {
	ID              uint64         `json:"db_id"`
	EntryAccession  string         `json:"db_submission_accession"`
	SourceAccession string         `json:"source_accession"`
	UserID          uint64         `json:"user_id"`
	Type            SubmissionType `json:"type"`
}

type SubmissionReview struct {
	ID             uint64      `json:"db_id"`
	Accession      string      `json:"accession" gorm:"uniqueIndex:compositeReviewIndex"`
	Category       Category    `json:"category" gorm:"uniqueIndex:compositeReviewIndex"`
	State          ReviewState `json:"state"`
	UserID         uint64      `json:"db_reviewer_id"`
	User           models.User `json:"reviewer"`
	SubmitterNotes string      `json:"submitter_notes"`
	ReviewerNotes  string      `json:"reviewer_notes"`
}

func init() {
	models.Models = append(models.Models, UserSubmission{})
	models.Models = append(models.Models, SubmissionReview{})
}

func CreateNewUserSubmission(db *gorm.DB, minimalEntry MinimalEntry, user models.User) (*Entry, error) {
	var newEntry Entry

	err := db.Transaction(func(tx *gorm.DB) error {
		var transactionErr error

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
								user.Info.Alias,
							},
							Reviewers: []string{},
							Date:      currentDate,
							Comment:   constants.NewEntryComment,
						},
					},
				},
			},
		}

		newAccession, transactionErr := GeneratePlaceholderAccession(db, "new")

		newEntry.Accession = newAccession

		if transactionErr != nil {
			slog.Error("[endpoints] [submission] Failed to generate new entry accession", "error", transactionErr)
			return transactionErr
		}

		tx.Create(&newEntry)

		var userSubmission UserSubmission

		userSubmission.UserID = user.ID
		userSubmission.EntryAccession = newEntry.Accession
		userSubmission.SourceAccession = "new"
		userSubmission.Type = NewSubmission

		transactionErr = tx.Create(&userSubmission).Error

		if transactionErr != nil {
			return transactionErr
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &newEntry, err
}

func CreateNewUserMutation(db *gorm.DB, accession string, user models.User) (*Entry, error) {
	var newEntry Entry

	err := db.Transaction(func(tx *gorm.DB) error {
		var entryExport export.Entry

		existingEntry, transactionErr := GetEntryFromAccession(tx, accession)

		if transactionErr != nil {
			return transactionErr
		}

		transactionErr = mapstructure.Decode(existingEntry, &entryExport)

		if transactionErr != nil {
			return transactionErr
		}

		transactionErr = mapstructure.Decode(&entryExport, &newEntry)

		if transactionErr != nil {
			return transactionErr
		}

		mutAccession, transactionErr := GeneratePlaceholderAccession(tx, "mut")

		if transactionErr != nil {
			return transactionErr
		}

		newEntry.Accession = mutAccession

		transactionErr = tx.Create(&newEntry).Error

		if transactionErr != nil {
			return transactionErr
		}

		var userSubmission UserSubmission

		userSubmission.UserID = user.ID
		userSubmission.EntryAccession = newEntry.Accession
		userSubmission.SourceAccession = accession
		userSubmission.Type = Mutation

		transactionErr = tx.Create(&userSubmission).Error

		if transactionErr != nil {
			return transactionErr
		}

		if transactionErr != nil {
			return transactionErr
		}

		return nil

	})

	if err != nil {
		return nil, err
	}

	return &newEntry, err
}

func CreateAntismashWorkerTask(db *gorm.DB, newEntry Entry) (*models.AntismashRun, error) {
	var start, stop int

	if newEntry.Loci[0].Location.Start == nil {
		start = 0
	} else {
		start = int(*newEntry.Loci[0].Location.Start)
	}

	if newEntry.Loci[0].Location.End == nil {
		stop = 0
	} else {
		stop = int(*newEntry.Loci[0].Location.End)
	}

	antismashTask := models.AntismashRun{
		LocusAccession: newEntry.Loci[0].Accession,
		EntryAccession: newEntry.Accession,
		Start:          start,
		Stop:           stop,
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
