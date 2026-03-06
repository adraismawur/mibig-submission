package biosynthesis

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/util/test_utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func TestReorderEntryBiosynthesisModules(t *testing.T) {
	testBiosynth := Biosynthesis{
		ID:      0,
		EntryID: 0,
		Classes: nil,
		Modules: []BiosyntheticModule{
			{
				ID:                  1,
				Index:               1,
				BiosynthesisID:      0,
				Type:                "pks-modular",
				Name:                "test1",
				Genes:               nil,
				Active:              false,
				Carriers:            nil,
				ModificationDomains: nil,
				ADomain:             nil,
				ATDomain:            nil,
				KSDomain:            nil,
			},
			{
				ID:                  2,
				Index:               2,
				BiosynthesisID:      0,
				Type:                "pks-modular",
				Name:                "test2",
				Genes:               nil,
				Active:              false,
				Carriers:            nil,
				ModificationDomains: nil,
				ADomain:             nil,
				ATDomain:            nil,
				KSDomain:            nil,
			},
		},
	}

	testDb := test_utils.CreateTestDB()
	models.Migrate(testDb)

	testDb.Session(&gorm.Session{FullSaveAssociations: true}).Create(&testBiosynth)

	idFrom := testBiosynth.Modules[0].ID
	idTo := testBiosynth.Modules[1].ID

	err := ReorderEntryBiosynthesisModules(testDb, idFrom, idTo)

	assert.Nil(t, err)

	var actualBioSynth Biosynthesis

	testDb.
		Table("biosyntheses").
		Preload("Modules.Carriers.Location").
		Preload("Modules.ModificationDomains.Location").
		Preload("Modules.ATDomain.Location").
		Preload("Modules.KSDomain.Location").
		First(&actualBioSynth)

	assert.Equal(t, uint64(2), actualBioSynth.Modules[0].Index)
	assert.Equal(t, uint64(1), actualBioSynth.Modules[1].Index)
}
