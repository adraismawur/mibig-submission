package export

type ExonLocation struct {
	ID             uint64 `json:"-"`
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
	ID          uint64       `json:"-"`
	GeneID      uint64       `json:"-"`
	Accession   string       `json:"accession"`
	Location    GeneLocation `json:"location" gorm:"ForeignKey:GeneID"`
	Translation string       `json:"translation"`
}

type GeneDeletion struct {
	ID        uint64 `json:"-"`
	GeneID    uint64 `json:"-"`
	Accession string `json:"accession"`
	Reason    string `json:"reason"`
}

type GeneAnnotation struct {
	ID        uint64 `json:"-"`
	GeneID    uint64 `json:"-"`         // GeneID is an internal identifier for the API DB
	Accession string `json:"accession"` // Accession is the gene ID, e.g. 'AEK75497.1'. This is confusing, but GeneID here is internal to the API
	Name      string `json:"name"`      // Name is the actual gene name, e.g. 'abyA1'
	Product   string `json:"product"`   // Product is the product of this gene, e.g. '3-oxoacyl-ACP synthase III'
}

type Gene struct {
	ID          uint64           `json:"-"`
	EntryID     uint64           `json:"-"`
	Additions   []GeneAddition   `json:"to_add,omitempty" gorm:"ForeignKey:GeneID"`
	Deletions   []GeneDeletion   `json:"to_delete,omitempty" gorm:"ForeignKey:GeneID"`
	Annotations []GeneAnnotation `json:"annotations,omitempty" gorm:"ForeignKey:GeneID"`
}
