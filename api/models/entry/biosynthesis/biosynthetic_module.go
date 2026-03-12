package biosynthesis

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/lib/pq"
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
