package entry

import (
	"github.com/adraismawur/mibig-submission/models/entry/export"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

func GetEntryExportFromAccession(db *gorm.DB, accession string) (*export.Entry, error) {
	entry, err := GetEntryFromAccession(db, accession)

	if err != nil {
		return nil, err
	}

	var entryExport export.Entry

	mapstructure.Decode(entry, &entryExport)

	return &entryExport, nil
}
