package entry

import (
	"github.com/adraismawur/mibig-submission/models/entry/locus"
	"github.com/adraismawur/mibig-submission/models/entry/taxonomy"
	"gorm.io/gorm"
)

type LociTax struct {
	ID       uint              `json:"-"`
	Loci     []locus.Locus     `json:"loci" gorm:"foreignKey:EntryID"`
	Taxonomy taxonomy.Taxonomy `json:"taxonomy" gorm:"ForeignKey:EntryID"`
}

func GetLociTax(db *gorm.DB, accession string) (*LociTax, error) {
	var result LociTax

	err := db.
		Table("entries").
		Where("accession = ?", accession).
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
	tx := db.Begin()

	err := tx.
		Table("entries").
		Model(&oldLociTax).
		Where("accession = ?", accession).
		Save(&newLociTax).
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

	err = tx.Model(&oldLociTax).
		Association("Taxonomy").
		Replace(&newLociTax.Taxonomy)

	if err != nil {
		return err
	}

	err = tx.Commit().Error

	if err != nil {
		return err
	}

	return nil
}
