package biosynthesis

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"log/slog"
)

type BiosyntheticPathwayProduct struct {
	ID                    uint64 `json:"db_id"`
	BiosyntheticPathwayID uint64 `json:"db_path_id"`
	Name                  string `json:"name"`
	Structure             string `json:"structure"`
	Comment               string `json:"comment"`
}

type BiosyntheticPathway struct {
	ID                uint64                       `json:"db_id"`
	BiosynthesisID    uint64                       `json:"db_biosynth_id"`
	Products          []BiosyntheticPathwayProduct `json:"products" gorm:"foreignKey:BiosyntheticPathwayID"`
	Steps             string                       `json:"steps"`
	References        pq.StringArray               `json:"references" gorm:"type:text[]"`
	IsSubCluster      bool                         `json:"is_sub_cluster"`
	ProducesPrecursor bool                         `json:"produces_precursor"`
}

func init() {
	models.Models = append(models.Models, &BiosyntheticPathway{})
	models.Models = append(models.Models, &BiosyntheticPathwayProduct{})
}

func GetBiosynthesisPath(db *gorm.DB, pathId int) (*BiosyntheticPathway, error) {
	var path BiosyntheticPathway

	err := db.
		Model(&BiosyntheticPathway{}).
		Where("id = $1", pathId).
		Preload("Products").
		First(&path).
		Error

	if err != nil {
		return nil, err
	}

	return &path, nil
}

func CreateBiosynthesisPath(db *gorm.DB, path BiosyntheticPathway) error {
	biosynth, err := GetBiosynthesisById(db, path.BiosynthesisID)

	if err != nil {
		return err
	}

	err = db.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Model(&biosynth).
		Association("Paths").
		Append(&path)

	if err != nil {
		return err
	}

	return nil
}

func UpdateBiosynthesisPath(db *gorm.DB, path BiosyntheticPathway) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		var oldPath BiosyntheticPathway
		err := tx.
			Model(&oldPath).
			Where("id = $1", path.ID).
			Save(&path).
			Error

		if err != nil {
			slog.Error("[biosynthetic_path] Failed to save pathway", "error", err)
			return err
		}

		err = tx.
			Model(&oldPath).
			Association("Products").
			Replace(&path.Products)

		if err != nil {
			slog.Error("[biosynthetic_path] Failed to replace products", "error", err)
			return err
		}

		return nil
	})

	return err
}

func DeleteBiosynthesisPath(db *gorm.DB, pathId int) error {
	err := db.
		Model(&BiosyntheticPathway{}).
		Delete("id = $1", pathId).
		Error

	if err != nil {
		return err
	}

	return nil
}
