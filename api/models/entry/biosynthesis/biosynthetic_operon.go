package biosynthesis

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/lib/pq"
)

type BiosyntheticOperon struct {
	ID             uint64         `json:"db_id"`
	BiosynthesisID uint64         `json:"db_biosynth_id"`
	Items          pq.StringArray `json:"items" gorm:"type:text[]"`
}

func init() {
	models.Models = append(models.Models, &BiosyntheticOperon{})
}
