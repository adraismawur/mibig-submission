package entry

import (
	"archive/tar"
	"compress/gzip"
	"github.com/adraismawur/mibig-submission/config"
	"gorm.io/gorm"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const DatabaseURL string = "https://dl.secondarymetabolites.org/mibig/mibig_json_4.0.tar.gz"

func downloadMibigDatabase(url string, dest string) error {

	if _, err := os.Stat(dest); err != nil {
		if !os.IsNotExist(err) {
			slog.Info("[source_db] db data directory detected, skipping download.", "error", err)
			return err
		}
	}

	f, err := os.Create(dest)

	defer f.Close()

	if err != nil {
		slog.Error("[source_db] Error creating temporary file", "error", err)
		return err
	}

	resp, err := http.Get(url)

	if err != nil || resp.StatusCode != http.StatusOK {
		slog.Error("[source_db] Error downloading database", "error", err)
		return err
	}

	_, err = io.Copy(f, resp.Body)

	if err != nil {
		slog.Error("[source_db] Error writing to temporary file", "error", err)
		return err
	}

	return nil
}

func unzipMibigDatabase(zipFile string, dest string) error {
	// create dest folders
	err := os.MkdirAll(dest, os.ModePerm)

	if err != nil {
		slog.Error("[source_db] Error creating output directory", "error", err)
		return err
	}

	zipReader, err := os.Open(zipFile)

	if err != nil {
		slog.Error("[source_db] Error opening tgz file", "error", err)
		return err
	}

	tgz, err := gzip.NewReader(zipReader)

	if err != nil {
		slog.Error("[source_db] Error creating tgz reader", "error", err)
		return err
	}

	defer func(tgz *gzip.Reader) {
		err := tgz.Close()
		if err != nil {
			slog.Error("[source_db] Error closing tgz file", "error", err)
		}
	}(tgz)

	tarReader := tar.NewReader(tgz)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			slog.Error("[source_db] Error reading tar file", "error", err)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			slog.Warn("[source_db] Unexpected directory. Skipping...", "path", header.Name)
			continue
		case tar.TypeReg:
			slog.Info("[source_db] Unzipping file: ", "file", header.Name)

			fileName := strings.Split(header.Name, "/")[1]

			dstPath := filepath.Join(dest, fileName)

			outfile, err := os.Create(dstPath)
			if err != nil {
				slog.Error("[source_db] Error creating file", "error", err)
				return err
			}

			_, err = io.Copy(outfile, tarReader)

			if err != nil {
				slog.Error("[source_db] Failed to copy archive file to destination", "file", header.Name, "destination", dstPath, "error", err)
				return err
			}

			err = outfile.Close()

			if err != nil {
				slog.Error("[source_db] Error closing file", "error", err)
				return err
			}
		}
	}

	return nil
}

func PreloadMibigDatabase(db *gorm.DB) {
	dataPath, err := config.GetConfig("DATA_PATH")

	if err != nil {
		slog.Error("[source_db] Could not get env variable for data path")
		slog.Error("[source_db] Did not preload MIBiG files")
		return
	}

	databaseZipDest := filepath.Join(dataPath, "mibig_db.tar.gz")
	databaseJzonDest := filepath.Join(dataPath, "json")

	err = downloadMibigDatabase(DatabaseURL, databaseZipDest)

	if err != nil {
		slog.Error("[source_db] Error downloading MIBiG data", "error", err)
	}

	err = unzipMibigDatabase(databaseZipDest, databaseJzonDest)

	if err != nil {
		slog.Error("[source_db] Error unzipping MIBiG data", "error", err)
	}

	LoadEntries(db, databaseJzonDest)
}
