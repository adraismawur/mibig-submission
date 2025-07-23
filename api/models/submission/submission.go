package submission

import "github.com/adraismawur/mibig-submission/models"

type Quality string

const (
	Questionable Quality = "questionable"
	Medium       Quality = "medium"
	High         Quality = "high"
)

type Status string

const (
	Pending Status = "pending"
	Active  Status = "active"
	Retired Status = "retired"
)

type Completeness string

const (
	Unknown  Completeness = "unknown"
	Complete Completeness = "complete"
)

type Submission struct {
	ID           uint         `json:"-"`
	Accession    string       `json:"accession"`
	Version      int          `json:"version"`
	Quality      Quality      `json:"quality"`
	Status       Status       `json:"status"`
	Completeness Completeness `json:"completeness"`
}

func init() {
	models.Models = append(models.Models, &Submission{})
}
