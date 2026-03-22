package biosynthesis

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/util/test_utils"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testBiosynthDomainBiosynthData = Biosynthesis{
	ID:             1,
	EntryAccession: "test",
	Classes:        nil,
	Modules: []BiosyntheticModule{
		{
			ID:             1,
			Index:          1,
			BiosynthesisID: 1,
			Type:           "a",
			Name:           "b",
			Genes: pq.StringArray{
				"c",
				"d",
			},
			Active: false,
			IntegratedMonomers: []IntegratedMonomer{
				{
					ID:                   1,
					BiosyntheticModuleID: 1,
					Name:                 "e",
					Structure:            "f",
					Evidence: []DomainSubstrateEvidence{
						{
							ID:         1,
							Method:     "g",
							References: []string{"h"},
						},
					},
				},
			},
			Carriers: []CarrierDomain{
				{
					ID:         1,
					Type:       "i",
					Subtype:    "j",
					Gene:       "k",
					LocationID: 1,
					Location: DomainLocation{
						ID:   1,
						From: 1,
						To:   2,
					},
					Inactive:      false,
					BetaBranching: false,
					Evidence: []DomainSubstrateEvidence{
						{
							ID:     2,
							Method: "l",
							References: pq.StringArray{
								"m",
								"n",
							},
						},
					},
				},
			},
			CDomainID: 1,
			CDomain: &CondensationDomain{
				ID:         1,
				Type:       "dd",
				Gene:       "ee",
				LocationID: 3,
				Location: DomainLocation{
					ID:   3,
					From: 1,
					To:   2,
				},
				References: pq.StringArray{
					"ff",
					"gg",
				},
			},
			ADomainID: 1,
			ADomain: &AdenylationDomain{
				ID:         1,
				Type:       "hh",
				Gene:       "ii",
				LocationID: 4,
				Location: DomainLocation{
					ID:   4,
					From: 1,
					To:   2,
				},
				Inactive: false,
				Evidence: []DomainSubstrateEvidence{
					{
						ID:     4,
						Method: "jj",
						References: pq.StringArray{
							"kk",
							"ll",
						},
					},
				},
				PrecursorBiosynthesis: pq.StringArray{
					"mm",
					"nn",
				},
				Substrates: []DomainSubstrate{
					{
						ID:        2,
						Name:      "oo",
						Details:   "pp",
						Structure: "qq",
					},
				},
			},
			ATDomainID: 1,
			ATDomain: &AcetyltransferaseDomain{
				ID:         1,
				Type:       "rr",
				Subtype:    "ss",
				Gene:       "tt",
				LocationID: 5,
				Location: DomainLocation{
					ID:   5,
					From: 1,
					To:   2,
				},
				Inactive: false,
				Substrates: []DomainSubstrate{
					{
						ID:        3,
						Name:      "uu",
						Details:   "vv",
						Structure: "ww",
					},
				},
				Evidence: []DomainSubstrateEvidence{
					{
						ID:     5,
						Method: "xx",
						References: pq.StringArray{
							"yy",
							"zz",
						},
					},
				},
			},
			KSDomainID: 1,
			KSDomain: &KetoSynthaseDomain{
				ID:         1,
				Type:       "aaa",
				Gene:       "bbb",
				LocationID: 6,
				Location: DomainLocation{
					ID:   6,
					From: 1,
					To:   2,
				},
			},
		},
	},
	Operons: nil,
	Paths:   nil,
}

var testInactive = true
var testBiosynthDomainTestData = ModificationDomain{
	ID:         1,
	Type:       "a",
	Subtype:    "b",
	Gene:       "c",
	LocationID: 1,
	Location: DomainLocation{
		ID:   1,
		From: 1,
		To:   2,
	},
	Inactive: &testInactive,
	Substrates: []DomainSubstrate{
		{
			ID:        1,
			Name:      "d",
			Details:   "e",
			Structure: "f",
		},
	},
	Evidence: []DomainSubstrateEvidence{
		{
			ID:         1,
			Method:     "g",
			References: []string{"h"},
		},
	},
	References:      []string{"i"},
	Stereochemistry: []string{"j"},
	Details:         "k",
}

func TestCreateModificationDomain(t *testing.T) {
	testDb := test_utils.CreateTestDB()
	models.Migrate(testDb)

	err := testDb.Create(&testBiosynthDomainBiosynthData).Error
	assert.Nil(t, err)

	err = CreateModificationDomain(testDb, int(testBiosynthDomainBiosynthData.Modules[0].ID), testBiosynthDomainTestData)

	assert.Nil(t, err)

	var actualModificationDomain ModificationDomain
	err = testDb.
		Table("modification_domains").
		Preload("Location").
		Preload("Substrates").
		Preload("Evidence").
		First(&actualModificationDomain).
		Error

	assert.Nil(t, err)

	assert.Equal(t, testBiosynthDomainTestData, actualModificationDomain)
}

func TestGetModificationDomain(t *testing.T) {
	testDb := test_utils.CreateTestDB()
	models.Migrate(testDb)

	err := testDb.Create(&testBiosynthDomainBiosynthData).Error
	assert.Nil(t, err)

	err = testDb.Create(&testBiosynthDomainTestData).Error
	assert.Nil(t, err)

	expectedModificationDomain := testBiosynthDomainTestData

	actualModificationDomain, err := GetModificationDomain(testDb, int(expectedModificationDomain.ID))

	assert.Nil(t, err)
	assert.Equal(t, expectedModificationDomain, *actualModificationDomain)
}

func TestUpdateModificationDomain(t *testing.T) {
	testDb := test_utils.CreateTestDB()
	models.Migrate(testDb)

	err := testDb.Create(&testBiosynthDomainBiosynthData).Error
	assert.Nil(t, err)

	err = testDb.Create(&testBiosynthDomainTestData).Error
	assert.Nil(t, err)

	updateInactive := false
	modificationDomainUpdate := ModificationDomain{
		ID:         1,
		Type:       "a_test",
		Subtype:    "b_test",
		Gene:       "c_test",
		LocationID: 1,
		Location: DomainLocation{
			ID:   1,
			From: 3,
			To:   4,
		},
		Inactive: &updateInactive,
		Substrates: []DomainSubstrate{
			{
				ID:        1,
				Name:      "d_test",
				Details:   "e_test",
				Structure: "f_test",
			},
		},
		Evidence: []DomainSubstrateEvidence{
			{
				ID:         1,
				Method:     "g_test",
				References: []string{"h_test"},
			},
		},
		References:      []string{"i_test"},
		Stereochemistry: []string{"j_test"},
		Details:         "k_test",
	}

	err = UpdateModificationDomain(testDb, modificationDomainUpdate)

	assert.Nil(t, err)

	var actualModificationDomain ModificationDomain
	err = testDb.
		Table("modification_domains").
		Preload("Location").
		Preload("Substrates").
		Preload("Evidence").
		First(&actualModificationDomain).
		Error

	assert.Nil(t, err)
	assert.Equal(t, modificationDomainUpdate, actualModificationDomain)
}

func TestDeleteModificationDomain(t *testing.T) {
	testDb := test_utils.CreateTestDB()
	models.Migrate(testDb)

	err := testDb.Create(&testBiosynthDomainBiosynthData).Error
	assert.Nil(t, err)

	err = testDb.Create(&testBiosynthDomainTestData).Error
	assert.Nil(t, err)

	err = DeleteModificationDomain(testDb, int(testBiosynthDomainTestData.ID))

	assert.Nil(t, err)

	var actualModificationDomain ModificationDomain
	err = testDb.
		Table("modification_domains").
		Preload("Location").
		Preload("Substrates").
		Preload("Evidence").
		Find(&actualModificationDomain).
		Error

	assert.Equal(t, actualModificationDomain.ID, uint64(0))
}
