package submission

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/util"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

var testDb *gorm.DB

var testSubmission = Submission{
	Accession:    "BGC0000433",
	Version:      5,
	Quality:      Questionable,
	Status:       Active,
	Completeness: Complete,
}

func TestMain(m *testing.M) {
	testDb = util.CreateTestDB()
	models.Migrate(testDb)
}

func TestParseSubmission(t *testing.T) {
	var submission Submission

	jsonString := util.ReadFile("testdata/test_submission.json")

	err := json.Unmarshal([]byte(jsonString), &submission)

	assert.NoError(t, err)

	assert.Equal(t, submission.Accession, "BGC0000001")
	assert.Equal(t, submission.Version, 4)
	assert.Equal(t, submission.Quality, Medium)
	assert.Equal(t, submission.Status, Active)
	assert.Equal(t, submission.Completeness, Unknown)
}

func TestLoadSubmission(t *testing.T) {
	actualSubmission := testSubmission

	submission, err := LoadSubmission(testDb, "submission_testdata/BGC0000433.json")

	assert.NoError(t, err)
	assert.Equal(t, submission, actualSubmission)
}

func TestGetSubmissionExists(t *testing.T) {
	testDb.Create(&testSubmission)

	actualResult, err := GetSubmissionExists(testDb, testSubmission.Accession)

	assert.NoError(t, err)
	assert.Equal(t, true, actualResult)
}
