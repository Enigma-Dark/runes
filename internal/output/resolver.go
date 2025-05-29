package output

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	DefaultPrefix = "ReplayTest"
	DefaultSuffix = ".t.sol"
	MaxFileCount  = 9999
)

// Config holds output configuration
type Config struct {
	BasePath     string
	ContractName string
	IsMultiple   bool
}

// ResolveOutputPath handles directory output with auto-incrementing names
func ResolveOutputPath(baseName string, isMultiple bool) string {
	// Check if baseName is a directory
	if info, err := os.Stat(baseName); err == nil && info.IsDir() {
		return generateIncrementingPath(baseName)
	}

	// Check if baseName ends with directory separator
	if strings.HasSuffix(baseName, "/") || strings.HasSuffix(baseName, "\\") {
		dirPath := strings.TrimSuffix(strings.TrimSuffix(baseName, "/"), "\\")
		if err := os.MkdirAll(dirPath, 0755); err == nil {
			return generateIncrementingPath(dirPath)
		}
	}

	// Default file handling
	if baseName == "" {
		if isMultiple {
			return "grouped_replays.t.sol"
		}
		return "replay.t.sol"
	}

	return baseName
}

// GenerateContractName creates a contract name from the output file path
func GenerateContractName(outputFile string, fallback string) string {
	baseName := filepath.Base(outputFile)
	baseName = strings.TrimSuffix(baseName, filepath.Ext(baseName))
	baseName = strings.TrimSuffix(baseName, ".t") // Remove .t suffix for .t.sol files

	if strings.Contains(baseName, "ReplayTest_") {
		return strings.ReplaceAll(baseName, "_", "")
	}

	return fallback
}

// generateIncrementingPath finds the next available filename with incrementing number
func generateIncrementingPath(dir string) string {
	for counter := 1; counter <= MaxFileCount; counter++ {
		filename := fmt.Sprintf("%s_%d%s", DefaultPrefix, counter, DefaultSuffix)
		fullPath := filepath.Join(dir, filename)

		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return fullPath
		}
	}

	// Fallback to timestamp-based name
	filename := fmt.Sprintf("%s_%s%s", DefaultPrefix, time.Now().Format("20060102150405"), DefaultSuffix)
	return filepath.Join(dir, filename)
}
