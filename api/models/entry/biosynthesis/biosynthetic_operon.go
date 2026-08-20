package biosynthesis

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/lib/pq"
)

type BiosyntheticOperonEvidence struct {
	ID         uint64         `json:"db_id"`
	Method     string         `json:"method"`
	References pq.StringArray `json:"references" gorm:"type:text[]"`
}

type BiosyntheticOperon struct {
	ID             uint64                       `json:"db_id"`
	BiosynthesisID uint64                       `json:"db_biosynth_id"`
	Genes          pq.StringArray               `json:"genes" gorm:"type:text[]"`
	Evidence       []BiosyntheticOperonEvidence `json:"evidence" gorm:"many2many:biosynthetic_operon_evidences"`
}

func init() {
	models.Models = append(models.Models, &BiosyntheticOperonEvidence{})
	models.Models = append(models.Models, &BiosyntheticOperon{})
}
