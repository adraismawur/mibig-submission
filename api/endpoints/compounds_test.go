package endpoints

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/models/entry"
	"github.com/adraismawur/mibig-submission/models/entry/biosynthesis"
	"github.com/adraismawur/mibig-submission/models/entry/compound"
	"github.com/adraismawur/mibig-submission/models/entry/consts"
	"github.com/adraismawur/mibig-submission/models/entry/taxonomy"
	"github.com/adraismawur/mibig-submission/util"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"testing"
)

var testEntry = entry.Entry{
	ID:           0,
	Accession:    "test",
	Version:      0,
	Changelog:    entry.Changelog{},
	Quality:      consts.Questionable,
	Status:       consts.Active,
	Completeness: consts.Unknown,
	Loci:         nil,
	Biosynthesis: biosynthesis.Biosynthesis{},
	Compounds: []compound.Compound{
		{
			ID:      1,
			EntryID: 1,
			Name:    "benzene",
			Evidence: []compound.CompoundEvidence{
				{
					ID:         1,
					CompoundID: 1,
					Method:     "test",
					References: []string{
						"doi:test",
					},
				},
			},
			BioActivities: []compound.BioActivities{
				{
					ID:         1,
					CompoundID: 1,
					Observed:   true,
					References: []string{
						"doi:test",
					},
				},
			},
			Structure: "c1ccccc1",
			DatabaseIDs: []string{
				"test",
			},
			Moieties: []string{
				"test",
			},
			Mass:    123,
			Formula: "C6H6",
		},
	},
	Taxonomy:         taxonomy.Taxonomy{},
	GeneInformation:  nil,
	LegacyReferences: nil,
	Embargo:          false,
}

func TestGetEntryCompounds(t *testing.T) {
	testDb.Session(&gorm.Session{FullSaveAssociations: true}).Create(&testEntry)

	c, r := util.CreateTestGinGetRequest("/entry/test/compounds")
	c.Params = []gin.Param{
		{
			Key:   "accession",
			Value: testEntry.Accession,
		},
	}

	user := models.User{
		ID:    1,
		Email: models.GenerateRandomEmail(),
		Roles: []models.UserRole{
			{
				Role: models.Admin,
			},
		},
	}

	testToken, _ := models.GenerateToken(user)

	util.AddTokenToHeader(c, testToken)

	getEntryCompounds(testDb, c)

	assert.Equal(t, http.StatusOK, c.Writer.Status())

	var response struct {
		ActualCompounds []compound.Compound `json:"compounds"`
	}

	err := json.Unmarshal(r.Body.Bytes(), &response)

	assert.Nil(t, err)
	assert.Equal(t, testEntry.Compounds, response.ActualCompounds)
}
