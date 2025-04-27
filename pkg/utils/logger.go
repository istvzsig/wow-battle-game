package utils

import (
	"log"
	"os"
	"path/filepath"
)

// NewLogger sets up a logger writing to a file in the "log" directory
func NewLogger(filename string) (*log.Logger, error) {
	logDir := "log"

	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	logPath := filepath.Join(logDir, filename)
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return log.New(logFile, "", log.LstdFlags|log.Lshortfile), nil
}
