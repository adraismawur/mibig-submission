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
	"strconv"
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

	EnsureEntryDefaults(entry)

	return entry, nil
}

// EnsureEntryDefaults ensures that an entry has sane default values, e.g. no NULL values or fields that should not
// be there
// This is largely done by gorm and the structure of the data already, but some cases cannot be covered using gorm
// or whatever creative structs have been made in this project and simply need to be handled through code
func EnsureEntryDefaults(entry *Entry) {
	for i := range entry.Loci {
		for j := range entry.Loci[i].Evidence {
			if entry.Loci[i].Evidence[j].References != nil {
				continue
			}

			references := make(pq.StringArray, 0)

			entry.Loci[i].Evidence[j].References = references
		}
	}
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
func LoadEntryTransaction(tx *gorm.DB, path string, skip bool) (*Entry, error) {
	jsonString := util.ReadFile(path)
	if jsonString == nil {
		slog.Error("[entry] Could not load entry file ", "path", path)
		return nil, errors.New("could not load entry file")
	}

	var err error
	var entry *Entry

	if entry, err = ParseEntry(jsonString); err != nil {
		if !skip {
			slog.Error("[entry] Failed to parse entry file ", "path", path)
			return nil, err
		}
		slog.Warn("[entry] Failed to parse entry file. Skipping this entry", "path", path)
		return nil, nil
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

			_, err = LoadEntryTransaction(db, fullPath, true)

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

func GetEntryBiosynthesis(db *gorm.DB, accession string) (*biosynthesis.Biosynthesis, error) {
	var entry Entry

	err := db.
		Table("entries").
		Where("accession = ?", accession).
		Preload("Biosynthesis.Classes").
		Preload("Biosynthesis.Modules.Carriers.Location").
		Preload("Biosynthesis.Modules.ModificationDomains.Location").
		Preload("Biosynthesis.Modules.ATDomain.Location").
		Preload("Biosynthesis.Modules.KSDomain.Location").
		First(&entry).
		Error

	if err != nil {
		return nil, err
	}

	return &entry.Biosynthesis, nil
}

func CreateEntryBiosynthesisModule(db *gorm.DB, accession string, module biosynthesis.BiosyntheticModule) error {

	var entry Entry

	err := db.
		Table("entries").
		Where("accession = ?", accession).
		Preload("Biosynthesis.Classes").
		Preload("Biosynthesis.Modules.Carriers.Location").
		Preload("Biosynthesis.Modules.ModificationDomains.Location").
		Preload("Biosynthesis.Modules.ADomain.Location").
		Preload("Biosynthesis.Modules.ATDomain.Location").
		Preload("Biosynthesis.Modules.KSDomain.Location").
		First(&entry).
		Error

	if err != nil {
		return err
	}

	// get new module number if it doesn't exist
	if module.Name == "" {
		module.Name = strconv.Itoa(len(entry.Biosynthesis.Modules) + 1)
	}

	err = db.Model(&entry.Biosynthesis).Association("Modules").Append(&module)

	if err != nil {
		return err
	}

	return nil
}

func UpdateEntryBiosynthesisModule(db *gorm.DB, accession string, newModule *biosynthesis.BiosyntheticModule) error {

	var entry Entry

	err := db.
		Table("entries").
		Where("accession = ?", accession).
		Preload("Biosynthesis.Classes").
		Preload("Biosynthesis.Modules.Carriers.Location").
		Preload("Biosynthesis.Modules.ModificationDomains.Location").
		Preload("Biosynthesis.Modules.ADomain.Location").
		Preload("Biosynthesis.Modules.ATDomain.Location").
		Preload("Biosynthesis.Modules.KSDomain.Location").
		First(&entry).
		Error

	if err != nil {
		return err
	}

	// find correct newModule
	found := false
	for _, module := range entry.Biosynthesis.Modules {
		if module.Name == module.Name {
			newModule.ID = module.ID
			newModule.BiosynthesisID = module.BiosynthesisID
			found = true
			break
		}
	}

	if !found {
		return errors.New("could not find newModule")
	}

	err = db.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Model(&entry.Biosynthesis).
		Association("Modules").
		Replace(newModule)

	if err != nil {
		return err
	}

	return nil
}

func DeleteEntryBiosynthesisModule(db *gorm.DB, accession string, name string) error {
	var entry Entry

	err := db.
		Table("entries").
		Where("accession = ?", accession).
		Preload("Biosynthesis.Modules").
		First(&entry).
		Error

	if err != nil {
		return err
	}

	err = db.Where("name = ?", name).Delete(entry.Biosynthesis.Modules).Error

	if err != nil {
		return err
	}

	return nil
}

func GetEntryBiosynthesisModule(db *gorm.DB, accession string, name string) (*biosynthesis.BiosyntheticModule, error) {
	var entry Entry

	err := db.
		Table("entries").
		Where("accession = ?", accession).
		Preload("Biosynthesis.Classes").
		Preload("Biosynthesis.Modules.Carriers.Location").
		Preload("Biosynthesis.Modules.ModificationDomains.Location").
		Preload("Biosynthesis.Modules.ADomain.Location").
		Preload("Biosynthesis.Modules.ATDomain.Location").
		Preload("Biosynthesis.Modules.KSDomain.Location").
		First(&entry).
		Error

	if err != nil {
		return nil, err
	}

	if len(entry.Biosynthesis.Modules) == 0 {
		return nil, nil
	}

	// TODO: use db to search. I can't figure out what gorm wants from me here. I am really regretting choosing gorm
	// at this point. Just use raw sql and save yourself a lot of trouble. Maybe if their documentation was better
	// I wouldn't spend hours of my life doing what I feel are very simple things. aaaaaaargh
	for _, module := range entry.Biosynthesis.Modules {
		if module.Name == name {
			return &module, nil
		}
	}

	// not found
	return nil, nil
}
