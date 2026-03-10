package export

import "github.com/lib/pq"

type BiosyntheticClass struct {
	ID             uint64         `json:"-"`
	BiosynthesisID uint64         `json:"-"`
	Class          string         `json:"class"`
	Subclass       string         `json:"subclass"`
	Cyclases       pq.StringArray `json:"cyclases" gorm:"type:text[]"`
	References     pq.StringArray `json:"references" gorm:"type:text[]"`
}
