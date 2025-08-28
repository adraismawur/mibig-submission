package genbank

import (
	"github.com/goccy/go-json"
	"io"
	"log/slog"
	"net/http"
)

const entrezBase = "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/"
const entrezEndpoint = "esummary.fcgi"

const entrezDbParam = "db"
const entrezRetmodeParam = "retmode"
const entrezIdParam = "id"

const entrezDb = "nucleotide"
const entrezRetmode = "json"

func GetGenbankAccession(accession string) (*EntrezSummaryResult, error) {
	url := entrezBase + entrezEndpoint

	request, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		slog.Error("[validations] [genbank] Error building http request", "error", err)
		return nil, err
	}

	query := request.URL.Query()
	query.Add(entrezDbParam, entrezDb)
	query.Add(entrezRetmodeParam, entrezRetmode)
	query.Add(entrezIdParam, accession)

	request.URL.RawQuery = query.Encode()

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		slog.Error("[validations] [genbank] Error sending http request to entrez", "error", err)
		return nil, err
	}

	data, err := io.ReadAll(response.Body)

	if err != nil {
		slog.Error("[validations] [genbank] Could not read response from entrez", "error", err)
		return nil, err
	}

	slog.Info("[validations] [genbank] Got entrez response", "read bytes", len(data))

	var entry EntrezSummaryResponse

	err = json.Unmarshal(data, &entry)

	if err != nil {
		slog.Error("[validations] [genbank] Could not umarshal entrez response", "error", err)
		return nil, err
	}

	uids := entry.Result["uids"].([]interface{})

	if len(uids) == 0 {
		slog.Info("[validations] [genbank] No entries for accession", "accession", accession)
		return nil, nil
	}

	uid := uids[0].(string)

	resultJson, err := json.Marshal(entry.Result[uid].(map[string]interface{}))

	if err != nil {
		slog.Error("[validations] [genbank] Could not marshal entry result", "error", err)
		return nil, err
	}

	var result EntrezSummaryResult

	err = json.Unmarshal(resultJson, &result)

	if err != nil {
		slog.Error("[validations] [genbank] Could not unmarshal entry result", "error", err)
		return nil, err
	}

	return &result, nil
}
