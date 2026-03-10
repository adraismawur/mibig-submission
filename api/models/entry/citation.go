package entry

import (
	"time"
)

type Citation struct {
	ID uint64 `json:"db_id"`
	ExportCitation
}

type ExportCitation struct {
	Title           string
	PublicationDate time.Time
	DOI             string
}
