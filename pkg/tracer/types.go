// These are the same BundlerCollectorTracer types from github.com/eth-infinitism/bundler ported for Go.

package tracer

import "github.com/ethereum/go-ethereum/common"

type Counts = map[string]float64

type AccessMap = map[common.Address]AccessInfo

// AccessInfo provides context on read and write counts by storage slots.
type AccessInfo struct {
	Reads  Counts `json:"reads"`
	Writes Counts `json:"writes"`
}

// NumberLevelInfo provides context on opcodes and storage access delimited by the use of NUMBER at the
// EntryPoint.
type NumberLevelInfo struct {
	Opcodes Counts    `json:"opcodes"`
	Access  AccessMap `json:"access"`
}

// CallInfo provides context on internal calls made during tracing.
type CallInfo struct {
	// Common
	Type string `json:"type"`

	// Method info
	From   common.Address `json:"from"`
	To     common.Address `json:"to"`
	Method string         `json:"method"`
	Value  any            `json:"value"`
	Gas    float64        `json:"gas"`

	// Exit info
	GasUsed float64 `json:"gasUsed"`
	Data    any     `json:"data"`
}

// LogInfo provides context from LOG opcodes during each step in the EVM trace.
type LogInfo struct {
	Topics []string `json:"topics"`
	Data   string   `json:"data"`
}

// BundlerCollectorReturn is the return value from performing an EVM trace with BundlerCollectorTracer.js.
type BundlerCollectorReturn struct {
	NumberLevels []NumberLevelInfo `json:"numberLevels"`
	Keccak       []string          `json:"keccak"`
	Calls        []CallInfo        `json:"calls"`
	Logs         []LogInfo         `json:"logs"`
	Debug        []any             `json:"debug"`
}
