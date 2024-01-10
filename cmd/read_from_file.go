package main

import (
	"io"
	"os"
)

func readFromFile(filePath string) string {
	f, err := os.Open(filePath)
	if err != nil {
		logger.Fatalf("unable to open file: %v", err)
		return ""
	}

	content, err := io.ReadAll(f)
	if err != nil {
		logger.Fatalf("unable to read file: %v", err)
		return ""
	}
	return string(content)
}
