package submission

type Location struct {
	From int64 `json:"from"`
	To   int64 `json:"to"`
}

type EvidenceMethod struct {
	ID   int64
	Name string
}

type EvidenceCitation struct {
	EvidenceId int64
	CitationId int64
}

type Evidence struct {
	ID     int64          `json:"-"`
	Method EvidenceMethod `json:"method"`
}

type Locus struct {
	ID        int64      `json:"id"`
	Accession string     `json:"accession"`
	Location  Location   `json:"location"`
	Evidence  []Evidence `json:"evidence"`
}
