package export

type Biosynthesis struct {
	ID      uint64               `json:"-"`
	EntryID uint64               `json:"-"`
	Classes []BiosyntheticClass  `json:"classes" gorm:"foreignKey:BiosynthesisID"`
	Modules []BiosyntheticModule `json:"modules,omitempty" gorm:"foreignKey:BiosynthesisID"`
}
