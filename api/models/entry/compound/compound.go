package compound

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/lib/pq"
)

type AssayMeasurement struct {
	ID            uint64  `json:"db_id"`
	Concentration float64 `json:"concentration"`
	Unit          string  `json:"unit"`
	Error         float64 `json:"error"`
	Replicates    float64 `json:"replicates"`
}

type AssayTestSystem struct {
	ID       uint64 `json:"db_id"`
	CellLine string `json:"cell_line"`
	Organism int64  `json:"organism"`
	Strain   string `json:"strain"`
}

type BioActivityAssay struct {
	ID            uint64           `json:"db_id"`
	BioActivityID uint64           `json:"db_bio_activity_id"`
	MeasurementID uint64           `json:"db_measurement_id"`
	Measurement   AssayMeasurement `json:"measurement"`
	Target        string           `json:"target"`
	Details       string           `json:"details"`
	TestSystemID  uint64           `json:"db_test_system_id"`
	TestSystem    AssayTestSystem  `json:"test_system"`
	References    pq.StringArray   `json:"references" gorm:"type:text[]"`
}

type BioActivities struct {
	ID         uint64             `json:"db_id"`
	CompoundID uint64             `json:"db_compound_id"`
	Name       string             `json:"name,omitempty"`
	Details    string             `json:"details,omitempty"`
	Observed   bool               `json:"observed"`
	Assays     []BioActivityAssay `json:"assays" gorm:"foreignKey:BioActivityID"`
	References pq.StringArray     `json:"references" gorm:"type:text[]"`
}

type CompoundEvidence struct {
	ID         uint64         `json:"db_id"`
	CompoundID uint64         `json:"db_compound_id"`
	Method     string         `json:"method"`
	References pq.StringArray `json:"references" gorm:"type:text[]"`
}

type Compound struct {
	ID             uint64             `json:"db_id"`
	EntryAccession string             `json:"db_entry_accession"`
	Name           string             `json:"name"`
	Evidence       []CompoundEvidence `json:"evidence" gorm:"foreignKey:CompoundID"`
	BioActivities  []BioActivities    `json:"bioactivities,omitempty" gorm:"foreignKey:CompoundID"`
	Structure      string             `json:"structure"`
	DatabaseIDs    pq.StringArray     `json:"databaseIds" gorm:"type:text[]"`
	Moieties       pq.StringArray     `json:"moieties,omitempty" gorm:"type:text[]"`
	Mass           float64            `json:"mass"`
	Formula        string             `json:"formula"`
}

func init() {
	models.Models = append(models.Models, &Compound{})
	models.Models = append(models.Models, &BioActivities{})
	models.Models = append(models.Models, &CompoundEvidence{})
	models.Models = append(models.Models, &AssayMeasurement{})
	models.Models = append(models.Models, &AssayTestSystem{})
	models.Models = append(models.Models, &BioActivityAssay{})
}
