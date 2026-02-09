package biosynthesis

import (
	"errors"
	"github.com/adraismawur/mibig-submission/models"
	"gorm.io/gorm"
	"strconv"
)

type Biosynthesis struct {
	ID      uint64               `json:"-"`
	EntryID uint64               `json:"-"`
	Classes []BiosyntheticClass  `json:"classes" gorm:"foreignKey:BiosynthesisID"`
	Modules []BiosyntheticModule `json:"modules,omitempty" gorm:"foreignKey:BiosynthesisID"`
}

func init() {
	models.Models = append(models.Models, Biosynthesis{})
}

func GetEntryBiosynthesis(db *gorm.DB, accession string) (*Biosynthesis, error) {
	var biosynth Biosynthesis

	err := db.
		Table("entries").
		Where("accession = ?", accession).
		Preload("Classes").
		Preload("Modules.Carriers.Location").
		Preload("Modules.ModificationDomains.Location").
		Preload("Modules.ATDomain.Location").
		Preload("Modules.KSDomain.Location").
		First(&biosynth).
		Error

	if err != nil {
		return nil, err
	}

	return &biosynth, nil
}

func CreateEntryBiosynthesisModule(db *gorm.DB, accession string, module BiosyntheticModule) error {
	var biosynth *Biosynthesis

	err := db.
		Table("entries").
		Where("accession = ?", accession).
		Preload("Classes").
		Preload("Modules.Carriers.Location").
		Preload("Modules.ModificationDomains.Location").
		Preload("Modules.ADomain.Location").
		Preload("Modules.ATDomain.Location").
		Preload("Modules.KSDomain.Location").
		First(&biosynth).
		Error

	if err != nil {
		return err
	}

	// get new module number if it doesn't exist
	if module.Name == "" {
		module.Name = strconv.Itoa(len(biosynth.Modules) + 1)
	}

	err = db.Model(&biosynth).Association("Modules").Append(&module)

	if err != nil {
		return err
	}

	return nil
}

func UpdateEntryBiosynthesisModule(db *gorm.DB, accession string, newModule *BiosyntheticModule) error {
	var biosynth Biosynthesis

	err := db.
		Table("entries").
		Where("accession = ?", accession).
		Preload("Classes").
		Preload("Modules.Carriers.Location").
		Preload("Modules.ModificationDomains.Location").
		Preload("Modules.ADomain.Location").
		Preload("Modules.ATDomain.Location").
		Preload("Modules.KSDomain.Location").
		First(&biosynth).
		Error

	if err != nil {
		return err
	}

	// find correct newModule
	found := false
	for _, module := range biosynth.Modules {
		if module.Name == module.Name {
			newModule.ID = module.ID
			newModule.BiosynthesisID = module.BiosynthesisID
			found = true
			break
		}
	}

	if !found {
		return errors.New("could not find newModule")
	}

	err = db.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Model(&biosynth).
		Association("Modules").
		Replace(newModule)

	if err != nil {
		return err
	}

	return nil
}

func DeleteEntryBiosynthesisModule(db *gorm.DB, accession string, name string) error {
	var biosynth Biosynthesis

	err := db.
		Table("entries").
		Where("accession = ?", accession).
		Preload("Biosynthesis.Modules").
		First(&biosynth).
		Error

	if err != nil {
		return err
	}

	err = db.Where("name = ?", name).Delete(biosynth.Modules).Error

	if err != nil {
		return err
	}

	return nil
}

func GetEntryBiosynthesisModule(db *gorm.DB, accession string, name string) (*BiosyntheticModule, error) {
	var biosynth Biosynthesis

	err := db.
		Table("entries").
		Where("accession = ?", accession).
		Preload("Classes").
		Preload("Modules.Carriers.Location").
		Preload("Modules.ModificationDomains.Location").
		Preload("Modules.ADomain.Location").
		Preload("Modules.ATDomain.Location").
		Preload("Modules.KSDomain.Location").
		First(&biosynth).
		Error

	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	for _, module := range biosynth.Modules {
		if module.Name == name {
			return &module, nil
		}
	}

	// not found
	return nil, nil
}
