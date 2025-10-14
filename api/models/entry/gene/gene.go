package gene

import "github.com/adraismawur/mibig-submission/models"

// TODO: there may be a way to unify this, locus.location and BiosyntheticModuleDomainLocation
type ExonLocation struct {
	GeneLocationID uint64 `json:"-"`
	From           uint64 `json:"from"`
	To             uint64 `json:"to"`
}

type GeneLocation struct {
	ID     uint64         `json:"-"`
	GeneID uint64         `json:"-"`
	Exons  []ExonLocation `json:"exons" gorm:"ForeignKey:GeneLocationID"`
	Strand int32          `json:"strand"`
}

type GeneAddition struct {
	GeneID      uint64       `json:"-"`
	Accession   string       `json:"id"`
	Location    GeneLocation `json:"location" gorm:"ForeignKey:GeneID"`
	Translation string       `json:"translation"`
}

type GeneAnnotation struct {
	GeneID    uint64 `json:"-"`
	Accession string `json:"id"`
	Name      string `json:"name"`
	Product   string `json:"product"`
}

type Gene struct {
	ID          uint64           `json:"-"`
	EntryID     uint64           `json:"-"`
	Additions   []GeneAddition   `json:"to_add,omitempty" gorm:"ForeignKey:GeneID"`
	Annotations []GeneAnnotation `json:"annotations,omitempty" gorm:"ForeignKey:GeneID"`
}

func init() {
	models.Models = append(models.Models, &Gene{})
	models.Models = append(models.Models, &GeneAddition{})
	models.Models = append(models.Models, &GeneLocation{})
	models.Models = append(models.Models, &ExonLocation{})
	models.Models = append(models.Models, &GeneAnnotation{})
}
