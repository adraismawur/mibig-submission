package biosynthesis

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

func init() {
	models.Models = append(models.Models, BiosyntheticClass{})
}
