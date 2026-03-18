package export

type ExonLocation struct {
	From uint64 `json:"from"`
	To   uint64 `json:"to"`
}

type GeneLocation struct {
	Exons  []ExonLocation `json:"exons" gorm:"ForeignKey:GeneLocationID"`
	Strand int32          `json:"strand"`
}

type GeneAddition struct {
	Accession   string       `json:"accession"`
	Location    GeneLocation `json:"location" gorm:"ForeignKey:GeneID"`
	Translation string       `json:"translation"`
}

type GeneDeletion struct {
	Accession string `json:"accession"`
	Reason    string `json:"reason"`
}

type GeneAnnotation struct {
	Accession string `json:"accession"` // Accession is the gene ID, e.g. 'AEK75497.1'. This is confusing, but GeneID here is internal to the API
	Name      string `json:"name"`      // Name is the actual gene name, e.g. 'abyA1'
	Product   string `json:"product"`   // Product is the product of this gene, e.g. '3-oxoacyl-ACP synthase III'
}

type Gene struct {
	Additions   []GeneAddition   `json:"to_add,omitempty" gorm:"ForeignKey:GeneID"`
	Deletions   []GeneDeletion   `json:"to_delete,omitempty" gorm:"ForeignKey:GeneID"`
	Annotations []GeneAnnotation `json:"annotations,omitempty" gorm:"ForeignKey:GeneID"`
}
