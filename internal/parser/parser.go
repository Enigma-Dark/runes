package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/Enigma-Dark/runes/internal/types"
)

// ParseReproducerFile parses an Echidna reproducer file
func ParseReproducerFile(filepath string) ([]types.ParsedCall, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filepath, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filepath, err)
	}

	var reproducer types.EchidnaReproducer
	if err := json.Unmarshal(data, &reproducer); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return parseTransactions(reproducer)
}

// parseTransactions converts raw transactions to parsed calls
func parseTransactions(transactions types.EchidnaReproducer) ([]types.ParsedCall, error) {
	var calls []types.ParsedCall

	for _, tx := range transactions {
		// Skip NoCall transactions (they represent time delays)
		if tx.Call.Tag == "NoCall" {
			continue
		}

		if tx.Call.Tag != "SolCall" {
			continue // Skip non-Solidity calls
		}

		call, err := parseCall(tx)
		if err != nil {
			return nil, fmt.Errorf("failed to parse call: %w", err)
		}

		if call != nil {
			calls = append(calls, *call)
		}
	}

	return calls, nil
}

// parseCall converts a transaction to a parsed call
func parseCall(tx types.Transaction) (*types.ParsedCall, error) {
	if len(tx.Call.Contents) < 2 {
		return nil, fmt.Errorf("invalid call contents")
	}

	// Extract function name
	functionName, ok := tx.Call.Contents[0].(string)
	if !ok {
		return nil, fmt.Errorf("function name is not a string")
	}

	// Extract parameters
	paramsInterface, ok := tx.Call.Contents[1].([]interface{})
	if !ok {
		return nil, fmt.Errorf("parameters are not an array")
	}

	var params []types.ParsedParam
	for _, paramInterface := range paramsInterface {
		param, err := parseParameter(paramInterface)
		if err != nil {
			return nil, fmt.Errorf("failed to parse parameter: %w", err)
		}
		params = append(params, param)
	}

	// Parse delay information
	hasDelay, delayValue := parseDelay(tx.Delay)

	return &types.ParsedCall{
		FunctionName: functionName,
		Parameters:   params,
		Dst:          tx.Dst,
		Src:          tx.Src,
		Value:        tx.Value,
		Gas:          tx.Gas,
		GasPrice:     tx.GasPrice,
		HasDelay:     hasDelay,
		DelayValue:   delayValue,
	}, nil
}

// parseDelay extracts delay information from the delay field
func parseDelay(delay []string) (bool, string) {
	if len(delay) < 2 {
		return false, "0"
	}

	// The delay is represented as two hex strings
	// We'll use the first one and convert from hex to decimal seconds
	delayHex := delay[0]
	if delayHex == "0x0000000000000000000000000000000000000000000000000000000000000000" || delayHex == "0x0" {
		return false, "0"
	}

	// Convert hex to decimal
	// Remove 0x prefix if present
	if strings.HasPrefix(delayHex, "0x") {
		delayHex = delayHex[2:]
	}

	// Parse hex string - we'll simplify by taking a reasonable subset
	// For very large hex values, we'll use a simplified approach
	if len(delayHex) > 8 {
		// Take last 8 characters for reasonable delay values
		delayHex = delayHex[len(delayHex)-8:]
	}

	// Convert to integer
	delayInt, err := strconv.ParseInt(delayHex, 16, 64)
	if err != nil || delayInt == 0 {
		return false, "0"
	}

	return true, strconv.FormatInt(delayInt, 10)
}

// parseParameter converts a raw parameter to a parsed parameter
func parseParameter(paramInterface interface{}) (types.ParsedParam, error) {
	paramMap, ok := paramInterface.(map[string]interface{})
	if !ok {
		return types.ParsedParam{}, fmt.Errorf("parameter is not a map")
	}

	tag, ok := paramMap["tag"].(string)
	if !ok {
		return types.ParsedParam{}, fmt.Errorf("parameter tag is not a string")
	}

	contents := paramMap["contents"]

	// Handle special case for AbiBool where contents might be a boolean directly
	if tag == "AbiBool" {
		if boolVal, ok := contents.(bool); ok {
			return types.ParsedParam{
				Type:  "bool",
				Value: strconv.FormatBool(boolVal),
			}, nil
		}
	}

	// For other types, contents should be an array
	contentsArray, ok := contents.([]interface{})
	if !ok {
		return types.ParsedParam{}, fmt.Errorf("parameter contents is not an array for type %s", tag)
	}

	if len(contentsArray) < 1 {
		return types.ParsedParam{}, fmt.Errorf("parameter contents has insufficient elements")
	}

	// Parse based on ABI type
	switch tag {
	case "AbiUInt":
		return parseUintParameter(contentsArray)
	case "AbiInt":
		return parseIntParameter(contentsArray)
	case "AbiAddress":
		return parseAddressParameter(contentsArray)
	case "AbiBool":
		return parseBoolParameter(contentsArray)
	case "AbiBytes":
		return parseBytesParameter(contentsArray)
	case "AbiString":
		return parseStringParameter(contentsArray)
	default:
		return types.ParsedParam{}, fmt.Errorf("unsupported ABI type: %s", tag)
	}
}

// parseUintParameter parses a uint parameter
func parseUintParameter(contents []interface{}) (types.ParsedParam, error) {
	if len(contents) < 2 {
		return types.ParsedParam{}, fmt.Errorf("uint parameter needs at least 2 elements")
	}

	bitSize, ok := contents[0].(float64)
	if !ok {
		return types.ParsedParam{}, fmt.Errorf("uint bit size is not a number")
	}

	value := contents[1]
	var valueStr string

	switch v := value.(type) {
	case string:
		valueStr = v
	case float64:
		valueStr = strconv.FormatFloat(v, 'f', 0, 64)
	default:
		valueStr = fmt.Sprintf("%v", v)
	}

	solType := fmt.Sprintf("uint%d", int(bitSize))

	return types.ParsedParam{
		Type:  solType,
		Value: valueStr,
	}, nil
}

// parseIntParameter parses an int parameter
func parseIntParameter(contents []interface{}) (types.ParsedParam, error) {
	if len(contents) < 2 {
		return types.ParsedParam{}, fmt.Errorf("int parameter needs at least 2 elements")
	}

	bitSize, ok := contents[0].(float64)
	if !ok {
		return types.ParsedParam{}, fmt.Errorf("int bit size is not a number")
	}

	value := contents[1]
	var valueStr string

	switch v := value.(type) {
	case string:
		valueStr = v
	case float64:
		valueStr = strconv.FormatFloat(v, 'f', 0, 64)
	default:
		valueStr = fmt.Sprintf("%v", v)
	}

	solType := fmt.Sprintf("int%d", int(bitSize))

	return types.ParsedParam{
		Type:  solType,
		Value: valueStr,
	}, nil
}

// parseAddressParameter parses an address parameter
func parseAddressParameter(contents []interface{}) (types.ParsedParam, error) {
	if len(contents) < 1 {
		return types.ParsedParam{}, fmt.Errorf("address parameter has no value")
	}

	value, ok := contents[0].(string)
	if !ok {
		return types.ParsedParam{}, fmt.Errorf("address value is not a string")
	}

	// Ensure proper address format
	if !strings.HasPrefix(value, "0x") {
		value = "0x" + value
	}

	return types.ParsedParam{
		Type:  "address",
		Value: value,
	}, nil
}

// parseBoolParameter parses a bool parameter
func parseBoolParameter(contents []interface{}) (types.ParsedParam, error) {
	if len(contents) < 1 {
		return types.ParsedParam{}, fmt.Errorf("bool parameter has no value")
	}

	value, ok := contents[0].(bool)
	if !ok {
		return types.ParsedParam{}, fmt.Errorf("bool value is not a boolean")
	}

	return types.ParsedParam{
		Type:  "bool",
		Value: strconv.FormatBool(value),
	}, nil
}

// parseBytesParameter parses a bytes parameter
func parseBytesParameter(contents []interface{}) (types.ParsedParam, error) {
	if len(contents) < 2 {
		return types.ParsedParam{}, fmt.Errorf("bytes parameter has insufficient elements")
	}

	size, ok := contents[0].(float64)
	if !ok {
		return types.ParsedParam{}, fmt.Errorf("bytes size is not a number")
	}

	value, ok := contents[1].(string)
	if !ok {
		return types.ParsedParam{}, fmt.Errorf("bytes value is not a string")
	}

	var solType string
	if size > 0 {
		solType = fmt.Sprintf("bytes%d", int(size))
	} else {
		solType = "bytes"
	}

	return types.ParsedParam{
		Type:  solType,
		Value: value,
	}, nil
}

// parseStringParameter parses a string parameter
func parseStringParameter(contents []interface{}) (types.ParsedParam, error) {
	if len(contents) < 1 {
		return types.ParsedParam{}, fmt.Errorf("string parameter has no value")
	}

	value, ok := contents[0].(string)
	if !ok {
		return types.ParsedParam{}, fmt.Errorf("string value is not a string")
	}

	return types.ParsedParam{
		Type:  "string",
		Value: fmt.Sprintf(`"%s"`, value), // Wrap in quotes for Solidity
	}, nil
}
