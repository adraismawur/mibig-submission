package gene

import (
	"github.com/adraismawur/mibig-submission/models"
	"gorm.io/gorm"
	"log/slog"
)

// TODO: there may be a way to unify this, locus.location and BiosyntheticModuleDomainLocation
type ExonLocation struct {
	ID             uint64 `json:"db_id"`
	GeneLocationID uint64 `json:"db_location_id"`
	From           uint64 `json:"from"`
	To             uint64 `json:"to"`
}

type GeneLocation struct {
	ID     uint64         `json:"db_id"`
	GeneID uint64         `json:"db_gene_id"`
	Exons  []ExonLocation `json:"exons" gorm:"ForeignKey:GeneLocationID"`
	Strand int32          `json:"strand"`
}

type GeneAddition struct {
	ID                uint64       `json:"db_id"`
	GeneInformationID uint64       `json:"gene_information_id"`
	Accession         string       `json:"accession"`
	Location          GeneLocation `json:"location" gorm:"ForeignKey:GeneID"`
	Translation       string       `json:"translation"`
}

type GeneDeletion struct {
	ID                uint64 `json:"db_id"`
	GeneInformationID uint64 `json:"gene_information_id"`
	Accession         string `json:"accession"`
	Reason            string `json:"reason"`
}

type GeneAnnotation struct {
	ID                uint64 `json:"db_id"`
	GeneInformationID uint64 `json:"gene_information_id"`
	Accession         string `json:"accession"` // Accession is the gene ID, e.g. 'AEK75497.1'. This is confusing, but GeneID here is internal to the API
	Name              string `json:"name"`      // Name is the actual gene name, e.g. 'abyA1'
	Product           string `json:"product"`   // Product is the product of this gene, e.g. '3-oxoacyl-ACP synthase III'
}

type GeneInformation struct {
	ID          uint64           `json:"db_id"`
	EntryID     uint64           `json:"entry_id"`
	Additions   []GeneAddition   `json:"to_add,omitempty" gorm:"ForeignKey:GeneInformationID"`
	Deletions   []GeneDeletion   `json:"to_delete,omitempty" gorm:"ForeignKey:GeneInformationID"`
	Annotations []GeneAnnotation `json:"annotations,omitempty" gorm:"ForeignKey:GeneInformationID"`
}

func init() {
	models.Models = append(models.Models, &GeneInformation{})
	models.Models = append(models.Models, &GeneAddition{})
	models.Models = append(models.Models, &GeneDeletion{})
	models.Models = append(models.Models, &GeneLocation{})
	models.Models = append(models.Models, &ExonLocation{})
	models.Models = append(models.Models, &GeneAnnotation{})
}

func GetEntryGeneInformation(db *gorm.DB, accession string) (*GeneInformation, error) {
	var geneInformation GeneInformation

	err := db.
		Table("gene_informations").
		Joins("JOIN entries ON gene_informations.entry_id = entries.id").
		Preload("Additions.Location.Exons").
		Preload("Deletions").
		Preload("Annotations").
		Where("accession = ?", accession).
		First(&geneInformation).
		Error

	if err != nil {
		return nil, err
	}

	return &geneInformation, nil
}

func GetGeneAdditionExists(db *gorm.DB, additionId int) (bool, error) {
	var count int64

	err := db.Table("gene_additions").Where("id = ?", additionId).Count(&count).Error

	if err != nil {
		slog.Error("[genes] Error finding gene addition", "addition_id", additionId, "error", err)
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

func GetGeneDeletionExists(db *gorm.DB, deletionId int) (bool, error) {
	var count int64

	err := db.Table("gene_deletions").Where("id = ?", deletionId).Count(&count).Error

	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

func GetGeneAnnotationExists(db *gorm.DB, annotationId int) (bool, error) {
	var count int64

	err := db.Table("gene_annotations").Where("id = ?", annotationId).Count(&count).Error

	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

func GetGeneAddition(db *gorm.DB, additionId int) (*GeneAddition, error) {
	var addition GeneAddition

	err := db.Table("gene_additions").
		Preload("Location.Exons").
		Where("id = ?", additionId).
		First(&addition).
		Error

	if err != nil {
		return nil, err
	}

	return &addition, err
}

func GetGeneDeletion(db *gorm.DB, deletionId int) (*GeneDeletion, error) {
	var deletion GeneDeletion

	err := db.Table("gene_deletions").Where("id = ?", deletionId).First(&deletion).Error

	if err != nil {
		return nil, err
	}

	return &deletion, err
}

func GetGeneAnnotation(db *gorm.DB, annotationId int) (*GeneAnnotation, error) {
	var annotation GeneAnnotation

	err := db.Table("gene_annotations").Where("id = ?", annotationId).First(&annotation).Error

	if err != nil {
		return nil, err
	}

	return &annotation, err
}

func UpdateOrCreateGeneAddition(db *gorm.DB, addition *GeneAddition) (*GeneAddition, error) {

	tx := db.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Begin()

	err := tx.
		Table("gene_additions").
		Save(&addition).
		Error

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit().Error

	if err != nil {
		return nil, err
	}

	return addition, nil
}

func UpdateOrCreateGeneDeletion(db *gorm.DB, deletion *GeneDeletion) (*GeneDeletion, error) {

	tx := db.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Begin()

	err := tx.
		Table("gene_deletions").
		Save(&deletion).
		Error

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit().Error

	if err != nil {
		return nil, err
	}

	return deletion, nil
}

func UpdateOrCreateGeneAnnotation(db *gorm.DB, annotation *GeneAnnotation) (*GeneAnnotation, error) {

	tx := db.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Begin()

	err := tx.
		Table("gene_annotations").
		Save(&annotation).
		Error

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit().Error

	if err != nil {
		return nil, err
	}

	return annotation, nil
}

func DeleteGeneAddition(db *gorm.DB, additionId int) error {
	var addition GeneAddition

	err := db.
		Table("gene_additions").
		Where("id = ?", additionId).
		First(&addition).
		Error

	if err != nil {
		return err
	}

	err = db.Delete(&addition).Error

	if err != nil {
		return err
	}

	return nil
}

func DeleteGeneDeletion(db *gorm.DB, additionId int) error {
	var deletion GeneDeletion

	err := db.
		Table("gene_deletions").
		Where("id = ?", additionId).
		First(&deletion).
		Error

	if err != nil {
		return err
	}

	err = db.Delete(&deletion).Error

	if err != nil {
		return err
	}

	return nil
}

func DeleteGeneAnnotation(db *gorm.DB, additionId int) error {
	var annotation GeneAnnotation

	err := db.
		Table("gene_annotations").
		Where("id = ?", additionId).
		First(&annotation).
		Error

	if err != nil {
		return err
	}

	err = db.Delete(&annotation).Error

	if err != nil {
		return err
	}

	return nil
}
