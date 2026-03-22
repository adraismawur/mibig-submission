package export

import (
	"github.com/lib/pq"
)

type AcetyltransferaseDomain struct {
	Type       string                    `json:"type"`
	Subtype    string                    `json:"subtype"`
	Gene       string                    `json:"gene"`
	Location   DomainLocation            `json:"location"`
	Inactive   bool                      `json:"inactive,omitempty"`
	Substrates []DomainSubstrate         `json:"substrates,omitempty" gorm:"many2many:acetyltransferase_substrates;"`
	Evidence   []DomainSubstrateEvidence `json:"evidence,omitempty" gorm:"many2many:acetyltransferase_evidences;"`
}

type AdenylationDomain struct {
	Type                  string                    `json:"type"`
	Gene                  string                    `json:"gene"`
	Location              DomainLocation            `json:"location"`
	Inactive              bool                      `json:"inactive"`
	Evidence              []DomainSubstrateEvidence `json:"evidence" gorm:"many2many:adenylation_evidences;"`
	PrecursorBiosynthesis pq.StringArray            `json:"precursor_biosynthesis" gorm:"type:text[]"`
	Substrates            []DomainSubstrate         `json:"substrates,omitempty" gorm:"many2many:adenylation_substrates;"`
}

type CarrierDomain struct {
	Type          string                    `json:"type"`
	Subtype       string                    `json:"subtype"`
	Gene          string                    `json:"gene"`
	Location      DomainLocation            `json:"location"`
	Inactive      bool                      `json:"inactive"`
	BetaBranching bool                      `json:"beta_branching"`
	Evidence      []DomainSubstrateEvidence `json:"evidence" gorm:"many2many:carrier_evidences;"`
}

type CondensationDomain struct {
	Type       string         `json:"type"`
	Gene       string         `json:"gene"`
	Location   DomainLocation `json:"location"`
	Subtype    string         `json:"subtype"`
	References pq.StringArray `json:"references" gorm:"type:text[]"`
}

type KetoSynthaseDomain struct {
	Type     string         `json:"type"`
	Gene     string         `json:"gene"`
	Location DomainLocation `json:"location"`
}

type DomainLocation struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type ModificationDomain struct {
	Type            string                    `json:"type"`
	Subtype         string                    `json:"subtype"`
	Gene            string                    `json:"gene"`
	Location        DomainLocation            `json:"location"`
	Inactive        *bool                     `json:"inactive,omitempty"`
	Substrates      []DomainSubstrate         `json:"substrates,omitempty" gorm:"many2many:modification_domain_substrates;"`
	Evidence        []DomainSubstrateEvidence `json:"evidence,omitempty" gorm:"many2many:modification_domain_evidences;"`
	References      pq.StringArray            `json:"references" gorm:"type:text[]"`
	Stereochemistry string                    `json:"stereochemistry,omitempty"`
	Details         string                    `json:"details,omitempty"`
}

type DomainSubstrate struct {
	Name      string `json:"name"`
	Details   string `json:"details,omitempty"`
	Structure string `json:"structure,omitempty"`
}

type DomainSubstrateEvidence struct {
	Method     string         `json:"method"`
	References pq.StringArray `json:"references" gorm:"type:text[]"`
}
