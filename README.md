<p>
  <img src="assets/logo.png" width="150" />
</p>

# Runes

A CLI tool that converts Echidna fuzzer reproducer files to executable Foundry test files.

## Overview

Echidna generates reproducer files in JSON format when it finds bugs or property violations. This tool parses those reproducer files and generates corresponding Foundry test files that can be executed with `forge test` to reproduce the exact same sequence of function calls.

## Features

- **JSON Parsing**: Parses Echidna reproducer files with complex ABI parameter encoding
- **Type Conversion**: Converts ABI types (AbiUInt, AbiInt, AbiBool, etc.) to proper Solidity types
- **Template Generation**: Uses Go templates to generate clean, readable Foundry test files
- **CLI Interface**: Simple command-line interface with sensible defaults
- **Directory Support**: Automatically finds and uses the oldest .txt file when given a directory
- **Actor Management**: Generates `_setUpActor()` calls for different users
- **Time Delays**: Includes `_delay()` calls for time-based testing
- **Configurable Output**: Customize contract names, test function names, and output paths

## Installation

### Prerequisites

- Go 1.21 or higher

### Install with Go

```bash
go install github.com/Enigma-Dark/runes@latest
```

### Build from source

```bash
git clone https://github.com/Enigma-Dark/runes.git
cd runes
go mod tidy
make build
```

The binary will be available at `build/runes`.

## Usage

### Basic Usage

Convert an Echidna reproducer file to a Foundry test:

```bash
./runes convert reproducer.txt
```

Convert all reproducer files in a directory (automatically selects oldest):

```bash
./runes convert /path/to/reproducers/
```

### Advanced Usage

Customize the output:

```bash
./runes convert reproducer.txt \
  --output MyTest.t.sol \
  --contract MyTestContract \
  --test testBugReproduction
```

### Command-line Options

- `--output, -o`: Output file path (default: `[input-name]_replay.t.sol`)
- `--contract, -c`: Contract name (default: `[input-name]Replay`)
- `--test, -t`: Test function name (default: `testReplay`)
- `--config`: Config file (default: `$HOME/.runes.yaml`)

## Input Format

The tool accepts:
1. **Single file**: A specific .txt reproducer file
2. **Directory**: A folder containing .txt files (automatically selects the oldest)

Echidna reproducer files are in JSON format and contain an array of transaction objects:

```json
[
  {
    "call": {
      "contents": [
        "deposit",
        [
          {"contents": [256, "3625"], "tag": "AbiUInt"},
          {"contents": [8, "0"], "tag": "AbiUInt"}
        ]
      ],
      "tag": "SolCall"
    },
    "dst": "0x7FA9385bE102ac3EAc297483Dd6233D62b3e1496",
    "delay": ["0x0000000000000000000000000000000000000000000000000000000000000000", "0x0000000000000000000000000000000000000000000000000000000000000000"],
    "gas": 1000000000,
    "value": "0x0000000000000000000000000000000000000000000000000000000000000000"
  }
]
```

## Output Format

The tool generates clean, readable Foundry test files in the style of modern property-based testing:

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Test} from "forge-std/Test.sol";

contract TestReplay is Test {
    // Actor addresses (adjust these to match your test setup)
    address constant USER1 = 0x0000000000000000000000000000000000010000;
    address constant USER2 = 0x0000000000000000000000000000000000020000;
    address constant USER3 = 0x0000000000000000000000000000000000030000;
    
    // TODO: Replace with your actual contract instance
    // YourContract Tester;
    
    function setUp() public {
        // TODO: Initialize your contract here
        // Tester = new YourContract();
    }
    
    function test_replay() public {
        _setUpActor(USER1);
        _delay(2);
        Tester.deposit(3625, 0, 1);
        _setUpActor(USER3);
        _delay(840);
        Tester.donateEulerOnlyEVaultToTargetEVault(109, 1, 11);
        Tester.setPrice(1, 0);
        Tester.borrowCV(115792089237316195423570985008687907853269984665640564039457584007913129639935, 0);
    }
    
    function _setUpActor(address actor) internal {
        vm.startPrank(actor);
        // Add any additional actor setup here if needed
    }
    
    function _delay(uint256 timeInSeconds) internal {
        vm.warp(block.timestamp + timeInSeconds);
    }
}
```

## Supported ABI Types

- `AbiUInt` - Unsigned integers (uint8, uint16, uint256, etc.)
- `AbiInt` - Signed integers (int8, int16, int256, etc.)
- `AbiAddress` - Ethereum addresses
- `AbiBool` - Boolean values (supports both array and direct boolean formats)
- `AbiBytes` - Fixed and dynamic byte arrays
- `AbiString` - String values

## Examples

### Example 1: Single File Conversion

```bash
./runes convert reproducer.txt
```

Output: `reproducer_replay.t.sol`

### Example 2: Directory Processing

```bash
./runes convert /path/to/reproducers/
```

Automatically finds the oldest .txt file and converts it.

### Example 3: Custom Output

```bash
./runes convert reproducer.txt \
  --output BugReproduction.t.sol \
  --contract BugReproduction \
  --test test_reproduce_bug
```

## Directory Processing

When you provide a directory path:

1. **Scans for .txt files**: Finds all .txt files in the directory
2. **Selects oldest**: Automatically selects the file with the oldest modification time
3. **Processes**: Converts the selected file to a Foundry test

This is particularly useful when working with Echidna corpus directories that contain multiple reproducer files.

## Development

### Project Structure

```
runes/
├── cmd/                 # CLI commands (cobra)
│   ├── root.go         # Root command setup
│   └── convert.go      # Convert command implementation
├── internal/
│   ├── types/          # Type definitions
│   ├── parser/         # JSON parsing logic
│   └── generator/      # Test file generation
├── main.go             # Entry point
├── go.mod              # Go module definition
└── README.md           # This file
```

### Running Tests

```bash
go test ./...
```

### Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

### Testing

```bash
make test          # Run all tests
make test-verbose  # Run tests with verbose output
```

## TODO & Roadmap

### Near-term improvements
- [ ] Add example reproducer files in `/examples`
- [ ] Support for more ABI types (AbiArray, AbiTuple)
- [ ] Better error messages with line numbers

### Future enhancements
- [ ] Support for multi-file test generation
- [ ] Integration with more fuzzing tools, like [Medusa](https://github.com/crytic/medusa)

### Known limitations
- Complex nested ABI types may need manual adjustment
- Generated tests require manual contract initialization

## License

MIT License - see [LICENSE](LICENSE) file for details. 