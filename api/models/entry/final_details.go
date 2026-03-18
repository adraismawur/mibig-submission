package entry

import (
	"github.com/adraismawur/mibig-submission/models/entry/consts"
	"gorm.io/gorm"
)

// FinalDetails describes whether an entry contains all genes responsible for
// production of components (completeness) and whether it is under embargo
type FinalDetails struct {
	Accession    string              `json:"accession" gorm:"primaryKey"`
	Completeness consts.Completeness `json:"completeness"`
	Embargo      bool                `json:"embargo"`
}

func UpdateFinalDetails(db *gorm.DB, details FinalDetails) error {
	err := db.
		Table("entries").
		Model(&FinalDetails{}).
		Where("accession = ?", details.Accession).
		Save(details).
		Error

	return err
}
