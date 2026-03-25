package entry

import (
	"github.com/adraismawur/mibig-submission/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log/slog"
)
import "github.com/lib/pq"

type ReleaseEntry struct {
	ID           uint64         `json:"db_id"`
	ReleaseID    uint64         `json:"release_id"`
	Contributors pq.StringArray `json:"contributors" gorm:"type:text[]"`
	Reviewers    pq.StringArray `json:"reviewers" gorm:"type:text[]"`
	Date         string         `json:"date"`
	Comment      string         `json:"comment"`
}

type Release struct {
	ID          uint64         `json:"db_id"`
	ChangelogID uint64         `json:"changelog_id"`
	Version     string         `json:"version"`
	Entries     []ReleaseEntry `json:"entries" gorm:"foreignKey:ReleaseID"`
	Date        string         `json:"date"`
}

type Changelog struct {
	ID             uint64    `json:"db_id"`
	EntryAccession string    `json:"db_entry_accession"`
	Releases       []Release `json:"releases" gorm:"foreignKey:ChangelogID"`
}

type SubmissionContributor struct {
	ID             uint64 `json:"db_id"`
	EntryAccession string `json:"entry_accession" gorm:"uniqueIndex:compositeSubmissionContributorIndex"`
	UserId         uint64 `json:"user_id" gorm:"uniqueIndex:compositeSubmissionContributorIndex"`
}

type SubmissionReviewer struct {
	ID             uint64 `json:"db_id"`
	EntryAccession string `json:"entry_accession" gorm:"uniqueIndex:compositeSubmissionReviewerIndex"`
	UserId         uint64 `json:"user_id" gorm:"uniqueIndex:compositeSubmissionReviewerIndex"`
}

func init() {
	models.Models = append(models.Models, &Changelog{})
	models.Models = append(models.Models, &Release{})
	models.Models = append(models.Models, &ReleaseEntry{})
	models.Models = append(models.Models, &SubmissionContributor{})
	models.Models = append(models.Models, &SubmissionReviewer{})
}

func AddContributor(db *gorm.DB, accession string, contributor uint64) {
	submissionContributor := SubmissionContributor{
		EntryAccession: accession,
		UserId:         contributor,
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.
			Clauses(clause.Insert{Modifier: "OR IGNORE"}).
			Create(&submissionContributor).
			Error

		return err
	})

	if err != nil {
		slog.Error("[Changelog] Could not add contributor")
	}
}

func AddReviewer(db *gorm.DB, accession string, reviewer uint64) {
	submissionContributor := SubmissionReviewer{
		EntryAccession: accession,
		UserId:         reviewer,
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.
			Clauses(clause.Insert{Modifier: "OR IGNORE"}).
			Create(&submissionContributor).
			Error

		return err
	})

	if err != nil {
		slog.Error("[Changelog] Could not add reveiwer")
	}
}
