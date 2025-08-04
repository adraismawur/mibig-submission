package entry

import "time"

type Citation struct {
	ID              int64 `json:"-"`
	Title           string
	PublicationDate time.Time
	DOI             string
}
