package taxonomy

import "github.com/adraismawur/mibig-submission/models"

type Taxonomy struct {
	ID      uint64 `json:"db_id"`
	EntryID uint64 `json:"entry_id"`
	Name    string `json:"name"`
	TaxID   uint64 `json:"ncbiTaxId"`
}

func init() {
	models.Models = append(models.Models, &Taxonomy{})
}
