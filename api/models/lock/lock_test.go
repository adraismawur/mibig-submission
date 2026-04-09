package lock

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/models/entry"
	"github.com/adraismawur/mibig-submission/util/test_utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestEntryCanCreateLockTrueNoLocks(t *testing.T) {
	db := test_utils.CreateTestDB()
	models.Migrate(db)

	testEntry := entry.Entry{
		Accession: "test",
	}

	err := db.Create(&testEntry).Error

	assert.Nil(t, err)

	canCreateLock, err := EntryCanCreateLock(db, testEntry.Accession, Locitax)

	assert.Nil(t, err)
	assert.True(t, canCreateLock)
}

func TestEntryCanCreateLockTrueOtherEntryLock(t *testing.T) {
	db := test_utils.CreateTestDB()
	models.Migrate(db)

	testEntry := entry.Entry{
		Accession: "test",
	}

	err := db.Create(&testEntry).Error

	assert.Nil(t, err)

	testLock := Lock{
		ID:             1,
		EntryAccession: "test2",
		Category:       Locitax,
		UnlocksAt:      time.Now().Add(time.Minute * 5).Unix(),
	}

	err = db.Create(&testLock).Error

	canCreateLock, err := EntryCanCreateLock(db, testEntry.Accession, Locitax)

	assert.Nil(t, err)
	assert.True(t, canCreateLock)
}

func TestEntryCanCreateLockFalseLockExists(t *testing.T) {
	db := test_utils.CreateTestDB()
	models.Migrate(db)

	testEntry := entry.Entry{
		Accession: "test",
	}

	err := db.Create(&testEntry).Error

	assert.Nil(t, err)

	testLock := Lock{
		ID:             1,
		EntryAccession: testEntry.Accession,
		Category:       Locitax,
		UnlocksAt:      time.Now().Add(time.Minute * 5).Unix(),
	}

	err = db.Create(&testLock).Error

	assert.Nil(t, err)

	canCreateLock, err := EntryCanCreateLock(db, testEntry.Accession, Locitax)

	assert.Nil(t, err)
	assert.False(t, canCreateLock)
}

func TestEntryCanCreateLockTrueOtherCategory(t *testing.T) {
	db := test_utils.CreateTestDB()
	models.Migrate(db)

	testEntry := entry.Entry{
		Accession: "test",
	}

	err := db.Create(&testEntry).Error

	assert.Nil(t, err)

	testLock := Lock{
		ID:             1,
		EntryAccession: testEntry.Accession,
		Category:       Locitax,
		UnlocksAt:      time.Now().Add(time.Minute * 5).Unix(),
	}

	err = db.Create(&testLock).Error

	assert.Nil(t, err)

	canCreateLock, err := EntryCanCreateLock(db, testEntry.Accession, Biosynth)

	assert.Nil(t, err)
	assert.True(t, canCreateLock)
}

func TestEntryCanCreateLockTrueExpired(t *testing.T) {
	db := test_utils.CreateTestDB()
	models.Migrate(db)

	testEntry := entry.Entry{
		Accession: "test",
	}

	err := db.Create(&testEntry).Error

	assert.Nil(t, err)

	testLock := Lock{
		ID:             1,
		EntryAccession: testEntry.Accession,
		Category:       Locitax,
		UnlocksAt:      time.Now().Add(-time.Minute * 5).Unix(),
	}

	err = db.Create(&testLock).Error

	assert.Nil(t, err)

	canCreateLock, err := EntryCanCreateLock(db, testEntry.Accession, Locitax)

	assert.Nil(t, err)
	assert.True(t, canCreateLock)
}

func TestCreateOrGetLock(t *testing.T) {
	db := test_utils.CreateTestDB()
	models.Migrate(db)

	testEntry := entry.Entry{
		Accession: "test",
	}

	err := db.Create(&testEntry).Error

	assert.Nil(t, err)

	testUser := models.User{
		ID: 1,
	}

	err = db.Create(&testUser).Error

	assert.Nil(t, err)

	expectedCategory := Locitax

	lock, err := CreateOrGetLock(db, testEntry.Accession, expectedCategory, testUser)

	assert.Nil(t, err)

	assert.NotNil(t, lock)
	assert.Equal(t, testEntry.Accession, lock.EntryAccession)
	assert.Equal(t, expectedCategory, lock.Category)
}

func TestCreateAlreadyExists(t *testing.T) {
	db := test_utils.CreateTestDB()
	models.Migrate(db)

	testEntry := entry.Entry{
		Accession: "test",
	}

	err := db.Create(&testEntry).Error

	assert.Nil(t, err)

	testUser := models.User{
		ID: 1,
	}

	err = db.Create(&testUser).Error

	assert.Nil(t, err)

	expectedCategory := Locitax
	expectedLock := Lock{
		ID:             0,
		EntryAccession: testEntry.Accession,
		Category:       expectedCategory,
		UnlocksAt:      time.Now().Add(time.Minute * 5).Unix(),
	}

	err = db.Create(&expectedLock).Error

	assert.Nil(t, err)

	lock, err := CreateOrGetLock(db, testEntry.Accession, expectedCategory, testUser)

	assert.Nil(t, err)
	assert.Equal(t, expectedLock.ID, lock.ID)
	assert.Equal(t, expectedLock.Category, lock.Category)
}

func TestCreateAlreadyExistsExpired(t *testing.T) {
	db := test_utils.CreateTestDB()
	models.Migrate(db)

	testEntry := entry.Entry{
		Accession: "test",
	}

	err := db.Create(&testEntry).Error

	assert.Nil(t, err)

	testUser := models.User{
		ID: 1,
	}

	err = db.Create(&testUser).Error

	assert.Nil(t, err)

	expectedCategory := Locitax
	existingLock := Lock{
		ID:             0,
		EntryAccession: testEntry.Accession,
		Category:       expectedCategory,
		UnlocksAt:      time.Now().Add(-time.Minute * 5).Unix(),
	}

	err = db.Create(&existingLock).Error

	assert.Nil(t, err)

	lock, err := CreateOrGetLock(db, testEntry.Accession, expectedCategory, testUser)

	assert.NotNil(t, lock)
	assert.Nil(t, err)
	assert.Equal(t, expectedCategory, lock.Category)
}

func TestDeleteLock(t *testing.T) {
	db := test_utils.CreateTestDB()
	models.Migrate(db)

	testEntry := entry.Entry{
		Accession: "test",
	}

	err := db.Create(&testEntry).Error

	assert.Nil(t, err)

	testUser := models.User{
		ID: 1,
	}

	err = db.Create(&testUser).Error

	assert.Nil(t, err)

	expectedCategory := Locitax
	existingLock := Lock{
		ID:             0,
		EntryAccession: testEntry.Accession,
		Category:       expectedCategory,
		UnlocksAt:      time.Now().Add(time.Minute * 5).Unix(),
		LockOwnerID:    testUser.ID,
		LockOwner:      testUser,
	}

	err = db.Create(&existingLock).Error

	assert.Nil(t, err)

	err = ReleaseLock(db, testEntry.Accession, existingLock.Category, testUser)

	assert.Nil(t, err)

	var actualLock Lock

	tx := db.Find(&actualLock)

	err = tx.Error

	assert.Nil(t, err)
	assert.Equal(t, int64(0), tx.RowsAffected)
}
