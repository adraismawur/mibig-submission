package entry

import (
	"github.com/adraismawur/mibig-submission/config"
	"gorm.io/gorm"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

const DatabaseURL string = "https://dl.secondarymetabolites.org/mibig/mibig_json_4.0.tar.gz"

func downloadMibigDatabase(url string, dest string) {

	if _, err := os.Stat(dest); err != nil {
		if !os.IsNotExist(err) {
			slog.Info("[source_db] db data directory detected, skipping download.")
			return
		}
	}

	f, err := os.Create(dest)

	defer f.Close()

	if err != nil {
		slog.Error("[source_db] Error creating temporary file")
		return
	}

	resp, err := http.Get(url)

	if err != nil || resp.StatusCode != http.StatusOK {
		slog.Error("[source_db] Error downloading database")
		return
	}

	_, err = io.Copy(f, resp.Body)

	if err != nil {
		slog.Error("[source_db] Error writing to temporary file")
		return
	}
}

func PreloadMibigDatabase(db *gorm.DB) {
	dataPath, err := config.GetConfig("DATA_PATH")

	if err != nil {
		slog.Error("[source_db] Could not get env variable for data path")
		slog.Error("[source_db] Did not preload mibig")
		return
	}

	databaseZipDest := filepath.Join(dataPath, "mibig_db.tar.gz")
	databaseJzonDest := filepath.Join(dataPath, "json")

	downloadMibigDatabase(DatabaseURL, databaseZipDest)
	// TODO: unzip it to data/json/*.json
	LoadEntries(db, databaseJzonDest)
}
