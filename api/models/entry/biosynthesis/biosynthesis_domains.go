package biosynthesis

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Acetyltransferase struct {
	ID                  uint64                    `json:"db_id"`
	BiosyntheticClassID uint64                    `json:"biosynthetic_class_id"`
	Type                string                    `json:"type"`
	Subtype             string                    `json:"subtype"`
	Gene                pq.StringArray            `json:"gene" gorm:"type:text[]"`
	LocationID          uint64                    `json:"db_location_id"`
	Location            DomainLocation            `json:"location"`
	Inactive            bool                      `json:"inactive,omitempty"`
	Substrates          []DomainSubstrate         `json:"substrates,omitempty" gorm:"many2many:acetyltransferase_substrates;"`
	Evidence            []DomainSubstrateEvidence `json:"evidence,omitempty" gorm:"many2many:acetyltransferase_evidences;"`
}

type Adenylation struct {
	ID                    uint64                     `json:"db_id"`
	BiosyntheticClassID   uint64                     `json:"biosynthetic_class_id"`
	Type                  string                     `json:"type"`
	Gene                  pq.StringArray             `json:"gene" gorm:"type:text[]"`
	LocationID            uint64                     `json:"db_location_id"`
	Location              DomainLocation             `json:"location"`
	Inactive              bool                       `json:"inactive"`
	Evidence              *[]DomainSubstrateEvidence `json:"evidence" gorm:"many2many:adenylation_evidences;"`
	PrecursorBiosynthesis pq.StringArray             `json:"precursor_biosynthesis" gorm:"type:text[]"`
	Substrates            *[]DomainSubstrate         `json:"substrates,omitempty" gorm:"many2many:adenylation_substrates;"`
}

type Thioesterase struct {
	ID                  uint64         `json:"db_id"`
	BiosyntheticClassID uint64         `json:"biosynthetic_class_id"`
	Type                string         `json:"type"`
	LocationID          uint64         `json:"db_location_id"`
	Location            DomainLocation `json:"location"`
	Subtype             string         `json:"subtype,omitempty"`
}

type DomainLocation struct {
	ID   uint64 `json:"db_id"`
	From int    `json:"from"`
	To   int    `json:"to"`
}

type DomainSubstrate struct {
	gorm.Model
	Name      string `json:"name"`
	Details   string `json:"details,omitempty"`
	Structure string `json:"structure,omitempty"`
}

type DomainSubstrateEvidence struct {
	ID         uint64         `json:"db_id""`
	Method     string         `json:"method"`
	References pq.StringArray `json:"references" gorm:"type:text[]"`
}
