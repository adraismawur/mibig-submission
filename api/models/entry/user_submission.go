package entry

import "github.com/adraismawur/mibig-submission/models"

type SubmissionState string

const (
	DraftSubmission    SubmissionState = "draft"
	DraftEdit                          = "edit"
	PendingReview                      = "pending_review"
	AcceptedSubmission                 = "accepted"
)

type UserSubmission struct {
	ID      uint            `json:"-"`
	EntryID uint            `json:"submission_id"`
	UserID  uint            `json:"user_id"`
	State   SubmissionState `json:"state"`
}

func init() {
	models.Models = append(models.Models, UserSubmission{})
}
