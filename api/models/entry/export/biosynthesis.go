package export

type Biosynthesis struct {
	Classes []BiosyntheticClass   `json:"classes" gorm:"foreignKey:BiosynthesisID"`
	Modules []BiosyntheticModule  `json:"modules,omitempty" gorm:"foreignKey:BiosynthesisID"`
	Operons []BiosyntheticOperon  `json:"operons,omitempty" gorm:"foreignKey:BiosynthesisID"`
	Paths   []BiosyntheticPathway `json:"paths,omitempty" gorm:"foreignKey:BiosynthesisID"`
}
