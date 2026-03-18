package export

type Taxonomy struct {
	Name  string `json:"name"`
	TaxID uint64 `json:"ncbiTaxId"`
}
