package entry_utils

import (
	"fmt"
	"gorm.io/gorm"
	"log/slog"
	"strconv"
)

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
