package entry

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/lib/pq"
)

type BiosyntheticClass struct {
	BiosynthesisID uint64         `json:"-"`
	Class          string         `json:"class"`
	Subclass       string         `json:"subclass"`
	Cyclases       pq.StringArray `json:"cyclases" gorm:"type:text[]"`
}

type Biosynthesis struct {
	ID      uint64              `json:"-"`
	EntryID uint64              `json:"-"`
	Classes []BiosyntheticClass `json:"classes" gorm:"foreignKey:BiosynthesisID"`
}

func init() {
	models.Models = append(models.Models, Biosynthesis{})
	models.Models = append(models.Models, BiosyntheticClass{})
}
