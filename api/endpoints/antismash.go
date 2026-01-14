package endpoints

import (
	"errors"
	"github.com/adraismawur/mibig-submission/config"
	"github.com/adraismawur/mibig-submission/middleware"
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/models/entry"
	"github.com/adraismawur/mibig-submission/models/entry/biosynthesis"
	"github.com/adraismawur/mibig-submission/util"
	"github.com/beevik/guid"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	path2 "path"
	"strconv"
	"strings"
	"time"
)

type AntismashResultFeature struct {
	Type       string `json:"type"`
	Qualifiers struct {
		Gene             []string `json:"gene,omitempty"`
		Product          []string `json:"product,omitempty"`
		DBCrossReference []string `json:"db_xref,omitempty"`
		Organism         []string `json:"organism,omitempty"`
		GeneKind         []string `json:"gene_kind,omitempty"`
		ProteinID        []string `json:"protein_id,omitempty"`
		LocusTags        []string `json:"locus_tags,omitempty"`
		Type             []string `json:"type,omitempty"`
		Domains          []string `json:"domains,omitempty"`
		Note             []string `json:"note,omitempty"`
	} `json:"qualifiers"`
}

// AntismashResult is a struct that contains information from an AntiSMASH
// result that is relevant for pre-filling an entry
type AntismashResult struct {
	Version string `json:"version"`
	Records []struct {
		ID       string                   `json:"id"`
		Sequence string                   `json:"sequence"`
		Features []AntismashResultFeature `json:"features"`
	} `json:"records"`
}

var moduleTypeMap = map[string]string{
	"type 1 polyketide synthase": "pks-modular",
	"pks":                        "pks-modular",
}

func init() {
	RegisterEndpointGenerator(AntismashEndpoint)
	middleware.AddProtectedRoute(http.MethodPost, "/antismash", models.Admin)
}

// AntismashEndpoint returns the antismash endpoint, used for submitting, checking and stopping antismash runs
func AntismashEndpoint(db *gorm.DB) Endpoint {
	return Endpoint{
		Routes: []Route{
			{
				Method: http.MethodGet,
				Path:   "/antismash",
				Handler: func(c *gin.Context) {
					getAntismashStatus(db, c)
				},
			},
			{
				Method: http.MethodGet,
				Path:   "/antismash/json/:accession",
				Handler: func(c *gin.Context) {
					getAntismashJson(db, c)
				},
			},
			{
				Method: http.MethodGet,
				Path:   "/antismash/list/:bgc_id",
				Handler: func(c *gin.Context) {
					getRecordAntismashAccessions(db, c)
				},
			},
			{
				Method: http.MethodPost,
				Path:   "/antismash",
				Handler: func(c *gin.Context) {
					startAntismashRun(db, c)
				},
			},
		},
	}
}

func getAntismashStatus(db *gorm.DB, c *gin.Context) {
	taskGuid := c.Query("guid")
	taskAccession := c.Query("accession")
	taskBGCID := c.Query("bgc_id")

	if taskGuid == "" && taskAccession == "" && taskBGCID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Must provide a guid, accession or bgc_id"})
		return
	}

	status := models.AntismashRun{}

	query := db

	if taskGuid != "" {
		query = query.Where("guid = ?", taskGuid)
	}

	if taskAccession != "" {
		query = query.Where("accession = ?", taskAccession)
	}

	if taskBGCID != "" {
		query = query.Where("bgc_id = ?", taskBGCID)
	}

	query = query.Find(&status)

	err := query.Error

	if err != nil {
		slog.Error("[Antismash] Get Antismash Status Error", "err", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Antismash Status Error"})
		return
	}

	if query.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "antismash run or runs not found"})
		return
	}

	c.JSON(http.StatusOK, status)
}

func getAntismashJson(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")

	if accession == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no accession given"})
		return
	}

	jsonPath := path2.Join(config.Envs["DATA_PATH"], "antismash", accession, accession+".json")

	c.File(jsonPath)
}

func getRecordAntismashAccessions(db *gorm.DB, c *gin.Context) {
	entryAccession := c.Param("bgc_id")

	if entryAccession == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no entry accession given"})
		return
	}

	var accession []string

	db.Table("antismash_runs").Where("bgc_id = ?", entryAccession).Select("accession").Find(&accession)

	c.JSON(http.StatusOK, accession)
}

func startAntismashRun(db *gorm.DB, c *gin.Context) {
	request := models.AntismashRun{}

	if err := c.ShouldBindJSON(&request); err != nil {
		slog.Error("[Antismash] Could not bind request json")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request.GUID = guid.NewString()

	slog.Info("[Antismash] Starting Antismash Run", "accession", request.Accession)

	db.Create(&request)
}

func AntismashWorker(db *gorm.DB) {
	for {
		time.Sleep(1 * time.Second)

		request := models.AntismashRun{}

		result := db.Where("state = ?", models.Pending).Find(&request)
		err := result.Error

		if err != nil {
			slog.Error("[AntismashWorker] antismash worker error:", "err", err)
			continue
		}

		if result.RowsAffected == 0 {
			//slog.Info("[AntismashWorker] Nothing to do")
			continue
		}

		if err != nil {
			slog.Error("[AntismashWorker] antismash worker error:", "err", err)
			continue
		}

		request.State = models.Downloading
		db.Save(&request)

		gbkPath, err := util.GetGBK(request.Accession)

		if err != nil {
			slog.Error("[AntismashWorker] Could not get GBK", "Accession", request.Accession, "error", err)

			request.State = models.Failed
			db.Save(&request)

			continue
		}

		request.State = models.Running
		db.Save(&request)

		outputDir := path2.Join(config.Envs["DATA_PATH"], "antismash", request.Accession)

		_, err = RunAntismash(*gbkPath, request.Accession, outputDir)

		if err != nil {
			slog.Error("[AntismashWorker] antismash worker error:", "err", err)

			request.State = models.Failed
			db.Save(&request)

			continue
		}

		jsonFile := path2.Join(outputDir, request.Accession+".json")

		antismashOutput, err := ReadAntismashJson(jsonFile)

		if err != nil {
			slog.Error("[AntismashWorker] Failed to read antismash output", "err", err)
			request.State = models.Failed
			db.Save(&request)
		}

		activeEntry, err := entry.GetEntryFromAccession(db, request.BGCID)
		if err != nil {
			slog.Error("[AntismashWorker] Failed to get entry from accession", "err", err)
			request.State = models.Failed
			db.Save(&request)
		}

		PrefillAntismash(activeEntry, antismashOutput)

		db.Save(activeEntry)

		request.State = models.Finished
		db.Save(&request)
	}
}

// RunAntismash is a helper function that runs antismash on a given GBK path
func RunAntismash(gbkPath string, accession string, outputDir string) (string, error) {

	err := os.MkdirAll(outputDir, 0755)

	if err != nil {
		slog.Error("[Antismash] Could not create output directory", "path", outputDir)
		return "", err
	}

	htmlPath := path2.Join(outputDir, "index.html")

	if _, err := os.Stat(htmlPath); err == nil {
		slog.Info("[util] [genbank] File already exists", "path", htmlPath)
		return "Output already exists", nil
	}

	ASCmd := exec.Command("antismash", gbkPath, "--output-dir", outputDir)

	slog.Info("[Antismash] Running Antismash Command", "gbk", gbkPath, "cmd", ASCmd)

	ASOut, err := ASCmd.Output()

	if err != nil {
		slog.Error("[Antismash] Error executing antismash", "gbkPath", gbkPath, "error", err)
		slog.Error("[Antismash] Output:", "output", string(ASOut))
		return "", err
	}

	return string(ASOut), nil
}

// ReadAntismashJson returns a reduced set of antismash result json for use in filling entry information
func ReadAntismashJson(asJsonPath string) (*AntismashResult, error) {
	var antismashResult AntismashResult

	data := util.ReadFile(asJsonPath)

	err := json.Unmarshal(data, &antismashResult)

	if err != nil {
		slog.Error("[Antismash] Could not unmarshal antismash result", "error", err)
		return nil, err
	}

	return &antismashResult, nil
}

func PrefillAntismash(entry *entry.Entry, antismashResult *AntismashResult) {
	var err error
	// moduleName is an integer here for easy counting, but will be converted to a string later
	moduleName := 1

	for _, record := range antismashResult.Records {
		for _, feature := range record.Features {
			if feature.Type == "source" {
				err = PrefillAntismashTaxonomy(entry, &feature)
				if err != nil {
					slog.Error("[Antismash] Could not parse taxonomy", "error", err)
				}
			}

			if feature.Type == "region" {
				err = PrefillAntismashBiosynthesisClass(entry, &feature)
			}

			if feature.Type == "aSModule" {
				//if len(feature.Qualifiers.GeneKind) == 0 || feature.Qualifiers.GeneKind[0] != "biosynthetic" {
				//	continue
				//}

				module, err := GenerateAntismashBiosynthesisModule(&feature, moduleName)

				if err != nil {
					slog.Error("[Antismash] Could not parse biosynthesis", "error", err)
					continue
				}

				val, exists := moduleTypeMap[module.Type]

				// case where we catch a CDS that is marked as biosynthetic, but we don't cover the type
				if !exists {
					continue
				}

				// convert the type into the type as it is in MIBiG
				module.Type = val

				entry.Biosynthesis.Modules = append(entry.Biosynthesis.Modules, *module)

				moduleName++
			}
		}
	}
}

func PrefillAntismashTaxonomy(entry *entry.Entry, feature *AntismashResultFeature) error {
	if feature.Type != "source" {
		return errors.New("antismash feature is not a source feature")
	}

	entry.Taxonomy.Name = feature.Qualifiers.Organism[0]

	// taxid is more complicated
	taxIDParts := feature.Qualifiers.DBCrossReference[0]
	taxIDString := strings.Split(taxIDParts, ":")[1]
	taxID, err := strconv.ParseInt(taxIDString, 10, 64)

	if err != nil {
		return err
	}

	entry.Taxonomy.TaxID = uint64(taxID)

	return nil
}

func PrefillAntismashBiosynthesisClass(entry *entry.Entry, feature *AntismashResultFeature) error {
	if feature.Type != "region" {
		return errors.New("antismash feature is not a region feature")
	}

	var class biosynthesis.BiosyntheticClass

	switch feature.Qualifiers.Product[0] {
	case "T1PKS":
		class.Class = "PKS"
		class.Subclass = "Type I"

	case "T2PKS":
		class.Class = "PKS"
		class.Subclass = "Type II"
	}
	entry.Biosynthesis.Classes = append(entry.Biosynthesis.Classes, class)

	return nil
}

func GenerateAntismashBiosynthesisModule(feature *AntismashResultFeature, moduleName int) (*biosynthesis.BiosyntheticModule, error) {
	//if feature.Type != "CDS" {
	//	return nil, errors.New("antismash feature is not a CDS feature")
	//}
	//
	//var module biosynthesis.BiosyntheticModule
	//
	//module.Genes = append(module.Genes, feature.Qualifiers.ProteinID[0])
	//module.Active = true
	//module.Type = feature.Qualifiers.Product[0]
	//module.Name = strconv.Itoa(moduleName)
	//
	//var err error
	//
	//module.ATDomain, err = GeneratePKSATDomain(feature, module)
	//
	//if err != nil {
	//	slog.Error("[Antismash] Could not generate PKS AT domain", "error", err)
	//	return nil, err
	//}
	//
	//return &module, nil

	if feature.Type != "aSModule" {
		return nil, errors.New("antismash feature is not a aSModule feature")
	}

	var module biosynthesis.BiosyntheticModule

	module.Genes = append(module.Genes, feature.Qualifiers.LocusTags[0])
	module.Active = true
	module.Type = feature.Qualifiers.Type[0]
	module.Name = strconv.Itoa(moduleName)

	var err error

	module.ATDomain, err = GeneratePKSATDomain(feature)

	if err != nil {
		slog.Error("[Antismash] Could not generate PKS AT domain", "error", err)
		return nil, err
	}

	for _, domainString := range feature.Qualifiers.Domains {

		parts := strings.Split(domainString, "_")
		domain := strings.Split(parts[2], ".")[0]

		switch domain {
		case "ACP":
			carrier, err := GenerateCarrier(module.Genes[0], "ACP")

			if err != nil {
				slog.Error("[Antismash] Could not generate PKS carrier", "error", err)
			}

			module.Carriers = append(module.Carriers, *carrier)
		case "PKS":
			domainType := strings.Split(parts[3], ".")[0]
			modificationDomain, err := GenerateModificationDomain(module.Genes[0], domainType)

			if err != nil {
				slog.Info("[Antismash] Could not generate modification domain", "error", err)
				continue
			}

			module.ModificationDomains = append(module.ModificationDomains, *modificationDomain)
		}
	}

	return &module, nil
}

func GeneratePKSATDomain(feature *AntismashResultFeature) (*biosynthesis.ATModuleDomain, error) {
	atDomain := biosynthesis.ATModuleDomain{}

	atDomain.Gene = feature.Qualifiers.LocusTags[0]
	atDomain.Location.From = -1
	atDomain.Location.To = -1

	return &atDomain, nil
}

func GenerateCarrier(gene string, subType string) (*biosynthesis.CarrierModuleDomain, error) {
	domain := biosynthesis.CarrierModuleDomain{}

	domain.Gene = gene
	domain.Subtype = subType
	domain.Location.From = -1
	domain.Location.To = -1

	return &domain, nil
}

var modificationDomainTypeMap = map[string]string{
	"DH": "dehydratase",
	"KR": "ketoreductase",
}

func GenerateModificationDomain(gene string, domainType string) (*biosynthesis.ModificationModuleDomain, error) {
	domain := biosynthesis.ModificationModuleDomain{}

	if mibigDomainType, ok := modificationDomainTypeMap[domainType]; !ok {
		return nil, errors.New("modification domain type not supported")
	} else {
		domain.DomainType = mibigDomainType
	}

	domain.Gene = gene
	domain.Location.From = -1
	domain.Location.To = -1

	return &domain, nil
}
