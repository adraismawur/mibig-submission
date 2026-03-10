package entry_utils

import (
	"fmt"
	"github.com/adraismawur/mibig-submission/models/entry"
	"gorm.io/gorm"
	"log/slog"
	"strconv"
)

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

	err = db.Model(&entry.Entry{}).
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

// GenerateNewAccession generates a sequentially increased new accession number from the previous highest
func GenerateNewAccession(db *gorm.DB) (string, error) {

	type EntryAccession struct {
		Accession string
	}
	var lastAccession EntryAccession

	err := db.Table("entries").Select("accession").Last(&lastAccession).Error

	if err != nil {
		slog.Error("[util] [entry] Could not get last accession", "err", err)
		return "", err
	}

	numericPart, err := strconv.Atoi(lastAccession.Accession[3:])

	if err != nil {
		slog.Error("[util] [entry] Could not convert last accession numeric part to int", "accession", lastAccession)
		return "", err
	}

	numericPart = numericPart + 1

	newAccession := fmt.Sprintf("BGC%07d", numericPart)

	return newAccession, nil
}
