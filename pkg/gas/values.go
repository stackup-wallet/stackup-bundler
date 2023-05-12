package gas

var (
	// The maximum value to start a binary search from. A lower value reduces the upper limit of RPC calls to
	// the debug_traceCall method. This value is the current gas limit of a block.
	MaxGasLimit = 30000000
)
