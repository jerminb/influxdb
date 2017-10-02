package metrics

import (
	"fmt"
	"sort"
)

// The groupRegistry type represents a set of metrics that are measured together.
type groupRegistry struct {
	desc        *GroupDesc
	descriptors []*Desc
	counters    []Counter
	timers      []Timer
}

func (g *groupRegistry) register(desc *Desc) error {
	p := sort.Search(len(g.descriptors), func(i int) bool {
		return g.descriptors[i].Name == desc.Name
	})

	if p != len(g.descriptors) {
		return fmt.Errorf("metric name '%s' already in use", desc.Name)
	}

	g.descriptors = append(g.descriptors, desc)
	sort.Slice(g.descriptors, func(i, j int) bool {
		return g.descriptors[i].Name < g.descriptors[j].Name
	})

	return nil
}

func (g *groupRegistry) mustRegister(desc *Desc) {
	if err := g.register(desc); err != nil {
		panic(err.Error())
	}
}

// MustRegisterCounter registers a new counter metric using the provided descriptor.
// If the metric name is not unique, MustRegisterCounter will panic.
//
// MustRegisterCounter is not safe to call from multiple goroutines.
func (g *groupRegistry) mustRegisterCounter(desc *Desc) ID {
	desc.mt = CounterMetricType
	g.mustRegister(desc)

	desc.id = newID(len(g.counters), g.desc.id)
	g.counters = append(g.counters, Counter{desc: desc})

	return desc.id
}

// MustRegisterTimer registers a new timer metric using the provided descriptor.
// If the metric name is not unique, MustRegisterTimer will panic.
//
// MustRegisterTimer is not safe to call from multiple goroutines.
func (g *groupRegistry) mustRegisterTimer(desc *Desc) ID {
	desc.mt = TimerMetricType
	g.mustRegister(desc)

	desc.id = newID(len(g.timers), g.desc.id)
	g.timers = append(g.timers, Timer{desc: desc})

	return desc.id
}

// newCollector returns a Collector with a copy of all the registered counters.
//
// newCollector is safe to call from multiple goroutines.
func (g *groupRegistry) newGroup() *Group {
	c := &Group{
		g:        g,
		counters: make([]Counter, len(g.counters)),
		timers:   make([]Timer, len(g.timers)),
	}

	copy(c.counters, g.counters)
	copy(c.timers, g.timers)
	return c
}
