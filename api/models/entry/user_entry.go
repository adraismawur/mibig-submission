package entry

import "github.com/adraismawur/mibig-submission/models"

type UserEntry struct {
	ID      uint   `json:-`
	EntryID uint   `json:-`
	UserID  string `json:-`
}

func init() {
	models.Models = append(models.Models, UserEntry{})
}
