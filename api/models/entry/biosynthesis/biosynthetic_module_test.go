package biosynthesis

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/util/test_utils"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
)

var inactive = false
var testModule = BiosyntheticModule{
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
	ModificationDomains: []ModificationDomain{
		{
			ID:         1,
			Type:       "o",
			Subtype:    "p",
			Gene:       "q",
			LocationID: 2,
			Location: DomainLocation{
				ID:   2,
				From: 1,
				To:   2,
			},
			Inactive: &inactive,
			Substrates: []DomainSubstrate{
				{
					ID:        1,
					Name:      "r",
					Details:   "s",
					Structure: "t",
				},
			},
			Evidence: []DomainSubstrateEvidence{
				{
					ID:     3,
					Method: "u",
					References: pq.StringArray{
						"w",
						"x",
					},
				},
			},
			References: pq.StringArray{
				"y",
				"z",
			},
			Stereochemistry: pq.StringArray{
				"aa",
				"bb",
			},
			Details: "cc",
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
}

func TestCreateEntryBiosynthesisModule(t *testing.T) {
	testDb := test_utils.CreateTestDB()
	models.Migrate(testDb)

	biosynth := Biosynthesis{
		ID:             1,
		EntryAccession: "test",
	}

	err := testDb.Create(&biosynth).Error

	assert.Nil(t, err)

	err = CreateEntryBiosynthesisModule(testDb, biosynth.EntryAccession, testModule)

	assert.Nil(t, err)

	var actualModule BiosyntheticModule

	err = testDb.
		Table("biosynthetic_modules").
		Where("id = $1", testModule.ID).
		Preload("IntegratedMonomers").
		Preload("IntegratedMonomers.Evidence").
		Preload("Carriers.Location").
		Preload("Carriers.Evidence").
		Preload("ModificationDomains.Location").
		Preload("ModificationDomains.Substrates").
		Preload("ModificationDomains.Evidence").
		Preload("CDomain.Location").
		Preload("ADomain.Location").
		Preload("ADomain.Evidence").
		Preload("ADomain.Substrates").
		Preload("ATDomain.Location").
		Preload("ATDomain.Substrates").
		Preload("ATDomain.Evidence").
		Preload("KSDomain.Location").
		First(&actualModule).
		Error

	assert.Nil(t, err)
	assert.Equal(t, testModule, actualModule)

	testDb.Delete(&BiosyntheticModule{}).Where("id = ?", testModule.ID)
}

func TestGetEntryBiosynthesisModule(t *testing.T) {
	testDb := test_utils.CreateTestDB()
	models.Migrate(testDb)

	biosynth := Biosynthesis{
		ID: 1,
	}

	err := testDb.Create(&biosynth).Error

	assert.Nil(t, err)

	err = testDb.Create(&testModule).Error

	assert.Nil(t, err)

	actualModule, err := GetEntryBiosynthesisModule(testDb, int(testModule.ID))

	assert.Nil(t, err)
	assert.Equal(t, testModule, *actualModule)

	testDb.Where("id = ?", testModule.ID).Delete(&BiosyntheticModule{})
}

func TestUpdateEntryBiosynthesisModule(t *testing.T) {
	testDb := test_utils.CreateTestDB()
	models.Migrate(testDb)

	biosynth := Biosynthesis{
		ID:             1,
		EntryAccession: "test",
	}

	err := testDb.Create(&biosynth).Error

	assert.Nil(t, err)

	err = testDb.Create(&testModule).Error

	assert.Nil(t, err)

	var expectedInactive = true
	var expectedModule = BiosyntheticModule{
		ID:             1,
		Index:          1,
		BiosynthesisID: 1,
		Type:           "a_test",
		Name:           "_test",
		Genes: pq.StringArray{
			"c_test",
			"d_test",
		},
		Active: true,
		IntegratedMonomers: []IntegratedMonomer{
			{
				ID:                   1,
				BiosyntheticModuleID: 1,
				Name:                 "e_test",
				Structure:            "f_test",
				Evidence: []DomainSubstrateEvidence{
					{
						ID:         1,
						Method:     "g_test",
						References: []string{"h_test"},
					},
				},
			},
		},
		Carriers: []CarrierDomain{
			{
				ID:         1,
				Type:       "i_test",
				Subtype:    "j_test",
				Gene:       "k_test",
				LocationID: 1,
				Location: DomainLocation{
					ID:   1,
					From: 3,
					To:   4,
				},
				Inactive:      true,
				BetaBranching: true,
				Evidence: []DomainSubstrateEvidence{
					{
						ID:     2,
						Method: "l_test",
						References: pq.StringArray{
							"m_test",
							"n_test",
						},
					},
				},
			},
		},
		ModificationDomains: []ModificationDomain{
			{
				ID:         1,
				Type:       "o_test",
				Subtype:    "p_test",
				Gene:       "q_test",
				LocationID: 2,
				Location: DomainLocation{
					ID:   2,
					From: 3,
					To:   4,
				},
				Inactive: &expectedInactive,
				Substrates: []DomainSubstrate{
					{
						ID:        1,
						Name:      "r_test",
						Details:   "s_test",
						Structure: "t_test",
					},
				},
				Evidence: []DomainSubstrateEvidence{
					{
						ID:     3,
						Method: "u_test",
						References: pq.StringArray{
							"w_test",
							"x_test",
						},
					},
				},
				References: pq.StringArray{
					"y_test",
					"z_test",
				},
				Stereochemistry: pq.StringArray{
					"aa_test",
					"bb_test",
				},
				Details: "cc_test",
			},
		},
		CDomainID: 1,
		CDomain: &CondensationDomain{
			ID:         1,
			Type:       "dd_test",
			Gene:       "ee_test",
			LocationID: 3,
			Location: DomainLocation{
				ID:   3,
				From: 3,
				To:   4,
			},
			References: pq.StringArray{
				"ff_test",
				"gg_test",
			},
		},
		ADomainID: 1,
		ADomain: &AdenylationDomain{
			ID:         1,
			Type:       "hh_test",
			Gene:       "ii_test",
			LocationID: 4,
			Location: DomainLocation{
				ID:   4,
				From: 3,
				To:   4,
			},
			Inactive: true,
			Evidence: []DomainSubstrateEvidence{
				{
					ID:     4,
					Method: "jj_test",
					References: pq.StringArray{
						"kk_test",
						"ll_test",
					},
				},
			},
			PrecursorBiosynthesis: pq.StringArray{
				"mm_test",
				"nn_test",
			},
			Substrates: []DomainSubstrate{
				{
					ID:        2,
					Name:      "oo_test",
					Details:   "pp_test",
					Structure: "qq_test",
				},
			},
		},
		ATDomainID: 1,
		ATDomain: &AcetyltransferaseDomain{
			ID:         1,
			Type:       "rr_test",
			Subtype:    "ss_test",
			Gene:       "tt_test",
			LocationID: 5,
			Location: DomainLocation{
				ID:   5,
				From: 3,
				To:   4,
			},
			Inactive: true,
			Substrates: []DomainSubstrate{
				{
					ID:        3,
					Name:      "uu_test",
					Details:   "vv_test",
					Structure: "ww_test",
				},
			},
			Evidence: []DomainSubstrateEvidence{
				{
					ID:     5,
					Method: "xx_test",
					References: pq.StringArray{
						"yy_test",
						"zz_test",
					},
				},
			},
		},
		KSDomainID: 1,
		KSDomain: &KetoSynthaseDomain{
			ID:         1,
			Type:       "aaa_test",
			Gene:       "bbb_test",
			LocationID: 6,
			Location: DomainLocation{
				ID:   6,
				From: 3,
				To:   4,
			},
		},
	}

	err = UpdateEntryBiosynthesisModule(testDb, expectedModule)

	assert.Nil(t, err)

	var actualModule BiosyntheticModule

	err = testDb.
		Table("biosynthetic_modules").
		Where("id = $1", testModule.ID).
		Preload("IntegratedMonomers").
		Preload("IntegratedMonomers.Evidence").
		Preload("Carriers.Location").
		Preload("Carriers.Evidence").
		Preload("ModificationDomains.Location").
		Preload("ModificationDomains.Substrates").
		Preload("ModificationDomains.Evidence").
		Preload("CDomain.Location").
		Preload("ADomain.Location").
		Preload("ADomain.Evidence").
		Preload("ADomain.Substrates").
		Preload("ATDomain.Location").
		Preload("ATDomain.Substrates").
		Preload("ATDomain.Evidence").
		Preload("KSDomain.Location").
		First(&actualModule).
		Error

	assert.Equal(t, expectedModule, actualModule)

	testDb.Where("id = ?", testModule.ID).Delete(&BiosyntheticModule{})
}

func TestDeleteEntryBiosynthesisModule(t *testing.T) {

}
