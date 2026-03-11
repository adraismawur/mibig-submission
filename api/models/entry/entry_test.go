package entry

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/models/entry/consts"
	"github.com/adraismawur/mibig-submission/util"
	"github.com/adraismawur/mibig-submission/util/test_utils"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

var testDb *gorm.DB

var testEntry = Entry{
	Accession:    "BGC0000433",
	Version:      5,
	Quality:      consts.Questionable,
	Status:       consts.Active,
	Completeness: consts.Complete,
}

func TestMain(m *testing.M) {
	testDb = test_utils.CreateTestDB()
	models.Migrate(testDb)
}

func TestParseEntry(t *testing.T) {
	var entry Entry

	jsonString := util.ReadFile("testdata/test_entry.json")

	err := json.Unmarshal([]byte(jsonString), &entry)

	assert.NoError(t, err)

	assert.Equal(t, entry.Accession, "BGC0000001")
	assert.Equal(t, entry.Version, 4)
	assert.Equal(t, entry.Quality, consts.Medium)
	assert.Equal(t, entry.Status, consts.Active)
	assert.Equal(t, entry.Completeness, consts.Unknown)
}

func TestLoadEntry(t *testing.T) {
	expectedEntry := testEntry

	actualEntry, err := LoadEntry(testDb, "entry_testdata/BGC0000433.json")

	assert.NoError(t, err)
	assert.Equal(t, expectedEntry, actualEntry)
}

func TestGetEntryExists(t *testing.T) {
	testDb.Create(&testEntry)

	actualResult, err := GetEntryExists(testDb, testEntry.Accession)

	assert.NoError(t, err)
	assert.Equal(t, true, actualResult)
}
