package models

import "time"

func init() {
	Models = append(Models, DatabaseMeta{})
}

type DatabaseMeta struct {
	FirstTimeSetupDone      bool
	LastGroundTruthDownload time.Time
}
