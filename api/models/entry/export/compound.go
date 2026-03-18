package export

import "github.com/lib/pq"

type BioActivities struct {
	Name       *string        `json:"name,omitempty"`
	Observed   bool           `json:"observed"`
	References pq.StringArray `json:"references" gorm:"type:text[]"`
}

type CompoundEvidence struct {
	Method     string         `json:"method"`
	References pq.StringArray `json:"references" gorm:"type:text[]"`
}

type Compound struct {
	Name          string             `json:"name"`
	Evidence      []CompoundEvidence `json:"evidence" gorm:"foreignKey:CompoundID"`
	BioActivities []BioActivities    `json:"bioactivities,omitempty" gorm:"foreignKey:CompoundID"`
	Structure     string             `json:"structure"`
	DatabaseIDs   pq.StringArray     `json:"databaseIds" gorm:"type:text[]"`
	Moieties      pq.StringArray     `json:"moieties,omitempty" gorm:"type:text[]"`
	Mass          float64            `json:"mass"`
	Formula       string             `json:"formula"`
}
