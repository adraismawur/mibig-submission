package util

import (
	"errors"
	"github.com/adraismawur/mibig-submission/config"
	"github.com/goccy/go-json"
	"io"
	"log/slog"
	"net/http"
	"os"
	path2 "path"
	"path/filepath"
)

const entrezBase = "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/"
const entrezDb = "nucleotide"

const entrezIdParam = "id"
const entrezDbParam = "db"
const entrezRetmodeParam = "retmode"
const entrezRettypeParam = "rettype"

const entrezSummaryEndpoint = "esummary.fcgi"
const entrezSummaryRetmode = "json"

const entrezFetchEndpoint = "efetch.fcgi"
const entrezFetchRetmode = "txt"
const entrezFetchRettype = "gb"

const GBKSubPath = "/gbk"

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
	CreateDate  string `json:"createdate"`
	UpdateDate  string `json:"update"`
	TaxID       int    `json:"taxied"`
	SLen        int    `json:"slen"`
	BioMol      string `json:"biol"`
	MolType     string `json:"moltype"`
	Topology    string `json:"topology"`
	SourceDB    string `json:"sourcedb"`
	ProjectId   string `json:"projectid"`
	Subtype     string `json:"subtype"`
	SubName     string `json:"subname"`
	GeneticCode string `json:"geneticcode"`
	Organism    string `json:"organism"`
	Strain      string `json:"strain"`
	BioSample   string `json:"biosample"`
	Statistics  string `json:"statistics"`
	Properties  struct {
		Na    string `json:"na"`
		Value string `json:"value"`
	} `json:"properties"`
	Oslt struct {
		Indexed bool   `json:"indexed"`
		Value   string `json:"value"`
	} `json:"oslt"`
	AccessionVersion string `json:"accessionversion"`
}

type EntrezSummaryResponse struct {
	Header EntrezSummaryHeader    `json:"header"`
	Error  string                 `json:"error,omitempty"`
	Result map[string]interface{} `json:"result"`
}

func GetGenbankAccessionSummary(accession string) (*EntrezSummaryResult, error) {
	url := entrezBase + entrezSummaryEndpoint

	request, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		slog.Error("[util] [genbank] Error building http request", "error", err)
		return nil, err
	}

	query := request.URL.Query()
	query.Add(entrezDbParam, entrezDb)
	query.Add(entrezRetmodeParam, entrezSummaryRetmode)
	query.Add(entrezIdParam, accession)

	request.URL.RawQuery = query.Encode()

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		slog.Error("[util] [genbank] Error sending http request to entrez", "error", err)
		return nil, err
	}

	data, err := io.ReadAll(response.Body)

	if err != nil {
		slog.Error("[util] [genbank] Could not read response from entrez", "error", err)
		return nil, err
	}

	slog.Info("[util] [genbank] Got entrez response", "read bytes", len(data))

	var entry EntrezSummaryResponse

	err = json.Unmarshal(data, &entry)

	if err != nil {
		slog.Error("[util] [genbank] Could not umarshal entrez response", "error", err)
		return nil, err
	}

	uids := entry.Result["uids"].([]interface{})

	if len(uids) == 0 {
		slog.Info("[util] [genbank] No entries for accession", "accession", accession)
		return nil, nil
	}

	uid := uids[0].(string)

	resultJson, err := json.Marshal(entry.Result[uid].(map[string]interface{}))

	if err != nil {
		slog.Error("[util] [genbank] Could not marshal entry result", "error", err)
		return nil, err
	}

	var result EntrezSummaryResult

	err = json.Unmarshal(resultJson, &result)

	if err != nil {
		slog.Error("[util] [genbank] Could not unmarshal entry result", "error", err)
		return nil, err
	}

	return &result, nil
}

// GetGBK downloads a gbk from entrez and returns the path it was saved at as a string
func GetGBK(accession string) (*string, error) {
	basePath := path2.Join(config.Envs["DATA_PATH"], GBKSubPath)

	err := os.MkdirAll(basePath, 0755)

	if err != nil {
		slog.Error("[util] [genbank] Could not create output directory", "path", basePath)
		return nil, err
	}

	path := filepath.Join(basePath, accession+".gbk")

	if _, err := os.Stat(path); err == nil {
		slog.Info("[util] [genbank] File already exists", "path", path)
		return &path, nil
	}

	url := entrezBase + entrezFetchEndpoint

	request, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		slog.Error("[util] [genbank] Error building http request", "error", err)
		return nil, err
	}

	query := request.URL.Query()
	query.Add(entrezDbParam, entrezDb)
	query.Add(entrezIdParam, accession)
	query.Add(entrezRetmodeParam, entrezFetchRetmode)
	query.Add(entrezRettypeParam, entrezFetchRettype)

	request.URL.RawQuery = query.Encode()

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		slog.Error("[util] [genbank] Error sending http request to entrez", "error", err)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusBadRequest {
			slog.Error("[util] [genbank] Error retrieving genbank from entrez.")
			slog.Error("[util] [genbank] ", "URL", request.URL)
			return nil, errors.New(response.Status)
		}
		return nil, errors.New(response.Status)
	}

	data, err := io.ReadAll(response.Body)

	if err != nil {
		slog.Error("[util] [genbank] Could not read response from entrez", "error", err)
		return nil, err
	}

	slog.Info("[util] [genbank] Got entrez response", "read bytes", len(data))

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)

	if err != nil {
		slog.Error("[util] [genbank] Error opening file for writing", "error", err)
		return nil, err
	}

	defer f.Close()

	_, err = f.Write(data)

	if err != nil {
		slog.Error("[util] [genbank] Error writing to file", "error", err)
		return nil, err
	}

	return &path, nil
}
