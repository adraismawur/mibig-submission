package export

import "github.com/lib/pq"

type ReleaseEntry struct {
	Contributors pq.StringArray `json:"contributors" gorm:"type:text[]"`
	Reviewers    pq.StringArray `json:"reviewers" gorm:"type:text[]"`
	Date         string         `json:"date"`
	Comment      string         `json:"comment"`
}

type Release struct {
	Version string         `json:"version"`
	Entries []ReleaseEntry `json:"entries" gorm:"foreignKey:ReleaseID"`
	Date    string         `json:"date"`
}

type Changelog struct {
	Releases []Release `json:"releases" gorm:"foreignKey:ChangelogID"`
}
