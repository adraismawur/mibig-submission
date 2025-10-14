package taxonomy

import "github.com/adraismawur/mibig-submission/models"

type Taxonomy struct {
	EntryID uint64 `json:"-"`
	Name    string `json:"name"`
	TaxID   uint64 `json:"ncbiTaxId"`
}

func init() {
	models.Models = append(models.Models, &Taxonomy{})
}
