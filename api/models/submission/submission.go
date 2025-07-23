package submission

import (
	"errors"
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/util"
	"github.com/goccy/go-json"
	"gorm.io/gorm"
	"log/slog"
)

type Quality string

const (
	Questionable Quality = "questionable"
	Medium       Quality = "medium"
	High         Quality = "high"
)

type Status string

const (
	Pending Status = "pending"
	Active  Status = "active"
	Retired Status = "retired"
)

type Completeness string

const (
	Unknown  Completeness = "unknown"
	Complete Completeness = "complete"
)

type Submission struct {
	ID           uint         `json:"-"`
	Accession    string       `json:"accession"`
	Version      int          `json:"version"`
	Quality      Quality      `json:"quality"`
	Status       Status       `json:"status"`
	Completeness Completeness `json:"completeness"`
}

func init() {
	models.Models = append(models.Models, &Submission{})
}

// ParseSubmission attempts to parse a submission json given as a byte array into a submission struct
func ParseSubmission(jsonString []byte) (*Submission, error) {
	var submission *Submission

	if err := json.Unmarshal(jsonString, &submission); err != nil {
		slog.Error("[sub] Failed to unmarshal annotation JSON", "error", err.Error())
		return nil, err
	}

	return submission, nil
}

// LoadSubmission attempts to read a file at a given path and load it as a submission into the database
func LoadSubmission(db *gorm.DB, path string) (*Submission, error) {
	jsonString := util.ReadFile(path)
	if jsonString == nil {
		slog.Error("[sub] Could not load submission file ", "path", path)
		return nil, errors.New("could not load submission file")
	}

	var err error
	var submission *Submission

	if submission, err = ParseSubmission(jsonString); err != nil {
		slog.Error("[sub] Failed to parse submission file ", "path", path)
		return nil, err
	}

	if err = db.Create(submission).Error; err != nil {
		slog.Error("[db] Failed to create submission ", "error", err.Error())
		return nil, err
	}

	return submission, nil
}

func GetSubmissionExists(db *gorm.DB, accession string) (bool, error) {
	var count int64

	err := db.Table("submissions").Where("accession = ?", accession).Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
