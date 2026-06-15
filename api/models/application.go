package models

import "time"

type ApplicationState struct {
	FirstTimeSetupDone       bool      `json:"first_time_setup_done"`
	LastGroundTruthDownload  time.Time `json:"last_ground_truth_download"`
	NewSubmissionsEnabled    bool      `json:"new_submissions_enabled"`
	SubmissionReviewsEnabled bool      `json:"submission_reviews_enabled"`
}

type ApplicationAnnouncement struct {
	Date      time.Time `json:"date" gorm:"default:CURRENT_TIMESTAMP"`
	ExpiresAt time.Time `json:"expires_at"`
	Text      string    `json:"text"`
	Severity  string    `json:"severity"`
}

func init() {
	Models = append(Models, ApplicationState{})
	Models = append(Models, ApplicationAnnouncement{})
}
