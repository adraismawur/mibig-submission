package biosynthesis

import (
	"github.com/adraismawur/mibig-submission/models"
)

type Biosynthesis struct {
	ID      uint64               `json:"-"`
	EntryID uint64               `json:"-"`
	Classes []BiosyntheticClass  `json:"classes" gorm:"foreignKey:BiosynthesisID"`
	Modules []BiosyntheticModule `json:"modules" gorm:"foreignKey:BiosynthesisID"`
}

func init() {
	models.Models = append(models.Models, Biosynthesis{})
}
