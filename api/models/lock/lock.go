package lock

import (
	"errors"
	"github.com/adraismawur/mibig-submission/config"
	"github.com/adraismawur/mibig-submission/models"
	"gorm.io/gorm"
	"time"
)

type LockingCategory string

const (
	Locitax         LockingCategory = "locitax"
	Biosynth        LockingCategory = "biosynth"
	Compounds       LockingCategory = "compounds"
	GeneInformation LockingCategory = "gene_information"
	Full            LockingCategory = "full"
)

type Lock struct {
	ID          uint64          `json:"db_id"`
	EntryID     uint64          `json:"db_entry_id" gorm:"uniqueIndex:compositeLockIndex"`
	Category    LockingCategory `json:"category" gorm:"uniqueIndex:compositeLockIndex"`
	UnlocksAt   time.Time       `json:"unlocks_at"`
	LockOwnerID uint64          `json:"db_lock_owner_id"`
	LockOwner   models.User     `json:"lock_owner"`
}

func init() {
	models.Models = append(models.Models, &Lock{})
}

func EntryCanCreateLock(db *gorm.DB, entryId int, category LockingCategory) (bool, error) {
	var lock Lock

	err := db.
		Model(&Lock{}).
		Where("entry_id = ? AND (category = ? OR category = ?)", entryId, category, Full).
		Find(&lock).
		Error

	if err != nil {
		return true, err
	}

	if lock.ID == 0 {
		return true, nil
	}

	if lock.UnlocksAt.Before(time.Now()) {
		return true, nil
	}

	return false, nil
}

func GetEntryLocks(db *gorm.DB, entryId int) (*[]Lock, error) {
	var locks []Lock

	now := time.Now()

	err := db.
		Model(&Lock{}).
		Where("entry_id = ? AND unlocks_at >= ?", entryId, now.UnixMilli()).
		Find(&locks).
		Error

	if err != nil {
		return nil, err
	}

	return &locks, nil
}

func GetEntryLock(db *gorm.DB, entryId int, category LockingCategory) (*Lock, error) {
	var lock Lock

	err := db.
		Model(&Lock{}).
		Where("entry_id = ? AND (category = ? OR category = ?)", entryId, category, Full).
		Find(&lock).
		Error

	if err != nil {
		return nil, err
	}

	return &lock, nil
}

func CreateOrGetLock(db *gorm.DB, entryId int, category LockingCategory, user models.User) (*Lock, error) {
	lockDuration, err := config.GetConfig(config.EnvLockDuration)

	if err != nil {
		return nil, err
	}

	// check if current lock is still active
	activeLock, err := GetEntryLock(db, entryId, category)

	if err != nil {
		return nil, err
	}

	if activeLock.UnlocksAt.After(time.Now()) {
		return activeLock, nil
	}

	// otherwise create a new lock
	canCreateLock, err := EntryCanCreateLock(db, entryId, category)

	if err != nil {
		return nil, err
	}

	if !canCreateLock {
		return nil, errors.New("entry/category is already locked")
	}

	parsedDuration, err := time.ParseDuration(lockDuration)
	newTime := time.Now().Add(parsedDuration)

	lock := Lock{
		ID:          0,
		EntryID:     uint64(entryId),
		Category:    category,
		UnlocksAt:   newTime,
		LockOwnerID: user.ID,
		LockOwner:   user,
	}

	err = db.
		Model(&Lock{}).
		Where("entry_id = ? AND category = ?", entryId, category).
		Assign(Lock{UnlocksAt: lock.UnlocksAt}).
		FirstOrCreate(&lock).
		Error

	if err != nil {
		return nil, err
	}

	return &lock, nil
}

func ReleaseLock(db *gorm.DB, entryId int, category LockingCategory, user models.User) error {
	var existingLock Lock

	err := db.
		Model(&Lock{}).
		Where("entry_id = ? AND category = ?", entryId, category).
		Find(&existingLock).
		Error

	userIsAdmin := false

	for _, role := range user.Roles {
		if role.Role == models.Admin {
			userIsAdmin = true
			break
		}
	}

	hasAccess := existingLock.LockOwnerID == user.ID || userIsAdmin

	if !hasAccess {
		return errors.New("user cannot release lock from entry category")
	}

	if err != nil {
		return err
	}

	if existingLock.ID == 0 {
		return nil
	}

	err = db.
		Delete(existingLock).
		Error

	if err != nil {
		return err
	}

	return nil
}
