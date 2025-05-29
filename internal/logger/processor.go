package logger

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ProcessingStats holds statistics about replay processing
type ProcessingStats struct {
	TotalFiles   int
	SuccessCount int
	FailureCount int
	FailedFiles  []FailedFile
	SuccessTests []SuccessTest
}

// FailedFile represents a file that failed to process
type FailedFile struct {
	FileName string
	Error    string
}

// SuccessTest represents a successfully processed test
type SuccessTest struct {
	FileName     string
	TestName     string
	LastFunction string
	CallCount    int
}

// ProcessorLogger handles logging for replay processing
type ProcessorLogger struct {
	stats ProcessingStats
}

// NewProcessorLogger creates a new processor logger
func NewProcessorLogger() *ProcessorLogger {
	return &ProcessorLogger{
		stats: ProcessingStats{
			FailedFiles:  make([]FailedFile, 0),
			SuccessTests: make([]SuccessTest, 0),
		},
	}
}

// LogFileStart logs the start of processing a file
func (l *ProcessorLogger) LogFileStart(filePath string) {
	l.stats.TotalFiles++
	fileName := filepath.Base(filePath)
	fmt.Printf("- Processing: %s\n", fileName)
}

// LogFileSuccess logs successful processing of a file
func (l *ProcessorLogger) LogFileSuccess(filePath, testName, lastFunction string, callCount int) {
	l.stats.SuccessCount++
	fileName := filepath.Base(filePath)

	l.stats.SuccessTests = append(l.stats.SuccessTests, SuccessTest{
		FileName:     fileName,
		TestName:     testName,
		LastFunction: lastFunction,
		CallCount:    callCount,
	})

	if lastFunction != "" {
		fmt.Printf("  ✓ %s -> %s (last: %s, %d calls)\n",
			fileName, testName, lastFunction, callCount)
	} else {
		fmt.Printf("  ✓ %s -> %s (%d calls)\n",
			fileName, testName, callCount)
	}
}

// LogFileFailure logs failed processing of a file
func (l *ProcessorLogger) LogFileFailure(filePath string, err error) {
	l.stats.FailureCount++
	fileName := filepath.Base(filePath)

	l.stats.FailedFiles = append(l.stats.FailedFiles, FailedFile{
		FileName: fileName,
		Error:    err.Error(),
	})

	fmt.Printf("  ✗ %s - %v\n", fileName, err)
}

// LogProcessingSummary logs a summary of all processing
func (l *ProcessorLogger) LogProcessingSummary() {
	fmt.Println("\n" + strings.Repeat("-", 50))
	fmt.Println("PROCESSING SUMMARY")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Printf("Total files: %d | Success: %d | Failed: %d\n",
		l.stats.TotalFiles, l.stats.SuccessCount, l.stats.FailureCount)

	if l.stats.SuccessCount > 0 {
		fmt.Println("\nGenerated tests:")
		for i, test := range l.stats.SuccessTests {
			fmt.Printf("  %d. %s (last call: %s)\n",
				i+1, test.TestName, test.LastFunction)
		}
	}

	if l.stats.FailureCount > 0 {
		fmt.Println("\nFailed files:")
		for i, failed := range l.stats.FailedFiles {
			fmt.Printf("  %d. %s: %s\n", i+1, failed.FileName, failed.Error)
		}
	}

	// Success rate
	if l.stats.TotalFiles > 0 {
		successRate := float64(l.stats.SuccessCount) / float64(l.stats.TotalFiles) * 100
		fmt.Println()
		fmt.Println(strings.Repeat("-", 50))
		fmt.Printf("Success rate: %.1f%%\n", successRate)
	}

	fmt.Println(strings.Repeat("-", 50))
}

// GetStats returns the current processing statistics
func (l *ProcessorLogger) GetStats() ProcessingStats {
	return l.stats
}
