package biosynthesis

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/lib/pq"
)

type BiosyntheticPathwayProduct struct {
	ID                    uint64 `json:"db_id"`
	BiosyntheticPathwayID uint64 `json:"db_biosynthetic_pathway_id"`
	Name                  string `json:"name"`
	Structure             string `json:"structure"`
	Comment               string `json:"comment"`
}

type BiosyntheticPathway struct {
	ID                uint64                       `json:"db_id"`
	BiosynthesisID    uint64                       `json:"db_biosynth_id"`
	Products          []BiosyntheticPathwayProduct `json:"products" gorm:"foreignKey:BiosyntheticPathwayID"`
	Steps             string                       `json:"steps"`
	References        pq.StringArray               `json:"references" gorm:"type:text[]"`
	IsSubCluster      bool                         `json:"is_sub_cluster"`
	ProducesPrecursor bool                         `json:"produces_precursor"`
}

func init() {
	models.Models = append(models.Models, &BiosyntheticPathway{})
	models.Models = append(models.Models, &BiosyntheticPathwayProduct{})
}
