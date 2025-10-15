package entry_utils

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/models/entry"
	"github.com/adraismawur/mibig-submission/util/test_utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateNewAccession(t *testing.T) {
	testDB := test_utils.CreateTestDB()
	models.Migrate(testDB)

	testEntry := entry.Entry{}

	testEntry.Accession = "BGC0011037"

	testDB.Create(&testEntry)

	expectedAccession := "BGC0011038"
	actualAccession, err := GenerateNewAccession(testDB)

	assert.Nil(t, err)
	assert.Equal(t, expectedAccession, actualAccession)
}
