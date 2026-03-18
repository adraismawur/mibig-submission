package biosynthesis

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ReleaseType struct {
	ID                  uint64         `json:"db_id"`
	BiosyntheticClassID uint64         `json:"db_class_id"`
	Name                string         `json:"name"`
	Details             string         `json:"details,omitempty"`
	References          pq.StringArray `json:"references,omitempty" gorm:"type:text[]"`
}

type Thioesterase struct {
	ID                  uint64         `json:"db_id"`
	BiosyntheticClassID uint64         `json:"db_biosynth_class_id"`
	Type                string         `json:"type"`
	Gene                string         `json:"gene"`
	LocationID          uint64         `json:"db_location_id"`
	Location            DomainLocation `json:"location"`
	Subtype             string         `json:"subtype,omitempty"`
}

type CleavageLocation struct {
	ID   uint64 `json:"db_id"`
	From int    `json:"from"`
	To   int    `json:"to"`
}

type RippPrecursorCrosslink struct {
	ID              uint64 `json:"db_id"`
	RippPrecursorID uint64 `json:"db_ripp_precursor_id"`
	From            int    `json:"from"`
	To              int    `json:"to"`
	Type            string `json:"type"`
	Details         string `json:"details"`
}

type RippPrecursor struct {
	ID                  uint64 `json:"db_id"`
	BiosyntheticClassID uint64 `json:"db_class_id"`
	Gene                string `json:"gene""`
	//CoreSequence             string                   `json:"core_sequence"`
	LeaderCleavageLocationID   uint64                   `json:"db_leader_cleavage_location_id"`
	LeaderCleavageLocation     *CleavageLocation        `json:"leader_cleavage_location,omitempty"`
	FollowerCleavageLocationID uint64                   `json:"db_follower_cleavage_location_id"`
	FollowerCleavageLocation   *CleavageLocation        `json:"follower_cleavage_location,omitempty"`
	Crosslinks                 []RippPrecursorCrosslink `json:"crosslinks,omitempty" gorm:"foreignKey:RippPrecursorID"`
}

type GlycosylTransferase struct {
	ID                  uint64                    `json:"db_id"`
	BiosyntheticClassID uint64                    `json:"db_biosynth_class_id"`
	Evidence            []DomainSubstrateEvidence `json:"evidence,omitempty" gorm:"many2many:glycosyl_transferase_evidences;"`
	References          pq.StringArray            `json:"references" gorm:"type:text[]"`
	Gene                string                    `json:"gene"`
	Specificity         string                    `json:"specificity"`
}

type SaccharideSubcluster struct {
	ID                  uint64         `json:"db_id"`
	BiosyntheticClassID uint64         `json:"db_class_id"`
	Specificity         string         `json:"specificity"`
	Genes               pq.StringArray `json:"genes" gorm:"type:text[]"`
	References          pq.StringArray `json:"references" gorm:"type:text[]"`
}

type BiosyntheticClass struct {
	ID             uint64 `json:"db_id"`
	BiosynthesisID uint64 `json:"db_biosynth_id"`
	// common elements
	Class    string `json:"class"`
	Subclass string `json:"subclass"`

	// NRPS
	Thioesterases []Thioesterase `json:"thioesterases,omitempty" gorm:"foreignKey:BiosyntheticClassID"`
	ReleaseTypes  []ReleaseType  `json:"release_types,omitempty" gorm:"foreignKey:BiosyntheticClassID"`

	// PKS
	Cyclases     pq.StringArray `json:"cyclases" gorm:"type:text[]"`
	KetideLength *int           `json:"ketideLength,omitempty"`
	// TODO: starter unit?
	// TODO: iterative?

	// ribosomal
	RIPPType   *string         `json:"ripp_type,omitempty"`
	Details    *string         `json:"details,omitempty"`
	Peptidases pq.StringArray  `json:"peptidases,omitempty" gorm:"type:text[]"`
	Precursors []RippPrecursor `json:"precursors,omitempty" gorm:"foreignKey:BiosyntheticClassID"`

	// saccharide
	GlycosylTransferases []GlycosylTransferase  `json:"glycosyltransferases,omitempty" gorm:"foreignKey:BiosyntheticClassID"`
	Subclusters          []SaccharideSubcluster `json:"subclusters,omitempty" gorm:"foreignKey:BiosyntheticClassID"`

	// terpene
	Prenyltransferases pq.StringArray `json:"prenyltransferases,omitempty" gorm:"type:text[]"`
	SynthasesCyclases  pq.StringArray `json:"synthases_cyclases,omitempty" gorm:"type:text[]"`
	Precursor          *string        `json:"precursor,omitempty"`
}

func init() {
	models.Models = append(models.Models, BiosyntheticClass{})
	// NRPS
	models.Models = append(models.Models, Thioesterase{})
	models.Models = append(models.Models, ReleaseType{})
	// ribosomal
	models.Models = append(models.Models, CleavageLocation{})
	models.Models = append(models.Models, RippPrecursorCrosslink{})
	models.Models = append(models.Models, RippPrecursor{})
	// saccharide
	models.Models = append(models.Models, GlycosylTransferase{})
	models.Models = append(models.Models, SaccharideSubcluster{})
}

func CreateBiosynthesisClass(db *gorm.DB, biosynthId uint64, class BiosyntheticClass) error {
	bioSynth, err := GetBiosynthesisById(db, biosynthId)

	if err != nil {
		return err
	}

	err = db.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Model(&bioSynth).
		Association("Classes").
		Append(&class)

	if err != nil {
		return err
	}

	return nil
}

func UpdateEntryBiosynthesisClass(db *gorm.DB, classId int, newClass BiosyntheticClass) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		var err error

		var oldClass BiosyntheticClass

		err = tx.
			Model(&oldClass).
			Where("id = ?", classId).
			Save(&newClass).
			Error

		if err != nil {
			return err
		}

		switch newClass.Class {
		case "NRPS":
			err = tx.
				Model(&oldClass).
				Association("ReleaseTypes").
				Replace(&newClass.ReleaseTypes)

			if err != nil {
				return err
			}

			err = tx.
				Model(&oldClass).
				Association("Thioesterases").
				Replace(&newClass.Thioesterases)

			if err != nil {
				return err
			}
		case "PKS":
			break
		case "ribosomal":
			err = tx.
				Model(&oldClass).
				Where("id = ?", classId).
				Association("Precursors").
				Replace(&newClass.Precursors)

			if err != nil {
				return err
			}
		case "saccharide":
			err = tx.
				Model(&oldClass).
				Association("GlycosylTransferases").
				Replace(&newClass.GlycosylTransferases)

			if err != nil {
				return err
			}

			err = tx.
				Model(&oldClass).
				Association("Subclusters").
				Replace(&newClass.Subclusters)

			if err != nil {
				return err
			}
		case "terpene":
			break
		case "other":
			break
		}

		return nil
	})

	return err
}

func DeleteEntryBiosynthesisClass(db *gorm.DB, id int) error {
	err := db.
		Model(&BiosyntheticClass{}).
		Delete("id = ?", id).
		Error

	return err
}

func GetEntryBiosynthesisClass(db *gorm.DB, id int) (*BiosyntheticClass, error) {
	var class BiosyntheticClass

	err := db.
		Table("biosynthetic_classes").
		Where("id = ?", id).
		First(&class).
		Error

	if err != nil {
		return nil, err
	}

	// oh god
	switch class.Class {
	case "NRPS":
		db.Model(&class).
			Preload("ReleaseTypes").
			Preload("Thioesterases.Location").
			Find(&class)
	case "PKS":
		break
	case "ribosomal":
		db.Model(&class).
			Preload("Precursors.LeaderCleavageLocation").
			Preload("Precursors.FollowerCleavageLocation").
			Preload("Precursors.Crosslinks").
			Find(&class)
	case "saccharide":
		db.Model(&class).
			Preload("GlycosylTransferases.Evidence").
			Preload("Subclusters").
			Find(&class)
	case "terpene":
		break
	case "other":
		break
	}

	return &class, nil
}
