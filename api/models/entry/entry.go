package entry

import (
	"errors"
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/models/entry/biosynthesis"
	"github.com/adraismawur/mibig-submission/models/entry/compound"
	"github.com/adraismawur/mibig-submission/models/entry/gene"
	"github.com/adraismawur/mibig-submission/models/entry/taxonomy"
	"github.com/adraismawur/mibig-submission/util"
	"github.com/goccy/go-json"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log/slog"
	"os"
	"path/filepath"
)

type Quality string

const (
	Questionable Quality = "questionable"
	Medium       Quality = "medium"
	High         Quality = "high"
)

type Status string

const (
	Pending Status = "pending"
	Active  Status = "active"
	Retired Status = "retired"
)

type Completeness string

const (
	Unknown  Completeness = "unknown"
	Complete Completeness = "complete"
)

type Entry struct {
	ID               uint                      `json:"-"`
	Accession        string                    `json:"accession"`
	Version          int                       `json:"version,omitempty"`
	Changelog        Changelog                 `json:"changelog" gorm:"foreignKey:EntryID"`
	Quality          Quality                   `json:"quality,omitempty"`
	Status           Status                    `json:"status,omitempty"`
	Completeness     Completeness              `json:"completeness"`
	Loci             []Locus                   `json:"loci" gorm:"foreignKey:EntryID"`
	Biosynthesis     biosynthesis.Biosynthesis `json:"biosynthesis" gorm:"foreignKey:EntryID"`
	Compounds        []compound.Compound       `json:"compounds" gorm:"ForeignKey:EntryID"`
	Taxonomy         taxonomy.Taxonomy         `json:"taxonomy" gorm:"ForeignKey:EntryID"`
	Genes            *gene.Gene                `json:"genes,omitempty" gorm:"ForeignKey:EntryID"`
	LegacyReferences pq.StringArray            `json:"legacy_references,omitempty" gorm:"type:text[]"`
	Embargo          bool                      `json:"-,omitempty"`
}

func init() {
	models.Models = append(models.Models, &Entry{})
}

// ParseEntry attempts to parse an entry json given as a byte array into an entry struct
func ParseEntry(jsonString []byte) (*Entry, error) {
	var entry *Entry

	if err := json.Unmarshal(jsonString, &entry); err != nil {
		slog.Error("[entry] Failed to unmarshal annotation JSON", "error", err.Error())
		return nil, err
	}

	return entry, nil
}

// LoadEntry attempts to read a file at a given path and load it as an entry into the database
func LoadEntry(db *gorm.DB, path string) (*Entry, error) {
	jsonString := util.ReadFile(path)
	if jsonString == nil {
		slog.Error("[entry] Could not load entry file ", "path", path)
		return nil, errors.New("could not load entry file")
	}

	var err error
	var entry *Entry

	if entry, err = ParseEntry(jsonString); err != nil {
		slog.Error("[entry] Failed to parse entry file ", "path", path)
		return nil, err
	}

	if err = db.Create(entry).Error; err != nil {
		slog.Error("[db] Failed to create entry ", "error", err.Error())
		return nil, err
	}

	return entry, nil
}

// LoadEntryTransaction attempts to read a file at a given path and load it as an entry into the database. This function
// uses a transaction to speed up batch writing
func LoadEntryTransaction(tx *gorm.DB, path string) (*Entry, error) {
	jsonString := util.ReadFile(path)
	if jsonString == nil {
		slog.Error("[entry] Could not load entry file ", "path", path)
		return nil, errors.New("could not load entry file")
	}

	var err error
	var entry *Entry

	if entry, err = ParseEntry(jsonString); err != nil {
		slog.Error("[entry] Failed to parse entry file ", "path", path)
		return nil, err
	}

	if err := tx.Create(entry).Error; err != nil {
		return nil, err
	}

	return entry, nil
}

// LoadEntries attempts to read all files at a given path and load them as entries into the database
func LoadEntries(db *gorm.DB, path string) {
	files, err := os.ReadDir(path)

	if err != nil {
		slog.Error("[db] Failed to read directory", "path", path)
		return
	}

	result, err := db.Table("entries").Select("accession").Rows()

	var accessions = map[string]bool{}

	// load a list of accessions that already exist
	var accession string
	for result.Next() {
		result.Scan(&accession)
		accessions[accession] = true
	}

	var _ *Entry

	db.Transaction(func(tx *gorm.DB) error {

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			baseName := file.Name()[:len(file.Name())-5]

			_, exists := accessions[baseName]

			if exists {
				continue
			}

			fullPath := filepath.Join(path, file.Name())

			_, err = LoadEntryTransaction(db, fullPath)

			if err != nil {
				slog.Error("[db] Failed to load entry", "path", fullPath)
				return err
			}
		}

		return nil
	})
}

func GetEntryExists(db *gorm.DB, accession string) (bool, error) {
	var count int64

	err := db.Table("entries").Where("accession = ?", accession).Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func GetEntryFromAccession(db *gorm.DB, accession string) (*Entry, error) {
	var entry Entry

	// TODO: there must be a better way to do this. no amount of googling on my part gets me anywhere though
	// ideally doing this amount of preloading is rare. This is done here on getting the entire entry
	err := db.
		Table("entries").
		Where("accession = ?", accession).
		Preload("Changelog.Releases.Entries").
		Preload("Loci.Location").
		Preload("Loci.Evidence").
		Preload("Biosynthesis.Classes").
		Preload("Biosynthesis.Modules.Carriers.Location").
		Preload("Biosynthesis.Modules.ModificationDomains.Location").
		Preload("Biosynthesis.Modules.ATDomain.Location").
		Preload("Biosynthesis.Modules.KSDomain.Location").
		Preload("Genes.Additions.Location.Exons").
		Preload("Genes.Annotations").
		Preload("Compounds.Evidence").
		Preload("Compounds.BioActivities").
		Preload(clause.Associations).
		First(&entry).
		Error

	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func GetUserEntries(db *gorm.DB, userId int) ([]string, error) {
	var accessions []string

	err := db.
		Table("user_entries").
		Select("accession").
		Joins("JOIN entries ON entries.id = user_entries.entry_id").
		Where("user_id = ?", userId).
		Find(&accessions).
		Error

	if err != nil {
		return nil, err
	}

	return accessions, nil
}
