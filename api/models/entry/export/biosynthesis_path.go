package export

import "github.com/lib/pq"

type BiosyntheticPathwayProduct struct {
	Name      string `json:"name"`
	Structure string `json:"structure"`
	Comment   string `json:"comment"`
}

type BiosyntheticPathway struct {
	Products          []BiosyntheticPathwayProduct `json:"products" gorm:"foreignKey:BiosyntheticPathwayID"`
	Steps             string                       `json:"steps"`
	References        pq.StringArray               `json:"references" gorm:"type:text[]"`
	IsSubCluster      bool                         `json:"is_sub_cluster"`
	ProducesPrecursor bool                         `json:"produces_precursor"`
}
