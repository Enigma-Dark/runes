package types

// EchidnaReproducer represents the root structure of an Echidna reproducer file
type EchidnaReproducer []Transaction

// Transaction represents a single transaction in the reproducer
type Transaction struct {
	Call     Call     `json:"call"`
	Delay    []string `json:"delay"`
	Dst      string   `json:"dst"`
	Gas      int64    `json:"gas"`
	GasPrice string   `json:"gasprice"`
	Src      string   `json:"src"`
	Value    string   `json:"value"`
}

// Call represents the function call within a transaction
type Call struct {
	Contents []interface{} `json:"contents,omitempty"`
	Tag      string        `json:"tag"`
}

// AbiParam represents a parameter with ABI type information
type AbiParam struct {
	Contents []interface{} `json:"contents"`
	Tag      string        `json:"tag"`
}

// SolCall represents a Solidity function call
type SolCall struct {
	FunctionName string     `json:"function_name"`
	Parameters   []AbiParam `json:"parameters"`
}

// ParsedCall represents a parsed function call for easier processing
type ParsedCall struct {
	FunctionName string
	Parameters   []ParsedParam
	Dst          string
	Src          string
	Value        string
	Gas          int64
	GasPrice     string
	HasDelay     bool   // Whether this call has an associated delay
	DelayValue   string // The delay value in seconds
}

// ParsedParam represents a parsed parameter
type ParsedParam struct {
	Type  string // Solidity type (uint256, uint8, etc.)
	Value string // The actual value
}

// ReplayGroup represents a group of calls that form one test function
type ReplayGroup struct {
	TestName string       // The name of the test function
	Calls    []ParsedCall // The sequence of calls in this replay
	FileName string       // Original file name for reference
}
