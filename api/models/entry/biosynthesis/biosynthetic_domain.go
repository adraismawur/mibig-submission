package biosynthesis

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type AcetyltransferaseDomain struct {
	ID         uint64                    `json:"db_id"`
	Type       string                    `json:"type"`
	Subtype    string                    `json:"subtype"`
	Gene       string                    `json:"gene"`
	LocationID uint64                    `json:"db_location_id"`
	Location   DomainLocation            `json:"location"`
	Inactive   bool                      `json:"inactive,omitempty"`
	Substrates []DomainSubstrate         `json:"substrates,omitempty" gorm:"many2many:acetyltransferase_substrates;"`
	Evidence   []DomainSubstrateEvidence `json:"evidence,omitempty" gorm:"many2many:acetyltransferase_evidences;"`
}

type AdenylationDomain struct {
	ID                    uint64                    `json:"db_id"`
	Type                  string                    `json:"type"`
	Gene                  string                    `json:"gene"`
	LocationID            uint64                    `json:"db_location_id"`
	Location              DomainLocation            `json:"location"`
	Inactive              bool                      `json:"inactive"`
	Evidence              []DomainSubstrateEvidence `json:"evidence" gorm:"many2many:adenylation_evidences;"`
	PrecursorBiosynthesis pq.StringArray            `json:"precursor_biosynthesis" gorm:"type:text[]"`
	Substrates            []DomainSubstrate         `json:"substrates,omitempty" gorm:"many2many:adenylation_substrates;"`
}

type CarrierDomain struct {
	ID            uint64                    `json:"db_id"`
	Type          string                    `json:"type"`
	Subtype       string                    `json:"subtype"`
	Gene          string                    `json:"gene"`
	LocationID    uint64                    `json:"db_location_id"`
	Location      DomainLocation            `json:"location"`
	Inactive      bool                      `json:"inactive"`
	BetaBranching bool                      `json:"beta_branching"`
	Evidence      []DomainSubstrateEvidence `json:"evidence" gorm:"many2many:carrier_evidences;"`
}

type CondensationDomain struct {
	ID         uint64         `json:"db_id"`
	Type       string         `json:"type"`
	Gene       string         `json:"gene"`
	LocationID uint64         `json:"db_location_id"`
	Location   DomainLocation `json:"location"`
	Subtype    string         `json:"subtype"`
	References pq.StringArray `json:"references" gorm:"type:text[]"`
}

type KetoSynthaseDomain struct {
	ID         uint64         `json:"db_id"`
	Type       string         `json:"type"`
	Gene       string         `json:"gene"`
	LocationID uint64         `json:"db_location_id"`
	Location   DomainLocation `json:"location"`
}

type DomainLocation struct {
	ID   uint64 `json:"db_id"`
	From int    `json:"from"`
	To   int    `json:"to"`
}

type ModificationDomain struct {
	ID              uint64                    `json:"db_id"`
	Type            string                    `json:"type"`
	Subtype         string                    `json:"subtype"`
	Gene            string                    `json:"gene"`
	LocationID      uint64                    `json:"db_location_id"`
	Location        DomainLocation            `json:"location"`
	Inactive        *bool                     `json:"inactive,omitempty"`
	Substrates      []DomainSubstrate         `json:"substrates,omitempty" gorm:"many2many:modification_domain_substrates;"`
	Evidence        []DomainSubstrateEvidence `json:"evidence,omitempty" gorm:"many2many:modification_domain_evidences;"`
	References      pq.StringArray            `json:"references,omitempty" gorm:"type:text[]"`
	Stereochemistry string                    `json:"stereochemistry,omitempty"`
	Details         string                    `json:"details,omitempty"`
}

type DomainSubstrate struct {
	ID            uint64 `json:"db_id"`
	Name          string `json:"name"`
	Details       string `json:"details,omitempty"`
	Proteinogenic *bool  `json:"proteinogenic,omitempty"`
	Structure     string `json:"structure,omitempty"`
}

type DomainSubstrateEvidence struct {
	ID         uint64         `json:"db_id"`
	Method     string         `json:"method"`
	References pq.StringArray `json:"references" gorm:"type:text[]"`
}

func init() {
	// we only need to load in the models that exist independently in the database
	// a lot of the above are modification domains and do not need their own table
	models.Models = append(models.Models, AdenylationDomain{})       // A domain exists on NRPS modules
	models.Models = append(models.Models, CondensationDomain{})      // C domain exists on NRPS modules
	models.Models = append(models.Models, CarrierDomain{})           // Carrier domains exist on NRPS modules
	models.Models = append(models.Models, KetoSynthaseDomain{})      // KS domain exists on PKS modules
	models.Models = append(models.Models, AcetyltransferaseDomain{}) // AT domain exists on PKS modules
	models.Models = append(models.Models, ModificationDomain{})
	models.Models = append(models.Models, DomainSubstrate{})
	models.Models = append(models.Models, DomainSubstrateEvidence{})
}

func CreateModificationDomain(db *gorm.DB, biosyntheticModuleID int, newModificationDomain ModificationDomain) error {
	bioSynth, err := GetEntryBiosynthesisModule(db, biosyntheticModuleID)

	if err != nil {
		return err
	}

	err = db.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Model(&bioSynth).
		Association("ModificationDomains").
		Append(&newModificationDomain)

	if err != nil {
		return err
	}

	return nil
}

func GetModificationDomains(db *gorm.DB, moduleID int) ([]ModificationDomain, error) {
	var module BiosyntheticModule

	err := db.
		Model(&BiosyntheticModule{}).
		Where("id = $1", moduleID).
		Preload("ModificationDomains.Location").
		Preload("ModificationDomains.Substrates").
		Preload("ModificationDomains.Evidence").
		First(&module).
		Error

	if err != nil {
		return nil, err
	}

	return module.ModificationDomains, nil
}

func GetModificationDomain(db *gorm.DB, modificationDomainID int) (*ModificationDomain, error) {
	var modificationDomain ModificationDomain

	err := db.
		Table("modification_domains").
		Where("id = $1", modificationDomainID).
		Preload("Location").
		Preload("Substrates").
		Preload("Evidence").
		First(&modificationDomain).
		Error

	if err != nil {
		return nil, err
	}

	return &modificationDomain, nil
}

func UpdateModificationDomain(db *gorm.DB, newModificationDomain ModificationDomain) error {
	err := db.Session(&gorm.Session{FullSaveAssociations: true}).Transaction(func(tx *gorm.DB) error {
		var transactionErr error

		oldModificationDomain, transactionErr := GetModificationDomain(tx, int(newModificationDomain.ID))

		if transactionErr != nil {
			return transactionErr
		}

		transactionErr = tx.
			Model(&oldModificationDomain).
			Save(newModificationDomain).
			Error

		if transactionErr != nil {
			return transactionErr
		}

		transactionErr = tx.
			Model(&oldModificationDomain.Location).
			Save(&newModificationDomain.Location).
			Error

		if transactionErr != nil {
			return transactionErr
		}

		transactionErr = tx.
			Model(&oldModificationDomain).
			Association("Substrates").
			Replace(newModificationDomain.Substrates)

		if transactionErr != nil {
			return transactionErr
		}

		transactionErr = tx.
			Model(&oldModificationDomain).
			Association("Evidence").
			Replace(newModificationDomain.Evidence)

		if transactionErr != nil {
			return transactionErr
		}

		return nil
	})

	return err
}

func DeleteModificationDomain(db *gorm.DB, modificationDomainID int) error {
	err := db.
		Model(&ModificationDomain{}).
		Delete("id = $1", modificationDomainID).
		Error

	if err != nil {
		return err
	}

	return nil
}
