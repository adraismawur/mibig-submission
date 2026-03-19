package biosynthesis

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"log/slog"
	"strconv"
)

type IntegratedMonomer struct {
	ID                   uint64         `json:"db_id"`
	BiosyntheticModuleID uint64         `json:"db_biosynth_module_id"`
	Name                 string         `json:"name"`
	Structure            string         `json:"structure"`
	References           pq.StringArray `json:"references" gorm:"type:text[]"`
}

type BiosyntheticModule struct {
	ID                  uint64                   `json:"db_id"`
	Index               uint64                   `json:"db_index"`
	BiosynthesisID      uint64                   `json:"db_biosynth_id"`
	Type                string                   `json:"type"`
	Name                string                   `json:"name"`
	Genes               pq.StringArray           `json:"genes" gorm:"type:text[]"`
	Active              bool                     `json:"active"`
	IntegratedMonomers  []IntegratedMonomer      `json:"integrated_monomers" gorm:"foreignKey:BiosyntheticModuleID"`
	Carriers            []CarrierDomain          `json:"carriers" gorm:"many2many:biosynth_carrier_domains"`
	ModificationDomains []ModificationDomain     `json:"modification_domains,omitempty" gorm:"many2many:biosynth_modification_domains"`
	CDomainID           uint64                   `json:"db_c_domain_id"`
	CDomain             *CondensationDomain      `json:"c_domain"`
	ADomainID           uint64                   `json:"db_a_domain_id"`
	ADomain             *AdenylationDomain       `json:"a_domain,omitempty"`
	ATDomainID          uint64                   `json:"db_at_domain_id"`
	ATDomain            *AcetyltransferaseDomain `json:"at_domain,omitempty"`
	KSDomainID          uint64                   `json:"db_ks_domain_id"`
	KSDomain            *KetoSynthaseDomain      `json:"ks_domain,omitempty"`
}

func init() {
	models.Models = append(models.Models, BiosyntheticModule{})
	models.Models = append(models.Models, IntegratedMonomer{})
}

func CreateEntryBiosynthesisModule(db *gorm.DB, entryAccession string, module BiosyntheticModule) error {
	var biosynth *Biosynthesis

	err := db.
		Table("biosyntheses").
		Where("entry_accession = $1", entryAccession).
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

	module.Index = uint64(len(biosynth.Modules)) + 0x1

	err = db.
		Model(&biosynth).
		Association("Modules").
		Append(&module)

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
		Where("id = $1", idFrom).
		First(&moduleFrom).
		Error

	if err != nil {
		slog.Error("Could not get first module in module reorder operation")
		tx.Rollback()
		return err
	}

	err = tx.
		Where("id = $1", idTo).
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

	err = tx.Select("index").Save(&moduleFrom).Error

	if err != nil {
		slog.Error("Could not save first module in module reorder operation")
		tx.Rollback()
		return err
	}

	err = tx.Select("index").Save(&moduleTo).Error

	if err != nil {
		slog.Error("Could not save second module in module reorder operation")
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func UpdateEntryBiosynthesisModule(db *gorm.DB, newModule *BiosyntheticModule) error {
	err := db.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Model(&BiosyntheticModule{}).
		Preload("Carriers.Location").
		Preload("ModificationDomains.Location").
		Preload("ADomain.Location").
		Preload("ATDomain.Location").
		Preload("KSDomain.Location").
		Where("id = $1", newModule.ID).
		Save(&newModule).
		Error

	// TODO: replace associations. guh

	if err != nil {
		return err
	}

	return nil
}

func DeleteEntryBiosynthesisModule(db *gorm.DB, id int) error {
	err := db.
		Model(&BiosyntheticModule{}).
		Delete("id = $1", id).
		Error

	if err != nil {
		return err
	}

	return nil
}

func GetEntryBiosynthesisModule(db *gorm.DB, id int) (*BiosyntheticModule, error) {
	var module BiosyntheticModule

	err := db.
		Table("biosynthetic_modules").
		Where("id = $1", id).
		Preload("Carriers.Location").
		Preload("ModificationDomains.Location").
		Preload("ADomain.Location").
		Preload("ATDomain.Location").
		Preload("KSDomain.Location").
		First(&module).
		Error

	if err != nil {
		return nil, err
	}

	return &module, nil
}

func GetEntryBiosynthesisModulesById(db *gorm.DB, biosynthId uint64) (*[]BiosyntheticModule, error) {
	var modules []BiosyntheticModule

	err := db.
		Table("biosynthetic_modules").
		Where("biosynthesis_id = $1", biosynthId).
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
