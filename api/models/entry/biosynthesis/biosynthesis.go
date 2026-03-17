package biosynthesis

import (
	"github.com/adraismawur/mibig-submission/models"
	"gorm.io/gorm"
)

type Biosynthesis struct {
	ID             uint64                `json:"db_id"`
	EntryAccession string                `json:"db_entry_accession"`
	Classes        []BiosyntheticClass   `json:"classes" gorm:"foreignKey:BiosynthesisID"`
	Modules        []BiosyntheticModule  `json:"modules,omitempty" gorm:"foreignKey:BiosynthesisID"`
	Operons        []BiosyntheticOperon  `json:"operons,omitempty" gorm:"foreignKey:BiosynthesisID"`
	Paths          []BiosyntheticPathway `json:"paths,omitempty" gorm:"foreignKey:BiosynthesisID"`
}

func init() {
	models.Models = append(models.Models, Biosynthesis{})
}

func GetEntryBiosynthesis(db *gorm.DB, accession string) (*Biosynthesis, error) {
	var biosynth Biosynthesis

	err := db.
		Table("biosyntheses").
		Where("entry_accession = ?", accession).
		Preload("Classes").
		Preload("Paths.Products").
		First(&biosynth).
		Error

	if err != nil {
		return nil, err
	}

	modules, err := GetEntryBiosynthesisModulesById(db, biosynth.ID)

	if err != nil {
		return nil, err
	}

	biosynth.Modules = *modules

	return &biosynth, nil
}

func GetBiosynthesisById(db *gorm.DB, id uint64) (*Biosynthesis, error) {
	var biosynth Biosynthesis

	err := db.
		Table("biosyntheses").
		Where("id = ?", id).
		Preload("Classes").
		First(&biosynth).
		Error

	if err != nil {
		return nil, err
	}

	return &biosynth, nil
}
