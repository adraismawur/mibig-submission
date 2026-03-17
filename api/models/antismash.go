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
	GUID           string            `json:"id" gorm:"primaryKey"`
	EntryAccession string            `json:"entry_accession" gorm:"primaryKey"`
	LocusAccession string            `json:"locus_accession"`
	State          AntismashRunState `json:"state"`
	SubmittedAt    time.Time         `json:"submitted_at"`
}

func init() {
	Models = append(Models, &AntismashRun{})
}
