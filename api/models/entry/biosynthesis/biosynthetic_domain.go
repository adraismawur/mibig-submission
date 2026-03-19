package biosynthesis

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/lib/pq"
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
	References      pq.StringArray            `json:"references" gorm:"type:text[]"`
	Stereochemistry pq.StringArray            `json:"stereochemistry" gorm:"type:text[]"`
	Details         string                    `json:"details,omitempty"`
}

type DomainSubstrate struct {
	ID        uint64 `json:"db_id"`
	Name      string `json:"name"`
	Details   string `json:"details,omitempty"`
	Structure string `json:"structure,omitempty"`
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
