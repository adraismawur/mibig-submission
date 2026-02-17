package entry

// FinalDetails describes whether an entry contains all genes responsible for
// production of components (completeness) and whether it is under embargo
type FinalDetails struct {
	Completeness Completeness `json:"completeness"`
	Embargo      bool         `json:"embargo"`
}

// FinalDetailsRequest adds a comments section that is not available in the mibig schema
// directly, but can be used to add information to the changelog
type FinalDetailsRequest struct {
	FinalDetails
	Comment string `json:"comments"`
}
