package biosynthesis

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/lib/pq"
)

// TODO: maybe unify with loci.location and gene location
type BiosyntheticModuleDomainLocation struct {
	BiosyntheticModuleID uint64 `json:"-"`
	From                 int    `json:"from"`
	To                   int    `json:"to"`
}

type BiosyntheticModuleDomain struct {
	ID                   uint64                           `json:"-"`
	BiosyntheticModuleID uint64                           `json:"-"`
	DomainType           string                           `json:"type"`
	Gene                 string                           `json:"gene"`
	Location             BiosyntheticModuleDomainLocation `json:"location" gorm:"foreignKey:BiosyntheticModuleID"`
}

type CarrierModuleDomain struct {
	BiosyntheticModuleDomain
	Subtype       string `json:"subtype,omitempty"`
	BetaBranching bool   `json:"beta_branching"`
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
	BiosynthesisID      uint64                     `json:"-"`
	Type                string                     `json:"type"`
	Name                string                     `json:"name"`
	Genes               pq.StringArray             `json:"genes" gorm:"type:text[]"`
	Active              bool                       `json:"active"`
	Carriers            []CarrierModuleDomain      `json:"carriers" gorm:"foreignKey:BiosyntheticModuleID"`
	ModificationDomains []ModificationModuleDomain `json:"modification_domains,omitempty" gorm:"foreignKey:BiosyntheticModuleID"`
	ATDomain            ATModuleDomain             `json:"at_domain,omitempty" gorm:"foreignKey:BiosyntheticModuleID"`
	KSDomain            KSModuleDomain             `json:"ks_domain,omitempty" gorm:"foreignKey:BiosyntheticModuleID"`
}

func init() {
	models.Models = append(models.Models, BiosyntheticModule{})
	models.Models = append(models.Models, CarrierModuleDomain{})
	models.Models = append(models.Models, ModificationModuleDomain{})
	models.Models = append(models.Models, ATModuleDomain{})
	models.Models = append(models.Models, KSModuleDomain{})
	models.Models = append(models.Models, BiosyntheticModuleDomainLocation{})
}
