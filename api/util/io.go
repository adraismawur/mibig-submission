package util

import (
	"io"
	"log/slog"
	"os"
)

// ReadFile reads a file into a byte array
func ReadFile(path string) []byte {
	var err error
	var file *os.File

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			slog.Error("[sub] Could not close file %s", "path", path)
		}
	}(file)

	if file, err = os.Open(path); err != nil {
		slog.Error("[sub] Could not open submission file", "path", path)
		return nil
	}

	var bytes []byte
	if bytes, err = io.ReadAll(file); err != nil {
		slog.Error("[sub] Could not read submission file", "path", path)
		return nil
	}

	return bytes
}
