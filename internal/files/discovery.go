package files

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// FileInfo holds file information for processing
type FileInfo struct {
	Path    string
	ModTime time.Time
}

// DiscoverReplayFiles resolves input path to a list of files to process
func DiscoverReplayFiles(inputPath string) ([]FileInfo, error) {
	info, err := os.Stat(inputPath)
	if err != nil {
		return nil, fmt.Errorf("path does not exist: %s", inputPath)
	}

	if !info.IsDir() {
		return []FileInfo{{
			Path:    inputPath,
			ModTime: info.ModTime(),
		}}, nil
	}

	return findNewestTxtGroup(inputPath)
}

// findNewestTxtGroup finds .txt files grouped by creation time and returns the newest group
func findNewestTxtGroup(dirPath string) ([]FileInfo, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	txtFiles := collectTxtFiles(dirPath, files)
	if len(txtFiles) == 0 {
		return nil, fmt.Errorf("no .txt files found in directory: %s", dirPath)
	}

	groups := groupByCreationTime(txtFiles)
	newestGroup := selectNewestGroup(groups)

	// Sort files within the group by name for consistent ordering
	sort.Slice(newestGroup, func(i, j int) bool {
		return filepath.Base(newestGroup[i].Path) < filepath.Base(newestGroup[j].Path)
	})

	fmt.Printf("\nFound %d groups of .txt files, selecting newest group with %d files (created at %s)\n",
		len(groups), len(newestGroup), newestGroup[0].ModTime.Truncate(time.Minute).Format("2006-01-02 15:04"))

	return newestGroup, nil
}

// collectTxtFiles gathers all .txt files from the directory
func collectTxtFiles(dirPath string, files []os.DirEntry) []FileInfo {
	var txtFiles []FileInfo

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".txt") {
			fullPath := filepath.Join(dirPath, file.Name())
			if info, err := os.Stat(fullPath); err == nil {
				txtFiles = append(txtFiles, FileInfo{
					Path:    fullPath,
					ModTime: info.ModTime(),
				})
			}
		}
	}

	return txtFiles
}

// groupByCreationTime groups files by creation time (truncated to the minute)
func groupByCreationTime(files []FileInfo) map[string][]FileInfo {
	groups := make(map[string][]FileInfo)

	for _, file := range files {
		timeKey := file.ModTime.Truncate(time.Minute).Format("2006-01-02 15:04")
		groups[timeKey] = append(groups[timeKey], file)
	}

	return groups
}

// selectNewestGroup finds the group with the most recent creation time
func selectNewestGroup(groups map[string][]FileInfo) []FileInfo {
	var newestTime time.Time
	var newestGroup []FileInfo

	for _, group := range groups {
		groupTime := group[0].ModTime.Truncate(time.Minute)
		if newestTime.IsZero() || groupTime.After(newestTime) {
			newestTime = groupTime
			newestGroup = group
		}
	}

	return newestGroup
}
