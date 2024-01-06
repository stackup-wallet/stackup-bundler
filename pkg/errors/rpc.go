package errors

var (
	REJECTED_BY_EP_OR_ACCOUNT  = -32500
	REJECTED_BY_PAYMASTER      = -32501
	BANNED_OPCODE              = -32502
	SHORT_DEADLINE             = -32503
	BANNED_OR_THROTTLED_ENTITY = -32504
	INVALID_ENTITY_STAKE       = -32505
	INVALID_AGGREGATOR         = -32506
	INVALID_SIGNATURE          = -32507
	INVALID_FIELDS             = -32602

	EXECUTION_REVERTED = -32521
)

// RPCError is a custom error that fits the JSON-RPC error spec.
type RPCError struct {
	code    int
	message string
	data    any
}

// New returns a new custom RPCError.
func NewRPCError(code int, message string, data any) error {
	return &RPCError{code, message, data}
}

// Code returns the message field of the JSON-RPC error object.
func (e *RPCError) Error() string {
	return e.message
}

// Data returns the data field of the JSON-RPC error object.
func (e *RPCError) Data() any {
	return e.data
}

// Code returns the code field of the JSON-RPC error object.
func (e *RPCError) Code() int {
	return e.code
}
