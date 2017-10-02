package metrics

// The Group type represents an instance of a set of metrics that are used for
// instrumenting a specific request.
type Group struct {
	g        *groupRegistry
	counters []Counter
	timers   []Timer
}

// Name returns the name of the group.
func (c *Group) Name() string { return c.g.desc.Name }

// GetCounter returns the counter identified by the id that was returned
// by MustRegisterCounter for the same group.
// Using an id from a different group will result in undefined behavior.
func (c *Group) GetCounter(id ID) *Counter { return &c.counters[id.id()] }

// GetTimer returns the timer identified by the id that was returned
// by MustRegisterTimer for the same group.
// Using an id from a different group will result in undefined behavior.
func (c *Group) GetTimer(id ID) *Timer { return &c.timers[id.id()] }
