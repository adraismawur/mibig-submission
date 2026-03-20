package gene

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/lib/pq"
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

type GeneFunctionAnnotationEvidence struct {
	ID                       uint64         `json:"db_id"`
	GeneFunctionAnnotationID uint64         `json:"db_gene_function_annotation_id"`
	Method                   string         `json:"method"`
	References               pq.StringArray `json:"references" gorm:"type:text[]"`
}

type GeneMutationPhenotypeAnnotation struct {
	ID         uint64         `json:"db_id"`
	Phenotype  string         `json:"phenotype"`
	Details    string         `json:"details"`
	References pq.StringArray `json:"references" gorm:"type:text[]"`
}

type GeneFunctionAnnotation struct {
	ID                  uint64                           `json:"db_id"`
	GeneAnnotationID    uint64                           `json:"db_gene_annotation_id"`
	Function            string                           `json:"function"`
	Details             string                           `json:"details"`
	Evidence            []GeneFunctionAnnotationEvidence `json:"evidence" gorm:"foreignKey:GeneFunctionAnnotationID"`
	MutationPhenotypeID uint64                           `json:"db_mutation_phenotype_id"`
	MutationPhenotype   *GeneMutationPhenotypeAnnotation `json:"mutation_phenotype"`
}

type GeneAnnotation struct {
	ID                uint64                   `json:"db_id"`
	GeneInformationID uint64                   `json:"db_gene_information_id"`
	Accession         string                   `json:"accession"` // Accession is the gene ID, e.g. 'AEK75497.1'. This is confusing, but GeneID here is internal to the API
	Name              string                   `json:"name"`      // Name is the actual gene name, e.g. 'abyA1'
	Product           string                   `json:"product"`   // Product is the product of this gene, e.g. '3-oxoacyl-ACP synthase III'
	Functions         []GeneFunctionAnnotation `json:"functions" gorm:"foreignKey:GeneAnnotationID"`
}

type GeneInformation struct {
	ID             uint64           `json:"db_id"`
	EntryAccession string           `json:"db_entry_accession"`
	Additions      []GeneAddition   `json:"to_add,omitempty" gorm:"ForeignKey:GeneInformationID"`
	Deletions      []GeneDeletion   `json:"to_delete,omitempty" gorm:"ForeignKey:GeneInformationID"`
	Annotations    []GeneAnnotation `json:"annotations,omitempty" gorm:"ForeignKey:GeneInformationID"`
}

func init() {
	models.Models = append(models.Models, &GeneInformation{})
	models.Models = append(models.Models, &GeneAddition{})
	models.Models = append(models.Models, &GeneDeletion{})
	models.Models = append(models.Models, &GeneLocation{})
	models.Models = append(models.Models, &ExonLocation{})
	models.Models = append(models.Models, &GeneAnnotation{})
	models.Models = append(models.Models, &GeneFunctionAnnotationEvidence{})
	models.Models = append(models.Models, &GeneMutationPhenotypeAnnotation{})
	models.Models = append(models.Models, &GeneFunctionAnnotation{})
}

func GetEntryGeneInformation(db *gorm.DB, accession string) (*GeneInformation, error) {
	var geneInformation GeneInformation

	err := db.
		Table("gene_informations").
		Joins("JOIN entries ON gene_informations.entry_accession = entries.accession").
		Preload("Additions.Location.Exons").
		Preload("Deletions").
		Preload("Annotations").
		Where("entry_accession = $1", accession).
		First(&geneInformation).
		Error

	if err != nil {
		return nil, err
	}

	return &geneInformation, nil
}

func GetGeneAdditionExists(db *gorm.DB, additionId int) (bool, error) {
	var count int64

	err := db.Table("gene_additions").Where("id = $1", additionId).Count(&count).Error

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

	err := db.Table("gene_deletions").Where("id = $1", deletionId).Count(&count).Error

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

	err := db.Table("gene_annotations").Where("id = $1", annotationId).Count(&count).Error

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
		Where("id = $1", additionId).
		First(&addition).
		Error

	if err != nil {
		return nil, err
	}

	return &addition, err
}

func GetGeneDeletion(db *gorm.DB, deletionId int) (*GeneDeletion, error) {
	var deletion GeneDeletion

	err := db.Table("gene_deletions").Where("id = $1", deletionId).First(&deletion).Error

	if err != nil {
		return nil, err
	}

	return &deletion, err
}

func GetGeneAnnotation(db *gorm.DB, annotationId int) (*GeneAnnotation, error) {
	var annotation GeneAnnotation

	err := db.
		Table("gene_annotations").
		Where("id = $1", annotationId).
		Preload("Functions.Evidence").
		Preload("Functions.MutationPhenotype").
		First(&annotation).
		Error

	if err != nil {
		return nil, err
	}

	return &annotation, err
}

func UpdateOrCreateGeneAddition(db *gorm.DB, entryAccession string, addition *GeneAddition) (*GeneAddition, error) {
	var returnAddition GeneAddition

	err := db.Session(&gorm.Session{FullSaveAssociations: true}).Transaction(func(tx *gorm.DB) error {
		var err error
		var geneInformation GeneInformation
		geneInformation.EntryAccession = entryAccession

		err = tx.
			Model(&GeneInformation{}).
			Where("entry_accession = $1", entryAccession).
			FirstOrCreate(&geneInformation).
			Error

		if err != nil {
			return err
		}

		addition.GeneInformationID = geneInformation.ID

		err = tx.
			Model(&returnAddition).
			Save(&addition).
			Error

		return err
	})

	if err != nil {
		return nil, err
	}

	return &returnAddition, err
}

func UpdateOrCreateGeneDeletion(db *gorm.DB, entryAccession string, deletion *GeneDeletion) (*GeneDeletion, error) {
	var returnDeletion GeneDeletion

	err := db.Session(&gorm.Session{FullSaveAssociations: true}).Transaction(func(tx *gorm.DB) error {
		var err error
		var geneInformation GeneInformation
		geneInformation.EntryAccession = entryAccession

		err = tx.
			Model(&GeneInformation{}).
			Where("entry_accession = $1", entryAccession).
			FirstOrCreate(&geneInformation).
			Error

		if err != nil {
			return err
		}

		deletion.GeneInformationID = geneInformation.ID

		err = tx.
			Model(&returnDeletion).
			Save(&deletion).
			Error

		return err
	})

	if err != nil {
		return nil, err
	}

	return &returnDeletion, err
}

func UpdateOrCreateGeneAnnotation(db *gorm.DB, entryAccession string, annotation *GeneAnnotation) (*GeneAnnotation, error) {
	var returnAnnotation GeneAnnotation

	err := db.Session(&gorm.Session{FullSaveAssociations: true}).Transaction(func(tx *gorm.DB) error {
		var err error
		var geneInformation GeneInformation
		geneInformation.EntryAccession = entryAccession

		err = tx.
			Model(&GeneInformation{}).
			Where("entry_accession = $1", entryAccession).
			FirstOrCreate(&geneInformation).
			Error

		if err != nil {
			return err
		}

		annotation.GeneInformationID = geneInformation.ID

		t := tx.
			Model(&returnAnnotation)

		if annotation.ID != 0 {
			t = t.Where("id = $1", annotation.ID)
		}

		err = t.Save(&annotation).Error

		if err != nil {
			return err
		}

		err = tx.
			Model(&returnAnnotation).
			Association("Functions").
			Replace(&annotation.Functions)

		if err != nil {
			return err
		}

		for i, function := range annotation.Functions {
			err = tx.
				Model(&returnAnnotation.Functions[i]).
				Association("Evidence").
				Replace(&function.Evidence)

			if err != nil {
				return err
			}

			err = tx.
				Model(&GeneMutationPhenotypeAnnotation{}).
				Where("id = $1", function.MutationPhenotypeID).
				Save(&function.MutationPhenotype).
				Error

			if err != nil {
				return err
			}
		}

		return err
	})

	if err != nil {
		return nil, err
	}

	return &returnAnnotation, err
}

func DeleteGeneAddition(db *gorm.DB, additionId int) error {
	var addition GeneAddition

	err := db.
		Table("gene_additions").
		Where("id = $1", additionId).
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
		Where("id = $1", additionId).
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
		Where("id = $1", additionId).
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
