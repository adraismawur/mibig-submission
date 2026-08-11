package export

import "github.com/lib/pq"

type BiosyntheticOperon struct {
	Genes    pq.StringArray `json:"genes" gorm:"type:text[]"`
	Evidence pq.StringArray `json:"evidence" gorm:"type:text[]"`
}
