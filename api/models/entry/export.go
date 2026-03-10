package entry

import (
	"github.com/adraismawur/mibig-submission/models/entry/export"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetEntryExportFromAccession(db *gorm.DB, accession string) (*export.Entry, error) {
	var entry export.Entry

	// TODO: there must be a better way to do this. no amount of googling on my part gets me anywhere though
	// ideally doing this amount of preloading is rare. This is done here on getting the entire entry
	err := db.
		Table("entries").
		Where("accession = ?", accession).
		Preload("Changelog.Releases.Entries").
		Preload("Loci.Location").
		Preload("Loci.Evidence").
		Preload("Biosynthesis.Classes").
		Preload("Biosynthesis.Modules.Carriers.Location").
		Preload("Biosynthesis.Modules.ModificationDomains.Location").
		Preload("Biosynthesis.Modules.ATDomain.Location").
		Preload("Biosynthesis.Modules.KSDomain.Location").
		Preload("GeneInformation.Additions.Location.Exons").
		Preload("GeneInformation.Annotations").
		Preload("Compounds.Evidence").
		Preload("Compounds.BioActivities").
		Preload(clause.Associations).
		First(&entry).
		Error

	if err != nil {
		return nil, err
	}

	return &entry, nil
}
