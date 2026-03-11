package biosynthesis

import (
	"errors"
	"github.com/adraismawur/mibig-submission/models"
	"gorm.io/gorm"
	"log/slog"
	"strconv"
)

type Biosynthesis struct {
	ID      uint64               `json:"db_id"`
	EntryID uint64               `json:"entry_id"`
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

func CreateBiosynthesisClass(db *gorm.DB, biosynthId uint64, class BiosyntheticClass) error {
	bioSynth, err := GetBiosynthesisById(db, biosynthId)

	if err != nil {
		return err
	}

	err = db.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Model(&bioSynth).
		Association("Classes").
		Append(&class)

	if err != nil {
		return err
	}

	return nil
}

func UpdateEntryBiosynthesisClass(db *gorm.DB, classId int, newClass BiosyntheticClass) error {
	var biosynth Biosynthesis

	err := db.
		Model(&biosynth).
		Where("id = ?", newClass.BiosynthesisID).
		Preload("Classes").
		First(&biosynth).
		Error

	if err != nil {
		return err
	}

	err = db.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Model(&biosynth).
		Association("Classes").
		Replace(&newClass)

	if err != nil {
		return err
	}

	return nil
}

func DeleteEntryBiosynthesisClass(db *gorm.DB, id int) error {
	err := db.
		Model(&BiosyntheticClass{}).
		Delete("id = ?", id).
		Error

	return err
}

func GetEntryBiosynthesisClass(db *gorm.DB, id int) (*BiosyntheticClass, error) {
	var class BiosyntheticClass

	err := db.
		Table("biosynthetic_classes").
		Where("id = ?", id).
		First(&class).
		Error

	if err != nil {
		return nil, err
	}

	// oh god
	switch class.Class {
	case "NRPS":
		db.Model(&class).
			Preload("ReleaseTypes").
			Preload("Thioesterases.Location").
			Find(&class)
	case "PKS":
		break
	case "ribosomal":
		db.Model(&class).
			Preload("Precursors.LeaderCleavageLocation").
			Preload("Precursors.FollowerCleavageLocation").
			Preload("Precursors.Crosslinks").
			Find(&class)
	case "saccharide":
		db.Model(&class).
			Preload("GlycosylTransferases.Evidence").
			Preload("Subclusters").
			Find(&class)
	case "terpene":
		break
	case "other":
		break
	}

	return &class, nil
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

// ReorderEntryBiosynthesisModules swaps the indexes of biosynthesis modules. This uses database IDs as input, not the
// indexes
func ReorderEntryBiosynthesisModules(db *gorm.DB, idFrom uint64, idTo uint64) error {

	tx := db.Table("biosynthetic_modules").Begin()

	var moduleFrom BiosyntheticModule
	var moduleTo BiosyntheticModule

	err := tx.
		Where("id = ?", idFrom).
		First(&moduleFrom).
		Error

	if err != nil {
		slog.Error("Could not get first module in module reorder operation")
		tx.Rollback()
		return err
	}

	err = tx.
		Where("id = ?", idTo).
		First(&moduleTo).
		Error

	if err != nil {
		slog.Error("Could not get second module in module reorder operation")
		tx.Rollback()
		return err
	}

	swap := moduleTo.Index

	moduleTo.Index = moduleFrom.Index
	moduleFrom.Index = swap

	err = tx.Save(&moduleFrom).Error

	if err != nil {
		slog.Error("Could not save first module in module reorder operation")
		tx.Rollback()
		return err
	}

	err = tx.Save(&moduleTo).Error

	if err != nil {
		slog.Error("Could not save second module in module reorder operation")
		tx.Rollback()
		return err
	}

	tx.Commit()

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
		Preload("Modules").
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

	for _, module := range biosynth.Modules {
		if module.Name == name {
			return &module, nil
		}
	}

	// not found
	return nil, nil
}

func GetEntryBiosynthesisModulesById(db *gorm.DB, biosynthId uint64) (*[]BiosyntheticModule, error) {
	var modules []BiosyntheticModule

	err := db.
		Table("biosynthetic_modules").
		Where("biosynthesis_id = ?", biosynthId).
		Preload("Carriers.Location").
		Preload("ModificationDomains.Location").
		Preload("ADomain.Location").
		Preload("ATDomain.Location").
		Preload("KSDomain.Location").
		Order("`index` ASC").
		Find(&modules).
		Error

	if err != nil {
		return nil, err
	}

	// not found
	return &modules, nil
}
