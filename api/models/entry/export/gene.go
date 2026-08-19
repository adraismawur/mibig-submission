package export

import "github.com/lib/pq"

type ExonLocation struct {
	From uint64 `json:"from"`
	To   uint64 `json:"to"`
}

type GeneLocation struct {
	Exons  []ExonLocation `json:"exons" gorm:"ForeignKey:GeneLocationID"`
	Strand int32          `json:"strand"`
}

type GeneAddition struct {
	Accession   string       `json:"id"`
	Location    GeneLocation `json:"location" gorm:"ForeignKey:GeneID"`
	Translation string       `json:"translation"`
}

type GeneDeletion struct {
	Accession string `json:"id"`
	Reason    string `json:"reason"`
}

type GeneFunctionAnnotationEvidence struct {
	Method     string         `json:"method"`
	References pq.StringArray `json:"references" gorm:"type:text[]"`
}

type GeneMutationPhenotypeAnnotation struct {
	Phenotype  string         `json:"phenotype"`
	Details    string         `json:"details"`
	References pq.StringArray `json:"references" gorm:"type:text[]"`
}

type GeneFunctionAnnotationFunction struct {
	Name    string `json:"name"`
	Details string `json:"details"`
}

type GeneFunctionAnnotation struct {
	Function          *GeneFunctionAnnotationFunction  `json:"function"`
	Details           string                           `json:"details"`
	Evidence          []GeneFunctionAnnotationEvidence `json:"evidence" gorm:"foreignKey:GeneFunctionAnnotationID"`
	MutationPhenotype *GeneMutationPhenotypeAnnotation `json:"mutation_phenotype"`
}

type GeneAnnotation struct {
	Accession string `json:"id"`      // Accession is the gene ID, e.g. 'AEK75497.1'. This is confusing, but GeneID here is internal to the API
	Name      string `json:"name"`    // Name is the actual gene name, e.g. 'abyA1'
	Product   string `json:"product"` // Product is the product of this gene, e.g. '3-oxoacyl-ACP synthase III'
}

type GeneInformation struct {
	Additions   []GeneAddition   `json:"to_add,omitempty" gorm:"ForeignKey:GeneID"`
	Deletions   []GeneDeletion   `json:"to_delete,omitempty" gorm:"ForeignKey:GeneID"`
	Annotations []GeneAnnotation `json:"annotations,omitempty" gorm:"ForeignKey:GeneID"`
}
