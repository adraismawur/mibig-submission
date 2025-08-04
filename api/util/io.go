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

	if file, err = os.Open(path); err != nil {
		slog.Error("[util] Could not open entry file", "path", path)
		return nil
	}

	defer file.Close()

	var bytes []byte
	if bytes, err = io.ReadAll(file); err != nil {
		slog.Error("[util] Could not read entry file", "path", path)
		return nil
	}

	return bytes
}
