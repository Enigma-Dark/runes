package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Enigma-Dark/runes/internal/files"
	"github.com/Enigma-Dark/runes/internal/generator"
	"github.com/Enigma-Dark/runes/internal/parser"
	"github.com/Enigma-Dark/runes/internal/replay"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEndToEnd tests the complete workflow from reproducer file to generated test
func TestEndToEnd(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a simple reproducer file
	reproducerContent := `[
		{
			"call": {
				"contents": ["deposit", [
					{"contents": [256, "1000"], "tag": "AbiUInt"}
				]],
				"tag": "SolCall"
			},
			"dst": "0x7FA9385bE102ac3EAc297483Dd6233D62b3e1496",
			"src": "0x0000000000000000000000000000000000010000",
			"delay": ["0x0", "0x0"],
			"gas": 1000000,
			"gasPrice": "0x0",
			"value": "0x0"
		}
	]`

	reproducerFile := filepath.Join(tmpDir, "test_reproducer.txt")
	err := os.WriteFile(reproducerFile, []byte(reproducerContent), 0644)
	require.NoError(t, err)

	outputFile := filepath.Join(tmpDir, "TestReplay.t.sol")

	// Step 1: Discover files
	discoveredFiles, err := files.DiscoverReplayFiles(reproducerFile)
	require.NoError(t, err)
	assert.Len(t, discoveredFiles, 1)

	// Step 2: Parse reproducer file
	calls, err := parser.ParseReproducerFile(reproducerFile)
	require.NoError(t, err)
	assert.Len(t, calls, 1)
	assert.Equal(t, "deposit", calls[0].FunctionName)

	// Step 3: Process into replay groups
	replayGroups, err := replay.ProcessFiles(discoveredFiles)
	require.NoError(t, err)
	assert.Len(t, replayGroups, 1)

	// Step 4: Generate test file
	config := generator.GenerateConfig{
		ContractName: "TestReplay",
		OutputFile:   outputFile,
		ReplayGroups: replayGroups,
		Template:     "enigmadark",
	}

	err = generator.GenerateFoundryTest(config)
	require.NoError(t, err)

	// Step 5: Verify generated file exists and has expected content
	assert.FileExists(t, outputFile)

	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	contentStr := string(content)
	assert.Contains(t, contentStr, "contract TestReplay")
	assert.Contains(t, contentStr, "deposit(1000)")
}

// TestErrorHandling tests basic error scenarios
func TestErrorHandling(t *testing.T) {
	// Test with non-existent file
	_, err := files.DiscoverReplayFiles("/does/not/exist")
	assert.Error(t, err)

	// Test with invalid JSON
	tmpDir := t.TempDir()
	invalidFile := filepath.Join(tmpDir, "invalid.txt")
	err = os.WriteFile(invalidFile, []byte(`{"invalid": "json"}`), 0644)
	require.NoError(t, err)

	_, err = parser.ParseReproducerFile(invalidFile)
	assert.Error(t, err)
}
