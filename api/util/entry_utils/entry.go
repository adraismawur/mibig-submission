package entry_utils

import (
	"fmt"
	"github.com/adraismawur/mibig-submission/models"
	"gorm.io/gorm"
	"log/slog"
	"strconv"
	"time"
)

func GeneratePlaceholderAccession(user models.User) string {
	newPart := "new"
	datePart := time.Now().Format("2006-01-02-15-04-05")

	return newPart + "-" + datePart
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
