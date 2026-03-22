package export

import "github.com/lib/pq"

type IntegratedMonomer struct {
	Name      string                    `json:"name"`
	Structure string                    `json:"structure"`
	Evidence  []DomainSubstrateEvidence `json:"evidence" gorm:"many2many:integrated_monomer_evidences"`
}

type BiosyntheticModule struct {
	Type                string                   `json:"type"`
	Name                string                   `json:"name"`
	Genes               pq.StringArray           `json:"genes" gorm:"type:text[]"`
	Active              bool                     `json:"active"`
	IntegratedMonomers  []IntegratedMonomer      `json:"integrated_monomers" gorm:"foreignKey:BiosyntheticModuleID"`
	Carriers            []CarrierDomain          `json:"carriers" gorm:"many2many:biosynth_carrier_domains"`
	ModificationDomains []ModificationDomain     `json:"modification_domains,omitempty" gorm:"many2many:biosynth_modification_domains"`
	CDomain             *CondensationDomain      `json:"c_domain"`
	ADomain             *AdenylationDomain       `json:"a_domain,omitempty"`
	ATDomain            *AcetyltransferaseDomain `json:"at_domain,omitempty"`
	KSDomain            *KetoSynthaseDomain      `json:"ks_domain,omitempty"`
}
