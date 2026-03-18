package export

import "github.com/lib/pq"

type BiosyntheticOperon struct {
	Items pq.StringArray `json:"items" gorm:"type:text[]"`
}
