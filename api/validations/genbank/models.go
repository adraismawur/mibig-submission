package genbank

type EntrezSummaryHeader struct {
	Type    string `json:"type"`
	Version string `json:"version"`
}

type EntrezSummaryResult struct {
	Uid         string `json:"uid"`
	Term        string `json:"term"`
	Caption     string `json:"caption"`
	Title       string `json:"title"`
	Extra       string `json:"extra"`
	Gi          int    `json:"gi"`
	Createdate  string `json:"createdate"`
	Updatedate  string `json:"updatedate"`
	Taxid       int    `json:"taxid"`
	Slen        int    `json:"slen"`
	Biomol      string `json:"biomol"`
	Moltype     string `json:"moltype"`
	Topology    string `json:"topology"`
	Sourcedb    string `json:"sourcedb"`
	Projectid   string `json:"projectid"`
	Subtype     string `json:"subtype"`
	Subname     string `json:"subname"`
	Geneticcode string `json:"geneticcode"`
	Organism    string `json:"organism"`
	Strain      string `json:"strain"`
	Biosample   string `json:"biosample"`
	Statistics  string `json:"statistics"`
	Properties  struct {
		Na    string `json:"na"`
		Value string `json:"value"`
	} `json:"properties"`
	Oslt struct {
		Indexed bool   `json:"indexed"`
		Value   string `json:"value"`
	} `json:"oslt"`
	Accessionversion string `json:"accessionversion"`
}

type EntrezSummaryResponse struct {
	Header EntrezSummaryHeader    `json:"header"`
	Error  string                 `json:"error,omitempty"`
	Result map[string]interface{} `json:"result"`
}
