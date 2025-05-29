# Tests

Simple test setup for the Runes CLI tool.

## What's Tested

- **Parser tests** (`internal/parser/parser_test.go`) - Core JSON parsing and ABI type handling
- **Integration test** (`integration_test.go`) - End-to-end workflow from file to generated test

## Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific tests
go test ./internal/parser
go test . -run TestEndToEnd
```

## Test Files

- `internal/parser/testdata/` - Sample reproducer files for testing
- Only essential tests are included - keeps it simple and maintainable 