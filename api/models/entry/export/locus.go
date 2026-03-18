package export

import "github.com/lib/pq"

type Location struct {
	Start *int64 `json:"from"`
	End   *int64 `json:"to"`
}

type LocusEvidence struct {
	Method     string         `json:"method"`
	References pq.StringArray `json:"references,omitempty" gorm:"type:text[]"`
}

type Locus struct {
	Accession   string          `json:"accession"`
	Location    Location        `json:"location" gorm:"foreignKey:LocusID"`
	Evidence    []LocusEvidence `json:"evidence" gorm:"foreignKey:LocusID"`
	DraftGenome bool            `json:"draft_genome,omitempty"`
}
