package models

import "time"

type AntismashRunState uint

const (
	Pending AntismashRunState = iota
	Downloading
	Running
	Failed
	Finished
)

type AntismashRun struct {
	GUID        string            `json:"id" gorm:"primaryKey"`
	BGCID       string            `json:"bgc_id" gorm:"primaryKey"`
	Accession   string            `json:"accession"`
	State       AntismashRunState `json:"state"`
	SubmittedAt time.Time         `json:"submitted_at"`
}

func init() {
	Models = append(Models, &AntismashRun{})
}
