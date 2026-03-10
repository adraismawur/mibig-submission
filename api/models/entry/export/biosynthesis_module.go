package export

import "github.com/lib/pq"

type BiosyntheticModuleDomainLocation struct {
	ID                   uint64 `json:"-"`
	BiosyntheticModuleID uint64 `json:"-"`
	From                 int    `json:"from"`
	To                   int    `json:"to"`
}

type BiosyntheticModuleDomain struct {
	ID                   uint64                           `json:"-"`
	BiosyntheticModuleID uint64                           `json:"-"`
	DomainType           string                           `json:"type"`
	Gene                 string                           `json:"gene"`
	Location             BiosyntheticModuleDomainLocation `gorm:"foreignKey:BiosyntheticModuleID"`
}

type CarrierModuleDomain struct {
	BiosyntheticModuleDomain
	Subtype       string `json:"subtype,omitempty"`
	BetaBranching bool   `json:"beta_branching"`
}

type AModuleDomainSubstrate struct {
	ID              uint64 `json:"-"`
	AModuleDomainID uint64 `json:"-"`
	Name            string `json:"name"`
	Structure       string `json:"structure"`
	Proteinogenic   bool   `json:"proteinogenic"`
}

type AModuleDomain struct {
	BiosyntheticModuleDomain
	References pq.StringArray           `json:"references" gorm:"type:text[]"`
	Substrates []AModuleDomainSubstrate `json:"substrates" gorm:"foreignKey:AModuleDomainID"`
}

type ATModuleDomain struct {
	BiosyntheticModuleDomain
	Substrates pq.StringArray `json:"substrates" gorm:"type:text[]"`
	Evidence   pq.StringArray `json:"evidence" gorm:"type:text[]"`
}

type KSModuleDomain struct {
	BiosyntheticModuleDomain
}

type ModificationModuleDomain struct {
	BiosyntheticModuleDomain
}

type BiosyntheticModule struct {
	ID                  uint64                     `json:"-"`
	Index               uint64                     `json:"-"`
	BiosynthesisID      uint64                     `json:"-"`
	Type                string                     `json:"type"`
	Name                string                     `json:"name"`
	Genes               pq.StringArray             `json:"genes" gorm:"type:text[]"`
	Active              bool                       `json:"active"`
	Carriers            []CarrierModuleDomain      `json:"carriers" gorm:"foreignKey:BiosyntheticModuleID"`
	ModificationDomains []ModificationModuleDomain `json:"modification_domains,omitempty" gorm:"foreignKey:BiosyntheticModuleID"`
	ADomain             *AModuleDomain             `json:"a_domain,omitempty" gorm:"foreignKey:BiosyntheticModuleID"`
	ATDomain            *ATModuleDomain            `json:"at_domain,omitempty" gorm:"foreignKey:BiosyntheticModuleID"`
	KSDomain            *KSModuleDomain            `json:"ks_domain,omitempty" gorm:"foreignKey:BiosyntheticModuleID"`
}
