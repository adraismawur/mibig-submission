package export

import (
	"github.com/adraismawur/mibig-submission/models/entry/consts"
	"github.com/lib/pq"
)

type Entry struct {
	ID               uint64              `json:"-"`
	Accession        string              `json:"accession"`
	Version          int                 `json:"version,omitempty"`
	Changelog        Changelog           `json:"changelog" gorm:"foreignKey:EntryID"`
	Quality          consts.Quality      `json:"quality,omitempty"`
	Status           consts.Status       `json:"status,omitempty"`
	Completeness     consts.Completeness `json:"completeness"`
	Loci             []Locus             `json:"loci" gorm:"foreignKey:EntryID"`
	Biosynthesis     Biosynthesis        `json:"biosynthesis" gorm:"foreignKey:EntryID"`
	Compounds        []Compound          `json:"compounds" gorm:"ForeignKey:EntryID"`
	Taxonomy         Taxonomy            `json:"taxonomy" gorm:"ForeignKey:EntryID"`
	Genes            *Gene               `json:"genes,omitempty" gorm:"ForeignKey:EntryID"`
	LegacyReferences pq.StringArray      `json:"legacy_references,omitempty" gorm:"type:text[]"`
}
