package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const knownExistingGb = "JF712342.1"
const knownNonExistingGb = "AA00000.1"

func TestGetGenbankAccession(t *testing.T) {

	actualExists, err := GetGenbankAccessionSummary(knownExistingGb)

	if err != nil {
		panic(err)
	}

	assert.NotNil(t, actualExists, "GetGenbankAccessionSummary with GB "+knownExistingGb+" should not return nil")
}

func TestGenbankAccessionDoesNotExist(t *testing.T) {
	actualExists, err := GetGenbankAccessionSummary(knownNonExistingGb)

	if err != nil {
		panic(err)
	}

	assert.Nil(t, actualExists, "GetGenbankAccessionSummary with GB "+knownNonExistingGb+" should return nil")
}

func TestGetGenbankAccessionGetLength(t *testing.T) {

	result, err := GetGenbankAccessionSummary(knownExistingGb)

	if err != nil {
		panic(err)
	}

	expectedLen := 696
	actualLen := result.SLen

	assert.Equal(t, expectedLen, actualLen)
}
