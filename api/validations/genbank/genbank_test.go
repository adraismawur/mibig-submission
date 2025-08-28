package genbank

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const knownExistingGb = "JF712342.1"
const knownNonExistingGb = "AA00000.1"

func TestGetGenbankAccession(t *testing.T) {

	actualExists, err := GetGenbankAccession(knownExistingGb)

	if err != nil {
		panic(err)
	}

	assert.NotNil(t, actualExists, "GetGenbankAccession with GB "+knownExistingGb+" should not return nil")
}

func TestGenbankAccessionDoesNotExist(t *testing.T) {
	actualExists, err := GetGenbankAccession(knownNonExistingGb)

	if err != nil {
		panic(err)
	}

	assert.Nil(t, actualExists, "GetGenbankAccession with GB "+knownNonExistingGb+" should return nil")
}

func TestGetGenbankAccessionGetLength(t *testing.T) {

	result, err := GetGenbankAccession(knownExistingGb)

	if err != nil {
		panic(err)
	}

	expectedLen := 696
	actualLen := result.Slen

	assert.Equal(t, expectedLen, actualLen)
}
