package metrics

//go:generate stringer -type=MetricType

import (
	"fmt"
	"sort"
)

type Registry struct {
	descriptors []*GroupDesc
	groups      []groupRegistry
}

const (
	// DefaultGroup is the identifier for the default group.
	DefaultGroup = GID(0)
)

// NewRegistry creates a new Registry with a single group identified by DefaultGroup.
func NewRegistry() *Registry {
	var r Registry
	r.MustRegisterGroup("global")
	return &r
}

func (r *Registry) register(gd *GroupDesc) error {
	p := sort.Search(len(r.descriptors), func(i int) bool {
		return r.descriptors[i].Name == gd.Name
	})

	if p != len(r.descriptors) {
		return fmt.Errorf("group name '%s' already in use", gd.Name)
	}

	r.descriptors = append(r.descriptors, gd)
	sort.Slice(r.descriptors, func(i, j int) bool {
		return r.descriptors[i].Name < r.descriptors[j].Name
	})

	gd.id = GID(len(r.groups))
	r.groups = append(r.groups, groupRegistry{desc: gd})

	return nil
}

func (r *Registry) mustRegister(gd *GroupDesc) {
	if err := r.register(gd); err != nil {
		panic(err.Error())
	}
}

// MustRegisterGroup registers a new group and panics if a group already exists with the same name.
//
// MustRegisterGroup is not safe to call from concurrent goroutines.
func (r *Registry) MustRegisterGroup(name string) GID {
	gd := &GroupDesc{Name: name}
	r.mustRegister(gd)
	return gd.id
}

func (r *Registry) mustGetGroupRegistry(id GID) *groupRegistry {
	if int(id) >= len(r.groups) {
		panic(fmt.Sprintf("invalid group ID"))
	}
	return &r.groups[id]
}

// MustRegisterCounter registers a new counter metric using the provided descriptor.
// If the metric name is not unique within the group, MustRegisterCounter will panic.
//
// MustRegisterCounter is not safe to call from concurrent goroutines.
func (r *Registry) MustRegisterCounter(name string, opts ...descOption) ID {
	desc := newDesc(name, opts...)
	return r.mustGetGroupRegistry(desc.gid).mustRegisterCounter(desc)
}

// MustRegisterTimer registers a new timer metric using the provided descriptor.
// If the metric name is not unique within the group, MustRegisterTimer will panic.
//
// MustRegisterTimer is not safe to call from concurrent goroutines.
func (r *Registry) MustRegisterTimer(name string, opts ...descOption) ID {
	desc := newDesc(name, opts...)
	return r.mustGetGroupRegistry(desc.gid).mustRegisterTimer(desc)
}

// NewScope create a new Scope for collecting a metrics.
//
// NewScope is safe to call from concurrent goroutines.
func (r *Registry) NewScope(name string, opts ...scopeOption) *Scope {
	s := &Scope{name: name, r: r}
	for _, o := range opts {
		o(s)
	}

	return s
}
