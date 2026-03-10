package locus

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/lib/pq"
)

type Location struct {
	ID      int64  `json:"db_id"`
	LocusID int64  `json:"db_locus_id"`
	Start   *int64 `json:"from"`
	End     *int64 `json:"to"`
}

type LocusEvidence struct {
	ID         int64          `json:"db_id"`
	LocusID    int64          `json:"db_locus_id"`
	Method     string         `json:"method"`
	References pq.StringArray `json:"references,omitempty" gorm:"type:text[]"`
}

type Locus struct {
	ID          int64           `json:"db_id"`
	EntryID     int64           `json:"db_entry_id"`
	Accession   string          `json:"accession"`
	Location    Location        `json:"location" gorm:"foreignKey:LocusID"`
	Evidence    []LocusEvidence `json:"evidence" gorm:"foreignKey:LocusID"`
	DraftGenome bool            `json:"draft_genome,omitempty"`
}

func init() {
	models.Models = append(models.Models, &Locus{})
	models.Models = append(models.Models, &Location{})
	models.Models = append(models.Models, &LocusEvidence{})
}
