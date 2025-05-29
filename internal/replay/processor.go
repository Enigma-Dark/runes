package replay

import (
	"fmt"

	"github.com/Enigma-Dark/runes/internal/files"
	"github.com/Enigma-Dark/runes/internal/logger"
	"github.com/Enigma-Dark/runes/internal/parser"
	"github.com/Enigma-Dark/runes/internal/types"
)

// ProcessFiles converts a list of replay files to ReplayGroups with detailed logging
func ProcessFiles(replayFiles []files.FileInfo) ([]types.ReplayGroup, error) {
	var allReplays []types.ReplayGroup
	log := logger.NewProcessorLogger()

	fmt.Printf("Processing %d replay files...\n\n", len(replayFiles))

	for _, file := range replayFiles {
		log.LogFileStart(file.Path)

		calls, err := parser.ParseReproducerFile(file.Path)
		if err != nil {
			log.LogFileFailure(file.Path, err)
			continue
		}

		if len(calls) == 0 {
			log.LogFileFailure(file.Path, fmt.Errorf("no valid calls found"))
			continue
		}

		replayGroup := types.ReplayGroup{
			TestName: "", // Will be set later when output file is known
			Calls:    calls,
			FileName: file.Path,
		}

		// Generate simple display name for logging only
		_, lastFunction := GenerateTestFunctionName(file.Path, "", calls)
		displayName := "test_replay"
		if lastFunction != "" {
			displayName = fmt.Sprintf("test_replay_%s", lastFunction)
		}

		allReplays = append(allReplays, replayGroup)
		log.LogFileSuccess(file.Path, displayName, lastFunction, len(calls))
	}

	// Print summary
	log.LogProcessingSummary()

	if len(allReplays) == 0 {
		return nil, fmt.Errorf("no valid replay files found - all %d files failed to process", log.GetStats().FailureCount)
	}

	return allReplays, nil
}

// GenerateTestFunctionName creates a test function name and returns both the name and the last function
func GenerateTestFunctionName(filePath string, number string, calls []types.ParsedCall) (testName string, lastFunction string) {
	// Create test prefix based on number
	testPrefix := "test_replay_"
	if number != "" {
		testPrefix = fmt.Sprintf("test_replay_%s_", number)
	}

	defaultName := fmt.Sprintf("%sdefault", testPrefix)

	// Find the last function call to use its name
	for i := len(calls) - 1; i >= 0; i-- {
		if calls[i].FunctionName != "" {
			lastFunction = calls[i].FunctionName
			testName = fmt.Sprintf("%s%s", testPrefix, calls[i].FunctionName)
			return
		}
	}

	// No function found, use default
	testName = defaultName
	return
}
