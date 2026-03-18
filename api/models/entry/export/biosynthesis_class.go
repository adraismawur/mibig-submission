package export

import "github.com/lib/pq"

type ReleaseType struct {
	Name       string         `json:"name"`
	Details    string         `json:"details,omitempty"`
	References pq.StringArray `json:"references,omitempty" gorm:"type:text[]"`
}

type Thioesterase struct {
	Type     string         `json:"type"`
	Gene     string         `json:"gene"`
	Location DomainLocation `json:"location"`
	Subtype  string         `json:"subtype,omitempty"`
}

type CleavageLocation struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type RippPrecursorCrosslink struct {
	From    int    `json:"from"`
	To      int    `json:"to"`
	Type    string `json:"type"`
	Details string `json:"details"`
}

type RippPrecursor struct {
	Gene string `json:"gene""`
	//CoreSequence             string                   `json:"core_sequence"`
	LeaderCleavageLocation   *CleavageLocation        `json:"leader_cleavage_location,omitempty"`
	FollowerCleavageLocation *CleavageLocation        `json:"follower_cleavage_location,omitempty"`
	Crosslinks               []RippPrecursorCrosslink `json:"crosslinks,omitempty" gorm:"foreignKey:RippPrecursorID"`
}

type GlycosylTransferase struct {
	Evidence    []DomainSubstrateEvidence `json:"evidence,omitempty" gorm:"many2many:glycosyl_transferase_evidences;"`
	References  pq.StringArray            `json:"references" gorm:"type:text[]"`
	Gene        string                    `json:"gene"`
	Specificity string                    `json:"specificity"`
}

type SaccharideSubcluster struct {
	Specificity string         `json:"specificity"`
	Genes       pq.StringArray `json:"genes" gorm:"type:text[]"`
	References  pq.StringArray `json:"references" gorm:"type:text[]"`
}

type BiosyntheticClass struct {
	// common elements
	Class    string `json:"class"`
	Subclass string `json:"subclass"`

	// NRPS
	Thioesterases []Thioesterase `json:"thioesterases,omitempty" gorm:"foreignKey:BiosyntheticClassID"`
	ReleaseTypes  []ReleaseType  `json:"release_types,omitempty" gorm:"foreignKey:BiosyntheticClassID"`

	// PKS
	Cyclases     pq.StringArray `json:"cyclases" gorm:"type:text[]"`
	KetideLength *int           `json:"ketideLength,omitempty"`
	// TODO: starter unit?
	// TODO: iterative?

	// ribosomal
	RIPPType   *string         `json:"ripp_type,omitempty"`
	Details    *string         `json:"details,omitempty"`
	Peptidases pq.StringArray  `json:"peptidases,omitempty" gorm:"type:text[]"`
	Precursors []RippPrecursor `json:"precursors,omitempty" gorm:"foreignKey:BiosyntheticClassID"`

	// saccharide
	GlycosylTransferases []GlycosylTransferase  `json:"glycosyltransferases,omitempty" gorm:"foreignKey:BiosyntheticClassID"`
	Subclusters          []SaccharideSubcluster `json:"subclusters,omitempty" gorm:"foreignKey:BiosyntheticClassID"`

	// terpene
	Prenyltransferases pq.StringArray `json:"prenyltransferases,omitempty" gorm:"type:text[]"`
	SynthasesCyclases  pq.StringArray `json:"synthases_cyclases,omitempty" gorm:"type:text[]"`
	Precursor          *string        `json:"precursor,omitempty"`
}
