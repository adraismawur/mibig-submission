package export

import (
	"github.com/adraismawur/mibig-submission/models/entry/consts"
	"github.com/lib/pq"
)

type Entry struct {
	Accession        string              `json:"accession"`
	Version          int                 `json:"version,omitempty"`
	Changelog        Changelog           `json:"changelog" gorm:"foreignKey:EntryAccession"`
	Quality          consts.Quality      `json:"quality,omitempty"`
	Status           consts.Status       `json:"status,omitempty"`
	Completeness     consts.Completeness `json:"completeness"`
	Loci             []Locus             `json:"loci" gorm:"foreignKey:EntryAccession"`
	Biosynthesis     Biosynthesis        `json:"biosynthesis" gorm:"foreignKey:EntryAccession"`
	Compounds        []Compound          `json:"compounds" gorm:"ForeignKey:EntryAccession"`
	Taxonomy         Taxonomy            `json:"taxonomy" gorm:"ForeignKey:EntryAccession"`
	GeneInformation  *GeneInformation    `json:"genes,omitempty" gorm:"ForeignKey:EntryAccession"`
	LegacyReferences pq.StringArray      `json:"legacy_references,omitempty" gorm:"type:text[]"`
}
