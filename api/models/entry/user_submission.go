package entry

import "github.com/adraismawur/mibig-submission/models"

type SubmissionState string

const (
	DraftSubmission    SubmissionState = "draft"
	PendingReview                      = "pending_review"
	AcceptedSubmission                 = "accepted"
)

type SubmissionType string

const (
	NewSubmission  SubmissionType = "new"
	SubmissionEdit SubmissionType = "edit"
)

type UserSubmission struct {
	ID      uint            `json:"-"`
	EntryID uint            `json:"submission_id"`
	UserID  uint            `json:"user_id"`
	Type    SubmissionType  `json:"type"`
	State   SubmissionState `json:"state"`
}

func init() {
	models.Models = append(models.Models, UserSubmission{})
}
