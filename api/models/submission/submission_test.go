package submission

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/util"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"io"
	"os"
	"testing"
)

var testDb *gorm.DB

func TestMain(m *testing.M) {
	testDb = util.CreateTestDB()
	models.Migrate(testDb)
}

func LoadJson(t *testing.T, path string) []byte {
	var err error
	var file *os.File

	if file, err = os.Open(
		"submission_testdata/annotations/BGC0000433.json",
	); err != nil {
		t.Fatal(err)
	}

	var jsonString []byte
	if jsonString, err = io.ReadAll(file); err != nil {
		t.Fatal(err)
	}

	return jsonString
}

func TestLoadJson(t *testing.T) {

	var submission Submission

	jsonString := LoadJson(t, "submission_testdata/annotations/BGC0000433.json")
	err := json.Unmarshal(jsonString, &submission)

	assert.NoError(t, err)

	assert.Equal(t, submission.Accession, "BGC0000433")
	assert.Equal(t, submission.Version, 5)
	assert.Equal(t, submission.Quality, Questionable)
	assert.Equal(t, submission.Status, Active)
	assert.Equal(t, submission.Completeness, Complete)
}

func TestSaveJson(t *testing.T) {
	var err error
	var submission Submission

	jsonString := LoadJson(t, "submission_testdata/annotations/BGC0000433.json")
	if err = json.Unmarshal(jsonString, &submission); err != nil {
		t.Fatal(err)
	}

	err = testDb.Model(&Submission{}).Create(&Submission{}).Error

	assert.NoError(t, err)

	actualSubmission := Submission{}
	err = testDb.Model(&Submission{}).First(&actualSubmission).Error
	assert.NoError(t, err)
	assert.Equal(t, submission, actualSubmission)
}
