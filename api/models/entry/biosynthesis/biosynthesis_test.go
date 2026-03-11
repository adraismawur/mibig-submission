package biosynthesis

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/util/test_utils"
	"github.com/lib/pq"
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

func TestCreateNRPSClass(t *testing.T) {
	biosynth := Biosynthesis{
		ID: 1,
	}

	testClass := BiosyntheticClass{
		ID:             1,
		BiosynthesisID: 1,
		Class:          "NRPS",
		Subclass:       "Type I",
		ReleaseTypes: &[]ReleaseType{
			{
				ID:                  1,
				BiosyntheticClassID: 1,
				Name:                "Macrolactamization",
				Details:             "a",
				References:          pq.StringArray{"doi:pending"},
			},
			{
				ID:                  2,
				BiosyntheticClassID: 1,
				Name:                "Reductive release",
				Details:             "b",
				References:          pq.StringArray{"doi:pending"},
			},
		},
		Thioesterases: &[]Thioesterase{
			{
				ID:                  3,
				BiosyntheticClassID: 1,
				Type:                "thioesterase",
				LocationID:          1,
				Location: DomainLocation{
					From: 1,
					To:   2,
				},
				Subtype: "Type I",
			},
			{
				ID:                  4,
				BiosyntheticClassID: 1,
				Type:                "thioesterase",
				LocationID:          2,
				Location: DomainLocation{
					From: 3,
					To:   4,
				},
				Subtype: "Type II",
			},
		},
	}

	testDb := test_utils.CreateTestDB()
	models.Migrate(testDb)
	testDb.Create(&biosynth)

	err := CreateBiosynthesisClass(testDb, biosynth.ID, testClass)

	assert.Nil(t, err)

	var actualBiosynthClass BiosyntheticClass

	err = testDb.
		Preload("ReleaseTypes").
		Preload("Thioesterases.Location").
		First(&actualBiosynthClass).
		Error

	assert.Nil(t, err)

	assert.Equal(t, testClass, actualBiosynthClass)
}

func TestCreatePKSClass(t *testing.T) {
	biosynth := Biosynthesis{
		ID: 1,
	}

	ketideLength := 3
	testClass := BiosyntheticClass{
		ID:             1,
		BiosynthesisID: 1,
		Class:          "PKS",
		Subclass:       "Type III",
		Cyclases:       pq.StringArray{"a", "b"},
		KetideLength:   &ketideLength,
	}

	testDb := test_utils.CreateTestDB()
	models.Migrate(testDb)
	testDb.Create(&biosynth)

	err := CreateBiosynthesisClass(testDb, biosynth.ID, testClass)

	assert.Nil(t, err)

	var actualBiosynthClass BiosyntheticClass

	err = testDb.
		First(&actualBiosynthClass).
		Error

	assert.Nil(t, err)

	assert.Equal(t, testClass, actualBiosynthClass)
}

func TestCreateRibosomalClass(t *testing.T) {
	biosynth := Biosynthesis{
		ID: 1,
	}

	rippType := "Borosin"
	testClass := BiosyntheticClass{
		ID:             1,
		BiosynthesisID: 1,
		Class:          "ribosomal",
		Subclass:       "RiPP",
		RIPPType:       &rippType,
		Peptidases:     pq.StringArray{"a", "b", "c"},
		Precursors: &[]RippPrecursor{
			{
				ID:                       1,
				BiosyntheticClassID:      1,
				Gene:                     "d",
				LeaderCleavageLocationID: 1,
				LeaderCleavageLocation: &CleavageLocation{
					ID:   1,
					From: 1,
					To:   2,
				},
				FollowerCleavageLocationID: 2,
				FollowerCleavageLocation: &CleavageLocation{
					ID:   2,
					From: 3,
					To:   4,
				},
				Crosslinks: []RippPrecursorCrosslink{
					{
						ID:              1,
						RippPrecursorID: 1,
						From:            5,
						To:              6,
						Type:            "ether",
						Details:         "test",
					},
				},
			},
		},
	}

	testDb := test_utils.CreateTestDB()
	models.Migrate(testDb)
	testDb.Create(&biosynth)

	err := CreateBiosynthesisClass(testDb, biosynth.ID, testClass)

	assert.Nil(t, err)

	var actualBiosynthClass BiosyntheticClass

	err = testDb.
		Preload("Precursors.LeaderCleavageLocation").
		Preload("Precursors.FollowerCleavageLocation").
		Preload("Precursors.Crosslinks").
		First(&actualBiosynthClass).
		Error

	assert.Nil(t, err)

	assert.Equal(t, testClass, actualBiosynthClass)
}

func TestCreateSaccharideClass(t *testing.T) {
	biosynth := Biosynthesis{
		ID: 1,
	}

	testClass := BiosyntheticClass{
		ID:             1,
		BiosynthesisID: 1,
		Class:          "saccharide",
		Subclass:       "",
		GlycosylTransferases: &[]GlycosylTransferase{
			{
				ID:                  1,
				BiosyntheticClassID: 1,
				Evidence: &[]GlycosylTransferaseEvidence{
					{
						ID:                    1,
						GlycosylTransferaseID: 1,
						Name:                  "a",
						References:            pq.StringArray{"a", "b"},
					},
				},
				References:  pq.StringArray{"a", "b"},
				Gene:        "a",
				Specificity: "C1CCCCC1",
			},
		},
		Subclusters: &[]SaccharideSubcluster{
			{
				ID:                  1,
				BiosyntheticClassID: 1,
				Specificity:         "C1CCCCC1",
				Genes:               pq.StringArray{"a", "b"},
				References:          pq.StringArray{"a", "b"},
			},
		},
	}

	testDb := test_utils.CreateTestDB()
	models.Migrate(testDb)
	testDb.Create(&biosynth)

	err := CreateBiosynthesisClass(testDb, biosynth.ID, testClass)

	assert.Nil(t, err)

	var actualBiosynthClass BiosyntheticClass

	err = testDb.
		Preload("GlycosylTransferases.Evidence").
		Preload("Subclusters").
		First(&actualBiosynthClass).
		Error

	assert.Nil(t, err)

	assert.Equal(t, testClass, actualBiosynthClass)
}

func TestCreateTerpeneClass(t *testing.T) {
	biosynth := Biosynthesis{
		ID: 1,
	}

	precursor := "GGPP"
	testClass := BiosyntheticClass{
		ID:                 1,
		BiosynthesisID:     1,
		Class:              "terpene",
		Subclass:           "Hemiterpene",
		Prenyltransferases: pq.StringArray{"a"},
		SynthasesCyclases:  pq.StringArray{"b"},
		Precursor:          &precursor,
	}

	testDb := test_utils.CreateTestDB()
	models.Migrate(testDb)
	testDb.Create(&biosynth)

	err := CreateBiosynthesisClass(testDb, biosynth.ID, testClass)

	assert.Nil(t, err)

	var actualBiosynthClass BiosyntheticClass

	err = testDb.
		First(&actualBiosynthClass).
		Error

	assert.Nil(t, err)

	assert.Equal(t, testClass, actualBiosynthClass)
}

func TestCreateOtherClass(t *testing.T) {
	biosynth := Biosynthesis{
		ID: 1,
	}

	classDetails := "test"
	testClass := BiosyntheticClass{
		ID:             1,
		BiosynthesisID: 1,
		Class:          "other",
		Subclass:       "cyclitol",
		Details:        &classDetails,
	}

	testDb := test_utils.CreateTestDB()
	models.Migrate(testDb)
	testDb.Create(&biosynth)

	err := CreateBiosynthesisClass(testDb, biosynth.ID, testClass)

	assert.Nil(t, err)

	var actualBiosynthClass BiosyntheticClass

	err = testDb.
		First(&actualBiosynthClass).
		Error

	assert.Nil(t, err)

	assert.Equal(t, testClass, actualBiosynthClass)
}

func TestGetNRPSClass(t *testing.T) {
	biosynth := Biosynthesis{
		ID: 1,
	}

	testClass := BiosyntheticClass{
		ID:             1,
		BiosynthesisID: 1,
		Class:          "NRPS",
		Subclass:       "Type I",
		ReleaseTypes: &[]ReleaseType{
			{
				ID:                  1,
				BiosyntheticClassID: 1,
				Name:                "Macrolactamization",
				Details:             "a",
				References:          pq.StringArray{"doi:pending"},
			},
			{
				ID:                  2,
				BiosyntheticClassID: 1,
				Name:                "Reductive release",
				Details:             "b",
				References:          pq.StringArray{"doi:pending"},
			},
		},
		Thioesterases: &[]Thioesterase{
			{
				ID:                  3,
				BiosyntheticClassID: 1,
				Type:                "thioesterase",
				LocationID:          1,
				Location: DomainLocation{
					From: 1,
					To:   2,
				},
				Subtype: "Type I",
			},
			{
				ID:                  4,
				BiosyntheticClassID: 1,
				Type:                "thioesterase",
				LocationID:          2,
				Location: DomainLocation{
					From: 3,
					To:   4,
				},
				Subtype: "Type II",
			},
		},
	}

	testDb := test_utils.CreateTestDB()
	models.Migrate(testDb)
	testDb.Create(&biosynth)
	testDb.Session(&gorm.Session{FullSaveAssociations: true}).Create(&testClass)

	actualBiosynthClass, err := GetEntryBiosynthesisClass(testDb, int(testClass.ID))

	assert.Nil(t, err)
	assert.Equal(t, &testClass, actualBiosynthClass)
}

func TestGetPKSClass(t *testing.T) {
	biosynth := Biosynthesis{
		ID: 1,
	}

	ketideLength := 3
	testClass := BiosyntheticClass{
		ID:             1,
		BiosynthesisID: 1,
		Class:          "PKS",
		Subclass:       "Type III",
		Cyclases:       pq.StringArray{"a", "b"},
		KetideLength:   &ketideLength,
	}

	testDb := test_utils.CreateTestDB()
	models.Migrate(testDb)
	testDb.Create(&biosynth)
	testDb.Create(&testClass)

	actualBiosynthClass, err := GetEntryBiosynthesisClass(testDb, int(testClass.ID))

	assert.Nil(t, err)
	assert.Equal(t, &testClass, actualBiosynthClass)
}

func TestGetRibosomalClass(t *testing.T) {
	biosynth := Biosynthesis{
		ID: 1,
	}

	rippType := "Borosin"
	testClass := BiosyntheticClass{
		ID:             1,
		BiosynthesisID: 1,
		Class:          "ribosomal",
		Subclass:       "RiPP",
		RIPPType:       &rippType,
		Peptidases:     pq.StringArray{"a", "b", "c"},
		Precursors: &[]RippPrecursor{
			{
				ID:                       1,
				BiosyntheticClassID:      1,
				Gene:                     "d",
				LeaderCleavageLocationID: 1,
				LeaderCleavageLocation: &CleavageLocation{
					ID:   1,
					From: 1,
					To:   2,
				},
				FollowerCleavageLocationID: 2,
				FollowerCleavageLocation: &CleavageLocation{
					ID:   2,
					From: 3,
					To:   4,
				},
				Crosslinks: []RippPrecursorCrosslink{
					{
						ID:              1,
						RippPrecursorID: 1,
						From:            5,
						To:              6,
						Type:            "ether",
						Details:         "test",
					},
				},
			},
		},
	}

	testDb := test_utils.CreateTestDB()
	models.Migrate(testDb)
	testDb.Create(&biosynth)
	testDb.Session(&gorm.Session{FullSaveAssociations: true}).Create(&testClass)

	actualBiosynthClass, err := GetEntryBiosynthesisClass(testDb, int(testClass.ID))

	assert.Nil(t, err)
	assert.Equal(t, &testClass, actualBiosynthClass)
}

func TestGetSaccharideClass(t *testing.T) {
	biosynth := Biosynthesis{
		ID: 1,
	}

	testClass := BiosyntheticClass{
		ID:             1,
		BiosynthesisID: 1,
		Class:          "saccharide",
		Subclass:       "",
		GlycosylTransferases: &[]GlycosylTransferase{
			{
				ID:                  1,
				BiosyntheticClassID: 1,
				Evidence: &[]GlycosylTransferaseEvidence{
					{
						ID:                    1,
						GlycosylTransferaseID: 1,
						Name:                  "a",
						References:            pq.StringArray{"a", "b"},
					},
				},
				References:  pq.StringArray{"a", "b"},
				Gene:        "a",
				Specificity: "C1CCCCC1",
			},
		},
		Subclusters: &[]SaccharideSubcluster{
			{
				ID:                  1,
				BiosyntheticClassID: 1,
				Specificity:         "C1CCCCC1",
				Genes:               pq.StringArray{"d", "d"},
				References:          pq.StringArray{"e", "f"},
			},
		},
	}

	testDb := test_utils.CreateTestDB()
	models.Migrate(testDb)
	testDb.Create(&biosynth)
	testDb.Session(&gorm.Session{FullSaveAssociations: true}).Create(&testClass)

	actualBiosynthClass, err := GetEntryBiosynthesisClass(testDb, int(testClass.ID))

	assert.Nil(t, err)
	assert.Equal(t, &testClass, actualBiosynthClass)
}

func TestGetTerpeneClass(t *testing.T) {
	biosynth := Biosynthesis{
		ID: 1,
	}

	precursor := "GGPP"
	testClass := BiosyntheticClass{
		ID:                 1,
		BiosynthesisID:     1,
		Class:              "terpene",
		Subclass:           "Hemiterpene",
		Prenyltransferases: pq.StringArray{"a"},
		SynthasesCyclases:  pq.StringArray{"b"},
		Precursor:          &precursor,
	}

	testDb := test_utils.CreateTestDB()
	models.Migrate(testDb)
	testDb.Create(&biosynth)
	testDb.Create(&testClass)

	actualBiosynthClass, err := GetEntryBiosynthesisClass(testDb, int(testClass.ID))

	assert.Nil(t, err)
	assert.Equal(t, &testClass, actualBiosynthClass)
}

func TestGetOtherClass(t *testing.T) {
	biosynth := Biosynthesis{
		ID: 1,
	}

	classDetails := "test"
	testClass := BiosyntheticClass{
		ID:             1,
		BiosynthesisID: 1,
		Class:          "other",
		Subclass:       "cyclitol",
		Details:        &classDetails,
	}

	testDb := test_utils.CreateTestDB()
	models.Migrate(testDb)
	testDb.Create(&biosynth)
	testDb.Create(&testClass)

	actualBiosynthClass, err := GetEntryBiosynthesisClass(testDb, int(testClass.ID))

	assert.Nil(t, err)
	assert.Equal(t, &testClass, actualBiosynthClass)
}

func TestUpdateNRPSClass(t *testing.T) {
	biosynth := Biosynthesis{
		ID: 1,
	}

	testClass := BiosyntheticClass{
		ID:             1,
		BiosynthesisID: 1,
		Class:          "NRPS",
		Subclass:       "Type I",
		ReleaseTypes: &[]ReleaseType{
			{
				ID:                  1,
				BiosyntheticClassID: 1,
				Name:                "Macrolactamization",
				Details:             "a",
				References:          pq.StringArray{"doi:pending"},
			},
			{
				ID:                  2,
				BiosyntheticClassID: 1,
				Name:                "Reductive release",
				Details:             "b",
				References:          pq.StringArray{"doi:pending"},
			},
		},
		Thioesterases: &[]Thioesterase{
			{
				ID:                  3,
				BiosyntheticClassID: 1,
				Type:                "thioesterase",
				LocationID:          1,
				Location: DomainLocation{
					From: 1,
					To:   2,
				},
				Subtype: "Type I",
			},
			{
				ID:                  4,
				BiosyntheticClassID: 1,
				Type:                "thioesterase",
				LocationID:          2,
				Location: DomainLocation{
					From: 3,
					To:   4,
				},
				Subtype: "Type II",
			},
		},
	}

	testDb := test_utils.CreateTestDB()
	models.Migrate(testDb)
	testDb.Create(&biosynth)
	testDb.Session(&gorm.Session{FullSaveAssociations: true}).Create(&testClass)

	classUpdate := BiosyntheticClass{
		ID:             1,
		BiosynthesisID: 1,
		Class:          "NRPS",
		Subclass:       "Type V",
		ReleaseTypes: &[]ReleaseType{
			{
				ID:                  1,
				BiosyntheticClassID: 1,
				Name:                "Hydrolysis",
				Details:             "c",
				References:          pq.StringArray{"e"},
			},
			{
				ID:                  2,
				BiosyntheticClassID: 1,
				Name:                "Claisen condensation",
				Details:             "d",
				References:          pq.StringArray{"f"},
			},
		},
		Thioesterases: &[]Thioesterase{
			{
				ID:                  3,
				BiosyntheticClassID: 1,
				Type:                "thioesterase",
				LocationID:          1,
				Location: DomainLocation{
					From: 5,
					To:   6,
				},
				Subtype: "Type II",
			},
			{
				ID:                  4,
				BiosyntheticClassID: 1,
				Type:                "thioesterase",
				LocationID:          2,
				Location: DomainLocation{
					From: 7,
					To:   8,
				},
				Subtype: "Type I",
			},
		},
	}

	err := UpdateEntryBiosynthesisClass(testDb, int(testClass.ID), classUpdate)

	var actualBiosynthClass BiosyntheticClass
	testDb.
		Model(&BiosyntheticClass{}).
		Preload("ReleaseTypes").
		Preload("Thioesterases.Location").
		First(&actualBiosynthClass)

	assert.Nil(t, err)
	assert.Equal(t, classUpdate, actualBiosynthClass)
}

func TestUpdatePKSClass(t *testing.T) {
	biosynth := Biosynthesis{
		ID: 1,
	}

	ketideLength := 3
	testClass := BiosyntheticClass{
		ID:             1,
		BiosynthesisID: 1,
		Class:          "PKS",
		Subclass:       "Type III",
		Cyclases:       pq.StringArray{"a", "b"},
		KetideLength:   &ketideLength,
	}

	testDb := test_utils.CreateTestDB()
	models.Migrate(testDb)
	testDb.Create(&biosynth)
	testDb.Create(&testClass)

	ketideLengthUpdate := 6
	classUpdate := BiosyntheticClass{
		ID:             1,
		BiosynthesisID: 1,
		Class:          "PKS",
		Subclass:       "Type II highly reducing",
		Cyclases:       pq.StringArray{"c", "d"},
		KetideLength:   &ketideLengthUpdate,
	}
	err := UpdateEntryBiosynthesisClass(testDb, int(testClass.ID), classUpdate)

	var actualBiosynthClass BiosyntheticClass
	testDb.Model(&BiosyntheticClass{}).First(&actualBiosynthClass)

	assert.Nil(t, err)
	assert.Equal(t, classUpdate, actualBiosynthClass)
}

func TestUpdateRibosomalClass(t *testing.T) {
	biosynth := Biosynthesis{
		ID: 1,
	}

	rippType := "Borosin"
	testClass := BiosyntheticClass{
		ID:             1,
		BiosynthesisID: 1,
		Class:          "ribosomal",
		Subclass:       "RiPP",
		RIPPType:       &rippType,
		Peptidases:     pq.StringArray{"a", "b", "c"},
		Precursors: &[]RippPrecursor{
			{
				ID:                       1,
				BiosyntheticClassID:      1,
				Gene:                     "d",
				LeaderCleavageLocationID: 1,
				LeaderCleavageLocation: &CleavageLocation{
					ID:   1,
					From: 1,
					To:   2,
				},
				FollowerCleavageLocationID: 2,
				FollowerCleavageLocation: &CleavageLocation{
					ID:   2,
					From: 3,
					To:   4,
				},
				Crosslinks: []RippPrecursorCrosslink{
					{
						ID:              1,
						RippPrecursorID: 1,
						From:            5,
						To:              6,
						Type:            "ether",
						Details:         "test",
					},
				},
			},
		},
	}

	testDb := test_utils.CreateTestDB()
	models.Migrate(testDb)
	testDb.Create(&biosynth)
	testDb.Session(&gorm.Session{FullSaveAssociations: true}).Create(&testClass)

	rippTypeUpdate := "Epipeptide"
	classUpdate := BiosyntheticClass{
		ID:             1,
		BiosynthesisID: 1,
		Class:          "ribosomal",
		Subclass:       "unmodified",
		RIPPType:       &rippTypeUpdate,
		Peptidases:     pq.StringArray{"e", "f", "g"},
		Precursors: &[]RippPrecursor{
			{
				ID:                       1,
				BiosyntheticClassID:      1,
				Gene:                     "h",
				LeaderCleavageLocationID: 1,
				LeaderCleavageLocation: &CleavageLocation{
					ID:   1,
					From: 6,
					To:   7,
				},
				FollowerCleavageLocationID: 2,
				FollowerCleavageLocation: &CleavageLocation{
					ID:   2,
					From: 8,
					To:   9,
				},
				Crosslinks: []RippPrecursorCrosslink{
					{
						ID:              1,
						RippPrecursorID: 1,
						From:            10,
						To:              11,
						Type:            "other",
						Details:         "testupdate",
					},
				},
			},
		},
	}

	err := UpdateEntryBiosynthesisClass(testDb, int(testClass.ID), classUpdate)

	var actualBiosynthClass BiosyntheticClass
	testDb.
		Model(&BiosyntheticClass{}).
		Preload("Precursors.LeaderCleavageLocation").
		Preload("Precursors.FollowerCleavageLocation").
		Preload("Precursors.Crosslinks").
		First(&actualBiosynthClass)

	assert.Nil(t, err)
	assert.Equal(t, classUpdate, actualBiosynthClass)
}

func TestUpdateSaccharideClass(t *testing.T) {
	biosynth := Biosynthesis{
		ID: 1,
	}

	testClass := BiosyntheticClass{
		ID:             1,
		BiosynthesisID: 1,
		Class:          "saccharide",
		Subclass:       "",
		GlycosylTransferases: &[]GlycosylTransferase{
			{
				ID:                  1,
				BiosyntheticClassID: 1,
				Evidence: &[]GlycosylTransferaseEvidence{
					{
						ID:                    1,
						GlycosylTransferaseID: 1,
						Name:                  "a",
						References:            pq.StringArray{"b", "c"},
					},
				},
				References:  pq.StringArray{"d", "e"},
				Gene:        "f",
				Specificity: "C1CCCCC1",
			},
		},
		Subclusters: &[]SaccharideSubcluster{
			{
				ID:                  1,
				BiosyntheticClassID: 1,
				Specificity:         "C1CCCCC1",
				Genes:               pq.StringArray{"g", "h"},
				References:          pq.StringArray{"i", "j"},
			},
		},
	}

	testDb := test_utils.CreateTestDB()
	models.Migrate(testDb)
	testDb.Create(&biosynth)
	testDb.Session(&gorm.Session{FullSaveAssociations: true}).Create(&testClass)

	classUpdate := BiosyntheticClass{
		ID:             1,
		BiosynthesisID: 1,
		Class:          "saccharide",
		Subclass:       "",
		GlycosylTransferases: &[]GlycosylTransferase{
			{
				ID:                  1,
				BiosyntheticClassID: 1,
				Evidence: &[]GlycosylTransferaseEvidence{
					{
						ID:                    1,
						GlycosylTransferaseID: 1,
						Name:                  "k",
						References:            pq.StringArray{"l", "m"},
					},
				},
				References:  pq.StringArray{"n", "o"},
				Gene:        "p",
				Specificity: "O=C(O)C[C@H](N)C(=O)N[C@H](C(=O)OC)Cc1ccccc1",
			},
		},
		Subclusters: &[]SaccharideSubcluster{
			{
				ID:                  1,
				BiosyntheticClassID: 1,
				Specificity:         "O=C(O)C[C@H](N)C(=O)N[C@H](C(=O)OC)Cc1ccccc1",
				Genes:               pq.StringArray{"q", "r"},
				References:          pq.StringArray{"s", "t"},
			},
		},
	}

	err := UpdateEntryBiosynthesisClass(testDb, int(testClass.ID), classUpdate)

	var actualBiosynthClass BiosyntheticClass
	testDb.
		Model(&BiosyntheticClass{}).
		Preload("GlycosylTransferases.Evidence").
		Preload("Subclusters").
		First(&actualBiosynthClass)

	assert.Nil(t, err)
	assert.Equal(t, classUpdate, actualBiosynthClass)
}

func TestUpdateTerpeneClass(t *testing.T) {
	biosynth := Biosynthesis{
		ID: 1,
	}

	precursor := "GGPP"
	testClass := BiosyntheticClass{
		ID:                 1,
		BiosynthesisID:     1,
		Class:              "terpene",
		Subclass:           "Hemiterpene",
		Prenyltransferases: pq.StringArray{"a"},
		SynthasesCyclases:  pq.StringArray{"b"},
		Precursor:          &precursor,
	}

	testDb := test_utils.CreateTestDB()
	models.Migrate(testDb)
	testDb.Create(&biosynth)
	testDb.Create(&testClass)

	newPrecursor := "IPP"
	classUpdate := BiosyntheticClass{
		ID:                 1,
		BiosynthesisID:     1,
		Class:              "terpene",
		Subclass:           "Monoterpene",
		Prenyltransferases: pq.StringArray{"c"},
		SynthasesCyclases:  pq.StringArray{"d"},
		Precursor:          &newPrecursor,
	}

	err := UpdateEntryBiosynthesisClass(testDb, int(testClass.ID), classUpdate)

	var actualBiosynthClass BiosyntheticClass
	testDb.Model(&BiosyntheticClass{}).First(&actualBiosynthClass)

	assert.Nil(t, err)
	assert.Equal(t, classUpdate, actualBiosynthClass)
}

func TestUpdateOtherClass(t *testing.T) {
	biosynth := Biosynthesis{
		ID: 1,
	}

	classDetails := "test"
	testClass := BiosyntheticClass{
		ID:             1,
		BiosynthesisID: 1,
		Class:          "other",
		Subclass:       "cyclitol",
		Details:        &classDetails,
	}

	testDb := test_utils.CreateTestDB()
	models.Migrate(testDb)
	testDb.Create(&biosynth)
	testDb.Create(&testClass)

	newClassDetails := "wow"
	classUpdate := BiosyntheticClass{
		ID:             1,
		BiosynthesisID: 1,
		Class:          "other",
		Subclass:       "aminocoumarin",
		Details:        &newClassDetails,
	}

	err := UpdateEntryBiosynthesisClass(testDb, int(testClass.ID), classUpdate)

	var actualBiosynthClass BiosyntheticClass
	testDb.Model(&BiosyntheticClass{}).First(&actualBiosynthClass)

	assert.Nil(t, err)
	assert.Equal(t, classUpdate, actualBiosynthClass)
}
