package tracer

var (
	// Loaded JS tracers for simulating various EntryPoint methods using debug_traceCall.
	Loaded, _ = NewTracers()
)
