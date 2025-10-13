package entry

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/lib/pq"
)

type Location struct {
	LocusID int64  `json:"locus_id"`
	Start   *int64 `json:"from"`
	End     *int64 `json:"to"`
}

type Evidence struct {
	LocusID    int64          `json:"locus_id"`
	Method     string         `json:"method"`
	References pq.StringArray `json:"references" gorm:"type:text[]"`
}

type Locus struct {
	ID          int64      `json:"-"`
	EntryID     int64      `json:"entry_id"`
	Accession   string     `json:"accession"`
	Location    Location   `json:"location" gorm:"foreignKey:LocusID"`
	Evidence    []Evidence `json:"evidence" gorm:"foreignKey:LocusID"`
	DraftGenome bool       `json:"draft_genome,omitempty"`
}

func init() {
	models.Models = append(models.Models, &Locus{})
	models.Models = append(models.Models, &Location{})
	models.Models = append(models.Models, &Evidence{})
}
