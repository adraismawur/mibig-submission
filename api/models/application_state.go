package models

import "time"

func init() {
	Models = append(Models, ApplicationState{})
}

type ApplicationState struct {
	FirstTimeSetupDone       bool      `json:"first_time_setup_done"`
	LastGroundTruthDownload  time.Time `json:"last_ground_truth_download"`
	NewSubmissionsEnabled    bool      `json:"new_submissions_enabled"`
	SubmissionReviewsEnabled bool      `json:"submission_reviews_enabled"`
}
