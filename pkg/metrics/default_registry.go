package metrics

var DefaultRegistry = NewRegistry()

func MustRegisterGroup(name string) GID {
	return DefaultRegistry.MustRegisterGroup(name)
}

// MustRegisterCounter registers a new counter metric with the DefaultRegistry
// using the provided descriptor.
// If the metric name is not unique, MustRegisterCounter will panic.
//
// MustRegisterCounter is not safe to call from multiple goroutines.
func MustRegisterCounter(name string, opts ...descOption) ID {
	return DefaultRegistry.MustRegisterCounter(name, opts...)
}

// MustRegisterTimer registers a new timer metric with the DefaultRegistry
// using the provided descriptor.
// If the metric name is not unique, MustRegisterTimer will panic.
//
// MustRegisterTimer is not safe to call from multiple goroutines.
func MustRegisterTimer(name string, opts ...descOption) ID {
	return DefaultRegistry.MustRegisterTimer(name, opts...)
}

// NewCollector returns a new collector from the DefaultRegistry
func NewScope(name string) *Scope {
	return DefaultRegistry.NewScope(name)
}
