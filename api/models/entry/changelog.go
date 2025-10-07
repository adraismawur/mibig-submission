package entry

import "github.com/adraismawur/mibig-submission/models"
import "github.com/lib/pq"

type ReleaseEntry struct {
	ReleaseID    uint64         `json:"-"`
	Contributors pq.StringArray `json:"contributors" gorm:"type:text[]"`
	Reviewers    pq.StringArray `json:"reviewers" gorm:"type:text[]"`
	Date         string         `json:"date"`
	Comment      string         `json:"comment"`
}

type Release struct {
	ID          uint64         `json:"-"`
	ChangelogID uint64         `json:"-"`
	Version     string         `json:"version"`
	Entries     []ReleaseEntry `json:"entries" gorm:"foreignKey:ReleaseID"`
	Date        string         `json:"date"`
}

type Changelog struct {
	ID       uint64    `json:"-"`
	EntryID  uint64    `json:"-"`
	Releases []Release `json:"releases" gorm:"foreignKey:ChangelogID"`
}

func init() {
	models.Models = append(models.Models, &Changelog{})
	models.Models = append(models.Models, &Release{})
	models.Models = append(models.Models, &ReleaseEntry{})
}
