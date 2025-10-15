package util

import (
	"github.com/adraismawur/mibig-submission/config"
	"github.com/adraismawur/mibig-submission/models"
	"gorm.io/gorm"
	"log/slog"
	"os"
	"os/exec"
	path2 "path"
	"time"
)

type AntismashRunState uint

const (
	Pending AntismashRunState = iota
	Downloading
	Running
	Failed
	Finished
)

type AntismashRun struct {
	GUID      string            `json:"id" gorm:"primaryKey"`
	Accession string            `json:"accession"`
	State     AntismashRunState `json:"state"`
}

func init() {
	models.Models = append(models.Models, &AntismashRun{})
}

func AntismashWorker(db *gorm.DB) {
	for {
		time.Sleep(1 * time.Second)

		request := AntismashRun{}

		result := db.Where("state = ?", Pending).Find(&request)
		err := result.Error

		if err != nil {
			slog.Error("[AntismashWorker] antismash worker error:", "err", err)
			continue
		}

		if result.RowsAffected == 0 {
			slog.Info("[AntismashWorker] Nothing to do")
			continue
		}

		if err != nil {
			slog.Error("[AntismashWorker] antismash worker error:", "err", err)
			continue
		}

		request.State = Downloading
		db.Save(&request)

		gbkPath, err := GetGBK(request.Accession)

		if err != nil {
			slog.Error("[AntismashWorker] Could not get GBK", "Accession", request.Accession, "error", err)

			request.State = Failed
			db.Save(&request)

			continue
		}

		request.State = Running
		db.Save(&request)

		_, err = RunAntismash(*gbkPath, request.Accession)

		if err != nil {
			slog.Error("[AntismashWorker] antismash worker error:", "err", err)

			request.State = Failed
			db.Save(&request)

			continue
		}

		request.State = Finished
		db.Save(&request)
	}
}

// RunAntismash is a helper function that runs antismash on a given GBK path
func RunAntismash(gbkPath string, accession string) (string, error) {
	outputDir := path2.Join(config.Envs["DATA_PATH"], "antismash", accession)

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
		slog.Error("[Antismash] Output:", string(ASOut))
		return "", err
	}

	return string(ASOut), nil
}
