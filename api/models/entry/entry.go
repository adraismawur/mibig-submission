package entry

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

type Entry struct {
	ID           uint         `json:"-"`
	Accession    string       `json:"accession"`
	Version      int          `json:"version"`
	Quality      Quality      `json:"quality"`
	Status       Status       `json:"status"`
	Completeness Completeness `json:"completeness"`
}

func init() {
	models.Models = append(models.Models, &Entry{})
}

// ParseEntry attempts to parse an entry json given as a byte array into an entry struct
func ParseEntry(jsonString []byte) (*Entry, error) {
	var entry *Entry

	if err := json.Unmarshal(jsonString, &entry); err != nil {
		slog.Error("[sub] Failed to unmarshal annotation JSON", "error", err.Error())
		return nil, err
	}

	return entry, nil
}

// LoadEntry attempts to read a file at a given path and load it as an entry into the database
func LoadEntry(db *gorm.DB, path string) (*Entry, error) {
	jsonString := util.ReadFile(path)
	if jsonString == nil {
		slog.Error("[sub] Could not load entry file ", "path", path)
		return nil, errors.New("could not load entry file")
	}

	var err error
	var entry *Entry

	if entry, err = ParseEntry(jsonString); err != nil {
		slog.Error("[sub] Failed to parse entry file ", "path", path)
		return nil, err
	}

	if err = db.Create(entry).Error; err != nil {
		slog.Error("[db] Failed to create entry ", "error", err.Error())
		return nil, err
	}

	return entry, nil
}

func GetEntryExists(db *gorm.DB, accession string) (bool, error) {
	var count int64

	err := db.Table("entries").Where("accession = ?", accession).Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
