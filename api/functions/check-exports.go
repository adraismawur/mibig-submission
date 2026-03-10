package functions

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/adraismawur/mibig-submission/config"
	"github.com/adraismawur/mibig-submission/db"
	"github.com/adraismawur/mibig-submission/models/entry"
	"github.com/goccy/go-json"
	diff2 "github.com/rogpeppe/go-internal/diff"
	"gorm.io/gorm"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type FileMismatchError struct {
	error
}

func listJsonFiles() ([]string, error) {
	dataPath, err := config.GetConfig("DATA_PATH")

	if err != nil {
		slog.Error("[check-exports] Could not get env variable for data path")
		slog.Error("[check-exports] Did not preload MIBiG files")
		return nil, err
	}

	databaseJsonPath := filepath.Join(dataPath, "json")

	dirEntries, err := os.ReadDir(databaseJsonPath)

	var files []string

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}

		fileParts := strings.Split(dirEntry.Name(), "/")

		files = append(files, fileParts[len(fileParts)-1])
	}

	return files, err
}

func checkEntryIdentity(db *gorm.DB, accession string) error {
	dataPath, err := config.GetConfig("DATA_PATH")

	if err != nil {
		slog.Error("[check-exports] Could not get env variable for data path")
		slog.Error("[check-exports] Did not preload MIBiG files")
		return err
	}

	databaseJsonPath := filepath.Join(dataPath, "json")

	diffPath := filepath.Join(databaseJsonPath, "diff")

	err = os.MkdirAll(diffPath, os.ModePerm)

	if err != nil {
		slog.Error("[check-exports] Could not create diff folder")
		return err
	}

	jsonPath := filepath.Join(databaseJsonPath, accession+".json")

	sourceJsonBytes, err := os.ReadFile(jsonPath)

	if err != nil {
		return err
	}

	var sourceJsonBytesIndented bytes.Buffer

	err = json.Indent(&sourceJsonBytesIndented, sourceJsonBytes, "", "  ")

	if err != nil {
		return err
	}

	sourceJsonBytes = sourceJsonBytesIndented.Bytes()

	dbEntry, err := entry.GetEntryExportFromAccession(db, accession)

	if err != nil {
		return err
	}

	dbEntryBytes, err := json.MarshalIndent(dbEntry, "", "  ")

	if err != nil {
		return err
	}

	sourceJsonString := string(sourceJsonBytes)
	dbJsonString := string(dbEntryBytes)

	if sourceJsonString != dbJsonString {
		diff := diff2.Diff("source", sourceJsonBytes, "db", dbEntryBytes)

		diffFilePath := filepath.Join(diffPath, accession+".diff")

		f, err := os.Create(diffFilePath)

		if err != nil {
			return err
		}
		if _, err = f.Write(diff); err != nil {
			_ = f.Close()
			return err
		}

		_ = f.Close()

		return FileMismatchError{
			errors.New("mismatch between source and db entry export"),
		}
	}

	return nil
}

func CheckExports() {
	slog.Info("[check-exports] Starting JSON export check")

	slog.Info("[check-exports] Overriding env variables")
	config.OverrideEnv(config.EnvDbDialect, "sqlite")
	config.OverrideEnv(config.EnvDbPath, "file::memory:?cache=shared")

	slog.Info("[check-exports] Setting up in-memory database")
	memDb, err := db.Connect()

	if err != nil {
		slog.Error("[check-exports] Could not connect to database")
		panic("Could not connect to database")
	}

	slog.Info("[check-exports] Downloading MIBiG database")
	err = entry.DownloadMIBiGdatabase()

	if err != nil {
		slog.Error("[check-exports] Could not download MIBiG database")
		panic("Could not download MIBiG database")
	}

	slog.Info("[check-exports] Loading MIBiG entries into database")

	err = entry.PreloadMibigDatabase(memDb)

	if err != nil {
		slog.Error("[check-exports] Could not load MIBiG database")
		panic("Could not load MIBiG database")
	}

	slog.Info("[check-exports] Listing existing JSON")
	files, err := listJsonFiles()

	if err != nil {
		slog.Error("[check-exports] Could not list JSON files")
		panic("Could not list JSON files")
	}

	numFiles := len(files)
	slog.Info(fmt.Sprintf("[check-exports] %d files found", numFiles))

	matches, mismatches, checkErrors := 0, 0, 0

	for _, file := range files {
		parts := strings.Split(file, ".")

		accession := parts[0]
		err = checkEntryIdentity(memDb, accession)

		if err != nil {
			if errors.As(err, &FileMismatchError{}) {
				mismatches++
				slog.Error("[check-exports] Source and DB entry mismatch", "accession", accession)
				continue
			}
			slog.Error("[check-exports] Could not check entry identity", "error", err)
			checkErrors++
			continue
		}

		matches++
	}

	slog.Info("[check-exports] End of JSON export check")
	slog.Info("[check-exports] Report:")
	slog.Info(fmt.Sprintf("[check-exports] Total number of entries: %d", numFiles))
	slog.Info(fmt.Sprintf("[check-exports] %d matches (%f%%)", matches, float64(matches*100)/float64(numFiles)))
	slog.Info(fmt.Sprintf("[check-exports] %d mismatches (%f%%)", mismatches, float64(mismatches*100)/float64(numFiles)))
	slog.Info(fmt.Sprintf("[check-exports] %d errors (%f%%)", checkErrors, float64(checkErrors*100)/float64(numFiles)))
}
