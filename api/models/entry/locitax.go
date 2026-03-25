package entry

import (
	"github.com/adraismawur/mibig-submission/models/entry/locus"
	"github.com/adraismawur/mibig-submission/models/entry/taxonomy"
	"gorm.io/gorm"
)

type LociTax struct {
	Accession string            `json:"accession" gorm:"primaryKey"`
	Loci      []locus.Locus     `json:"loci" gorm:"foreignKey:EntryAccession"`
	Taxonomy  taxonomy.Taxonomy `json:"taxonomy" gorm:"ForeignKey:EntryAccession"`
}

func GetLociTax(db *gorm.DB, accession string) (*LociTax, error) {
	var result LociTax

	err := db.
		Table("entries").
		Where("accession = $1", accession).
		Preload("Loci.Location").
		Preload("Loci.Evidence").
		Preload("Taxonomy").
		Find(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func UpdateLociTax(db *gorm.DB, accession string, oldLociTax LociTax, newLociTax LociTax) error {
	err := db.Session(&gorm.Session{FullSaveAssociations: true}).Transaction(func(tx *gorm.DB) error {
		err := tx.
			Model(&oldLociTax.Taxonomy).
			Where("entry_accession = $1", accession).
			Save(&newLociTax.Taxonomy).
			Error

		if err != nil {
			return err
		}

		err = tx.Model(&oldLociTax).
			Association("Loci").
			Replace(&newLociTax.Loci)

		if err != nil {
			return err
		}

		for _, loci := range oldLociTax.Loci {
			err = tx.
				Model(&loci).
				Association("Evidence").
				Replace(&loci.Evidence)

			if err != nil {
				return err
			}
		}

		err = tx.Model(&oldLociTax).
			Association("Taxonomy").
			Replace(&newLociTax.Taxonomy)

		if err != nil {
			return err
		}

		return nil
	})

	return err
}
