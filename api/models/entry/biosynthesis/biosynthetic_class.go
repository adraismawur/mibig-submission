package biosynthesis

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/lib/pq"
)

type ReleaseType struct {
	ID                  uint64         `json:"db_id"`
	BiosyntheticClassID uint64         `json:"db_biosynthetic_class_id"`
	Name                string         `json:"name"`
	Details             string         `json:"details,omitempty"`
	References          pq.StringArray `json:"references,omitempty" gorm:"type:text[]"`
}

type CleaveageLocation struct {
	ID    uint64 `json:"db_id"`
	Start int    `json:"start"`
	To    int    `json:"to"`
}

type RippPrecursorCrosslink struct {
	ID              uint64 `json:"db_id"`
	RippPrecursorID uint64 `json:"ripp_precursor_id"`
	From            int    `json:"from"`
	To              int    `json:"to"`
	Type            string `json:"type"`
	Details         string `json:"details"`
}

type RippPrecursor struct {
	ID                  uint64 `json:"db_id"`
	BiosyntheticClassID uint64 `json:"biosynthetic_class_id"`
	Gene                string `json:"gene""`
	//CoreSequence             string                   `json:"core_sequence"`
	//LeaderCleavageLocation *CleaveageLocation `json:"leader_cleavage_location" gorm:"many2many:ripp_leader_locations;"`
	//FollowerCleavageLocation *CleaveageLocation       `json:"follower_cleavage_location" gorm:"many2many:ripp_follower_locations;"`
	Crosslinks []RippPrecursorCrosslink `json:"crosslinks,omitempty" gorm:"foreignKey:RippPrecursorID"`
}

type BiosyntheticClass struct {
	ID                uint64             `json:"db_id"`
	BiosynthesisID    uint64             `json:"db_biosynth_id"`
	Class             string             `json:"class"`
	Subclass          string             `json:"subclass"`
	Cyclases          pq.StringArray     `json:"cyclases" gorm:"type:text[]"`
	ReleaseType       *ReleaseType       `json:"release_type,omitempty" gorm:"foreignKey:BiosyntheticClassID"`
	AcetylTransferase *Acetyltransferase `json:"acetyltransferase,omitempty" gorm:"foreignKey:BiosyntheticClassID"`
	Adenylation       *Adenylation       `json:"adenylation,omitempty" gorm:"foreignKey:BiosyntheticClassID"`
	Thioesterase      *[]Thioesterase    `json:"thioesterase,omitempty" gorm:"foreignKey:BiosyntheticClassID"`
	Details           *string            `json:"details,omitempty"`
	KetideLength      *int               `json:"ketideLength,omitempty"`
	RippType          *string            `json:"ripp_type,omitempty"`
	Peptidases        pq.StringArray     `json:"peptidases,omitempty" gorm:"type:text[]"`
	Precursors        *[]RippPrecursor   `json:"precursors,omitempty" gorm:"foreignKey:BiosyntheticClassID"`
}

func init() {
	models.Models = append(models.Models, DomainSubstrate{})
	models.Models = append(models.Models, DomainSubstrateEvidence{})
	models.Models = append(models.Models, Acetyltransferase{})
	models.Models = append(models.Models, ReleaseType{})
	models.Models = append(models.Models, BiosyntheticClass{})
	models.Models = append(models.Models, Adenylation{})
	models.Models = append(models.Models, CleaveageLocation{})
	models.Models = append(models.Models, RippPrecursorCrosslink{})
	models.Models = append(models.Models, RippPrecursor{})
	models.Models = append(models.Models, Thioesterase{})
}
