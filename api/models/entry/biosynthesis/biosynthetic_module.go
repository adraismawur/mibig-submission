package biosynthesis

import (
	"github.com/adraismawur/mibig-submission/config"
	"github.com/adraismawur/mibig-submission/models"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"log/slog"
)

type IntegratedMonomer struct {
	ID                   uint64                    `json:"db_id"`
	BiosyntheticModuleID uint64                    `json:"db_biosynth_module_id"`
	Name                 string                    `json:"name"`
	Structure            string                    `json:"structure"`
	Evidence             []DomainSubstrateEvidence `json:"evidence" gorm:"many2many:integrated_monomer_evidences"`
}

type BiosyntheticModule struct {
	ID                  uint64                   `json:"db_id"`
	Index               uint64                   `json:"db_index"`
	BiosynthesisID      uint64                   `json:"db_biosynth_id"`
	Type                string                   `json:"type"`
	Name                string                   `json:"name"`
	Genes               pq.StringArray           `json:"genes" gorm:"type:text[]"`
	Iterations          *uint64                  `json:"iterations,omitempty"`
	Active              bool                     `json:"active"`
	IntegratedMonomers  []IntegratedMonomer      `json:"integrated_monomers" gorm:"foreignKey:BiosyntheticModuleID"`
	Carriers            []CarrierDomain          `json:"carriers" gorm:"many2many:biosynth_carrier_domains"`
	ModificationDomains []ModificationDomain     `json:"modification_domains,omitempty" gorm:"many2many:biosynth_modification_domains"`
	CDomainID           uint64                   `json:"db_c_domain_id"`
	CDomain             *CondensationDomain      `json:"c_domain,omitempty"`
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
	bioSynth, err := GetEntryBiosynthesis(db, entryAccession)

	if err != nil {
		return err
	}

	err = db.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Model(&bioSynth).
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

func UpdateEntryBiosynthesisModule(db *gorm.DB, newModule BiosyntheticModule) error {
	err := db.Session(&gorm.Session{FullSaveAssociations: true}).Transaction(func(tx *gorm.DB) error {
		var err error

		oldModule, err := GetEntryBiosynthesisModule(tx, int(newModule.ID))

		err = tx.
			Model(&oldModule).
			Where("biosynthetic_modules.id = $1", newModule.ID).
			Omit("index").
			Save(&newModule).
			Error

		if err != nil {
			return err
		}

		err = tx.
			Model(&oldModule).
			Association("IntegratedMonomers").
			Replace(&newModule.IntegratedMonomers)

		if err != nil {
			return err
		}

		err = tx.
			Model(&oldModule).
			Association("Carriers").
			Replace(&newModule.Carriers)

		if err != nil {
			return err
		}

		for _, carrier := range newModule.Carriers {

			err = tx.
				Model(&DomainLocation{}).
				Where("id = $1", carrier.Location.ID).
				Save(&carrier.Location).
				Error

			if err != nil {
				return err
			}
		}

		for _, modificationDomain := range newModule.ModificationDomains {

			err = tx.
				Model(&DomainLocation{}).
				Where("id = $1", modificationDomain.Location.ID).
				Save(&modificationDomain.Location).
				Error

			if err != nil {
				return err
			}
		}

		if newModule.CDomain != nil {
			err = tx.
				Model(&CondensationDomain{}).
				Where("id = $1", newModule.CDomain.ID).
				Save(&newModule.CDomain).
				Error

			if err != nil {
				return err
			}

			err = tx.
				Model(&DomainLocation{}).
				Where("id = $1", newModule.CDomain.Location.ID).
				Save(&newModule.CDomain.Location).
				Error

			if err != nil {
				return err
			}
		}

		if newModule.ADomain != nil {
			err = tx.
				Model(&AdenylationDomain{}).
				Where("id = $1", newModule.ADomain.ID).
				Save(&newModule.ADomain).
				Error

			if err != nil {
				return err
			}

			err = tx.
				Model(&DomainLocation{}).
				Where("id = $1", newModule.ADomain.Location.ID).
				Save(&newModule.ADomain.Location).
				Error

			if err != nil {
				return err
			}

			err = tx.
				Model(&newModule.ADomain).
				Association("Evidence").
				Replace(&newModule.ADomain.Evidence)

			if err != nil {
				return err
			}

			err = tx.
				Model(&newModule.ADomain).
				Association("Substrates").
				Replace(&newModule.ADomain.Substrates)

			if err != nil {
				return err
			}
		}

		if newModule.ATDomain != nil {
			err = tx.
				Model(&AcetyltransferaseDomain{}).
				Where("id = $1", newModule.ATDomain.ID).
				Save(&newModule.ATDomain).
				Error

			if err != nil {
				return err
			}

			err = tx.
				Model(&DomainLocation{}).
				Where("id = $1", newModule.ATDomain.Location.ID).
				Save(&newModule.ATDomain.Location).
				Error

			if err != nil {
				return err
			}

			err = tx.
				Model(&newModule.ATDomain).
				Association("Substrates").
				Replace(&newModule.ATDomain.Substrates)

			if err != nil {
				return err
			}

			err = tx.
				Model(&newModule.ATDomain).
				Association("Evidence").
				Replace(&newModule.ATDomain.Evidence)

			if err != nil {
				return err
			}
		}

		if newModule.KSDomain != nil {
			err = tx.
				Model(&KetoSynthaseDomain{}).
				Where("id = $1", newModule.KSDomain.ID).
				Save(&newModule.KSDomain).
				Error

			if err != nil {
				return err
			}

			err = tx.
				Model(&DomainLocation{}).
				Where("id = $1", newModule.KSDomain.Location.ID).
				Save(&newModule.KSDomain.Location).
				Error

			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
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
		Preload("IntegratedMonomers").
		Preload("IntegratedMonomers.Evidence").
		Preload("Carriers.Location").
		Preload("Carriers.Evidence").
		Preload("ModificationDomains.Location").
		Preload("ModificationDomains.Substrates").
		Preload("ModificationDomains.Evidence").
		Preload("CDomain.Location").
		Preload("ADomain.Location").
		Preload("ADomain.Evidence").
		Preload("ADomain.Substrates").
		Preload("ATDomain.Location").
		Preload("ATDomain.Substrates").
		Preload("ATDomain.Evidence").
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

	q := db.
		Table("biosynthetic_modules").
		Where("biosynthesis_id = $1", biosynthId)

	dbDialect, err := config.GetConfig(config.EnvDbDialect)

	if err != nil {
		return nil, err
	}

	switch config.EnvValue(dbDialect) {
	case config.DbDialectSqlite:
		q = q.Order("`index` ASC")
		break
	default:
		q = q.Order("index")
		break
	}

	if dbDialect == string(config.DbDialectPostgres) {
	}

	err = q.Find(&modules).Error

	if err != nil {
		return nil, err
	}

	// not found
	return &modules, nil
}
