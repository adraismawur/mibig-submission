package consts

type Quality string

const (
	Questionable Quality = "questionable"
	Medium       Quality = "medium"
	High         Quality = "high"
)

type Status string

const (
	Pending Status = "pending"
	Active  Status = "active"
	Retired Status = "retired"
)

type Completeness string

const (
	Unknown    Completeness = "unknown"
	Complete   Completeness = "complete"
	Incomplete Completeness = "incomplete"
)
