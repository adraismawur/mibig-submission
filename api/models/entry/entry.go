package entry

import (
	"errors"
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/models/entry/biosynthesis"
	"github.com/adraismawur/mibig-submission/models/entry/compound"
	"github.com/adraismawur/mibig-submission/models/entry/consts"
	"github.com/adraismawur/mibig-submission/models/entry/gene"
	"github.com/adraismawur/mibig-submission/models/entry/locus"
	"github.com/adraismawur/mibig-submission/models/entry/taxonomy"
	"github.com/adraismawur/mibig-submission/util"
	"github.com/goccy/go-json"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
)

type MinimalEntry struct {
	Locus     locus.Locus
	Compounds []compound.Compound
}

type Gene struct {
	ID             uint64
	EntryAccession string
	Name           string
}

type Entry struct {
	Accession        string                    `json:"accession" gorm:"primaryKey"`
	Version          int                       `json:"version,omitempty"`
	Changelog        Changelog                 `json:"changelog" gorm:"foreignKey:EntryAccession"`
	Quality          consts.Quality            `json:"quality,omitempty"`
	Status           consts.Status             `json:"status,omitempty"`
	Completeness     consts.Completeness       `json:"completeness"`
	Loci             []locus.Locus             `json:"loci" gorm:"foreignKey:EntryAccession"`
	Biosynthesis     biosynthesis.Biosynthesis `json:"biosynthesis" gorm:"foreignKey:EntryAccession"`
	Compounds        []compound.Compound       `json:"compounds" gorm:"ForeignKey:EntryAccession"`
	Taxonomy         taxonomy.Taxonomy         `json:"taxonomy" gorm:"ForeignKey:EntryAccession"`
	GeneInformation  *gene.GeneInformation     `json:"genes,omitempty" gorm:"ForeignKey:EntryAccession"`
	LegacyReferences pq.StringArray            `json:"legacy_references,omitempty" gorm:"type:text[]"`
	Comment          string                    `json:"comment"`

	// internal data starts here

	GeneList []Gene `json:"-" gorm:"ForeignKey:EntryAccession"`
	Embargo  bool   `json:"embargo,omitempty"`
}

func init() {
	models.Models = append(models.Models, &Entry{})
	models.Models = append(models.Models, &Gene{})
}

func ParseEntryFallback(entryJson []byte, entry *Entry) error {
	// first read the json into a map
	entryMap := make(map[string]any)

	err := json.Unmarshal(entryJson, &entryMap)

	if err != nil {
		return err
	}

	// first problem encountered: sometimes the bioactivity name does not correspond to the schema:
	compounds := entryMap["compounds"].([]interface{})

	for _, compound := range compounds {
		c := compound.(map[string]interface{})

		bioactivities, ok := c["bioactivities"].([]interface{})
		if !ok {
			continue
		}
		for _, bioactivity := range bioactivities {
			b := bioactivity.(map[string]interface{})
			nameType := reflect.TypeOf(b["name"])
			if nameType.Kind() != reflect.String {
				slog.Warn("[entry] Found fix: bad bioactivity name", "entry", entry.Accession)
				actualName := b["name"].(map[string]interface{})["activity"]
				b["name"] = actualName
			}
		}
	}

	// second problem: some ripp precursor core sequences are arrays of strings
	// most cases have an array with one tring, but some are arrays of multiple strings.
	// convert these latter ones into multiple precursors
	if entryMap["biosynthesis"] == nil {
		return nil
	}

	biosynth := entryMap["biosynthesis"].(map[string]interface{})

	classes := biosynth["classes"].([]interface{})

	for _, class := range classes {
		c := class.(map[string]interface{})

		var additionalPrecursors []interface{}

		precursors, ok := c["precursors"].([]interface{})
		if !ok {
			continue
		}

		if len(precursors) == 0 {
			continue
		}

		for i, _ := range precursors {
			p := precursors[i].(map[string]interface{})
			structureType := reflect.TypeOf(p["core_sequence"])
			if structureType.Kind() != reflect.String {
				pArray := p["core_sequence"].([]interface{})
				if len(pArray) <= 1 {
					slog.Warn("[entry] Found fix: bad ripp precursor core_sequence, single element", "entry", entry.Accession)
					p["core_sequence"] = pArray[0].(interface{})
					slog.Warn("[entry] Corrected to " + p["core_sequence"].(string))
				} else {
					slog.Warn("Found fix: bad ripp precursor core_sequence, multiple elements", "entry", entry.Accession)
					for _, coreSequence := range pArray[1:] {
						newPrecursor := map[string]interface{}{}

						if p["gene"] != nil {
							newPrecursor["gene"] = p["gene"]
						}
						if p["leader_cleavage_location"] != nil {
							newPrecursor["leader_cleavage_location"] = p["leader_cleavage_location"]
						}
						if p["follower_cleavage_location"] != nil {
							newPrecursor["follower_cleavage_location"] = p["follower_cleavage_location"]
						}
						if p["crosslinks"] != nil {
							newPrecursor["crosslinks"] = p["crosslinks"]
						}
						if p["recognition_motif"] != nil {
							newPrecursor["recognition_motif"] = p["recognition_motif"]
						}

						newPrecursor["core_sequence"] = coreSequence

						additionalPrecursors = append(additionalPrecursors, newPrecursor)
					}

					asd := pArray[0]

					p["core_sequence"] = asd
				}
			}
		}
		for _, precursor := range additionalPrecursors {
			c["precursors"] = append(precursors, precursor)
		}
	}

	j, err := json.Marshal(entryMap)
	//var tempEntryMap map[string]any
	err = json.Unmarshal(j, &entry)

	//err = mapstructure.Decode(&tempEntryMap, &entry)

	if err != nil {
		return err
	}

	return nil
}

// ParseEntry attempts to parse an entry json given as a byte array into an entry struct
func ParseEntry(jsonString []byte) (*Entry, error) {
	entry := Entry{}

	if err := json.Unmarshal(jsonString, &entry); err != nil {
		slog.Warn("[entry] Failed to unmarshal annotation JSON. Attempting to recover...", "error", err.Error())

		err = ParseEntryFallback(jsonString, &entry)

		if err != nil {
			slog.Error("[entry] Failed to unmarshal annotation JSON using fallback.", "error", err.Error())
			return nil, err
		}
	}

	EnsureEntryDefaults(&entry)

	PopulateBiosynthIndexes(&entry)

	return &entry, nil
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

// PopulateBiosynthIndexes fills in the Index fields in biosynth module lists. These fields are later used to maintain
// order in the list of modules
func PopulateBiosynthIndexes(entry *Entry) {
	for i := range entry.Biosynthesis.Modules {
		entry.Biosynthesis.Modules[i].Index = uint64(i + 1)
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
		} else {
			slog.Warn("[entry] Failed to parse entry file. Skipping this entry", "path", path)
		}
		return nil, nil
	}

	if err := tx.Create(entry).Error; err != nil {
		return nil, err
	}

	slog.Info("[entry] Successfully preloaded entry", "path", path)

	return entry, nil
}

// LoadEntries attempts to read all files at a given path and load them as entries into the database
func LoadEntries(db *gorm.DB, path string) error {
	files, err := os.ReadDir(path)

	if err != nil {
		slog.Error("[db] Failed to read directory", "path", path)
		return err
	}

	result, err := db.Table("entries").Select("accession").Rows()

	if err != nil {
		slog.Error("[db] Failed to read entries table", "path", path)
		return err
	}

	var accessions = map[string]bool{}

	// load a list of accessions that already exist
	var accession string

	for result.Next() {

		err = result.Scan(&accession)

		if err != nil {
			slog.Error("[db] Failed to scan entry table row", "path", path)
			return err
		}

		accessions[accession] = true
	}

	var _ *Entry

	err = db.Transaction(func(tx *gorm.DB) error {

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

			_, err = LoadEntryTransaction(tx, fullPath, true)

			if err != nil {
				slog.Error("[db] Failed to load entry", "path", fullPath)
				return err
			}
		}

		return nil
	})

	if err != nil {
		slog.Error("[db] Failed to load entries table", "path", path)
		return err
	}

	return nil
}

func GetEntryExists(db *gorm.DB, accession string) (bool, error) {
	var count int64

	err := db.Table("entries").Where("accession = $1", accession).Count(&count).Error

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
		Where("accession = $1", accession).
		Preload("Changelog.Releases.Entries").
		Preload("Loci.Location").
		Preload("Loci.Evidence").
		Preload("Biosynthesis.Classes.ReleaseTypes").
		Preload("Biosynthesis.Classes.Thioesterases.Location").
		Preload("Biosynthesis.Classes.Precursors.LeaderCleavageLocation").
		Preload("Biosynthesis.Classes.Precursors.FollowerCleavageLocation").
		Preload("Biosynthesis.Classes.Precursors.Crosslinks").
		Preload("Biosynthesis.Classes.GlycosylTransferases.Evidence").
		Preload("Biosynthesis.Classes.Subclusters").
		Preload("Biosynthesis.Modules.IntegratedMonomers").
		Preload("Biosynthesis.Modules.IntegratedMonomers.Evidence").
		Preload("Biosynthesis.Modules.Carriers.Location").
		Preload("Biosynthesis.Modules.Carriers.Evidence").
		Preload("Biosynthesis.Modules.ModificationDomains.Location").
		Preload("Biosynthesis.Modules.ModificationDomains.Substrates").
		Preload("Biosynthesis.Modules.ModificationDomains.Evidence").
		Preload("Biosynthesis.Modules.CDomain.Location").
		Preload("Biosynthesis.Modules.ADomain.Location").
		Preload("Biosynthesis.Modules.ADomain.Evidence").
		Preload("Biosynthesis.Modules.ADomain.Substrates").
		Preload("Biosynthesis.Modules.ATDomain.Location").
		Preload("Biosynthesis.Modules.ATDomain.Substrates").
		Preload("Biosynthesis.Modules.ATDomain.Evidence").
		Preload("Biosynthesis.Modules.KSDomain.Location").
		Preload("Biosynthesis.Operons").
		Preload("Biosynthesis.Paths.Products").
		Preload("GeneInformation.Additions.Location.Exons").
		Preload("GeneInformation.Annotations").
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

func GetEntryGenes(db *gorm.DB, accession string) (*[]string, error) {
	genes := make([]string, 0)

	db.Model(&Gene{}).
		Select("name").
		Where("genes.entry_accession = $1", accession).
		Find(&genes)

	return &genes, nil
}
