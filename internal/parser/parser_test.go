package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseReproducerFile(t *testing.T) {
	// Test with valid reproducer file
	calls, err := ParseReproducerFile("testdata/valid_reproducer.json")
	require.NoError(t, err)
	assert.Len(t, calls, 2)

	// First call should be deposit
	assert.Equal(t, "deposit", calls[0].FunctionName)
	assert.Len(t, calls[0].Parameters, 2)

	// Second call should be withdraw
	assert.Equal(t, "withdraw", calls[1].FunctionName)
	assert.True(t, calls[1].HasDelay) // 0x1E should be parsed as delay
}

func TestParseReproducerFile_Errors(t *testing.T) {
	// Non-existent file
	_, err := ParseReproducerFile("nonexistent.json")
	assert.Error(t, err)

	// Invalid JSON
	_, err = ParseReproducerFile("testdata/invalid_reproducer.json")
	assert.Error(t, err)

	// Empty file should work but return no calls
	calls, err := ParseReproducerFile("testdata/empty_reproducer.json")
	require.NoError(t, err)
	assert.Len(t, calls, 0)
}

func TestParseParameter_BasicTypes(t *testing.T) {
	// Test uint
	param, err := parseParameter(map[string]interface{}{
		"tag":      "AbiUInt",
		"contents": []interface{}{256.0, "1000"},
	})
	require.NoError(t, err)
	assert.Equal(t, "uint256", param.Type)
	assert.Equal(t, "1000", param.Value)

	// Test address
	param, err = parseParameter(map[string]interface{}{
		"tag":      "AbiAddress",
		"contents": []interface{}{"0x1234567890123456789012345678901234567890"},
	})
	require.NoError(t, err)
	assert.Equal(t, "address", param.Type)

	// Test bool
	param, err = parseParameter(map[string]interface{}{
		"tag":      "AbiBool",
		"contents": true,
	})
	require.NoError(t, err)
	assert.Equal(t, "bool", param.Type)
	assert.Equal(t, "true", param.Value)
}
