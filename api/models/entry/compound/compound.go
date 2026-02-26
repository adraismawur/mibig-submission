package compound

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/lib/pq"
)

type BioActivities struct {
	CompoundID uint64 `json:"-"`
	//Name       string         `json:"name"`
	Observed   bool           `json:"observed"`
	References pq.StringArray `json:"references" gorm:"type:text[]"`
}

type CompoundEvidence struct {
	CompoundID uint64         `json:"-"`
	Method     string         `json:"method"`
	References pq.StringArray `json:"references" gorm:"type:text[]"`
}

type Compound struct {
	ID            uint64             `json:"id"`
	EntryID       uint64             `json:"-"`
	Name          string             `json:"name"`
	Evidence      []CompoundEvidence `json:"evidence" gorm:"foreignKey:CompoundID"`
	BioActivities []BioActivities    `json:"bioactivities,omitempty" gorm:"foreignKey:CompoundID"`
	Structure     string             `json:"structure"`
	DatabaseIDs   pq.StringArray     `json:"databaseIds" gorm:"type:text[]"`
	Moieties      pq.StringArray     `json:"moieties,omitempty" gorm:"type:text[]"`
	Mass          float64            `json:"mass"`
	Formula       string             `json:"formula"`
}

func init() {
	models.Models = append(models.Models, &Compound{})
	models.Models = append(models.Models, &BioActivities{})
	models.Models = append(models.Models, &CompoundEvidence{})
}
