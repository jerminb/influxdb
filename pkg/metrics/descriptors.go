package metrics

type GroupDesc struct {
	Name string
	id   GID
}

type MetricType int

const (
	CounterMetricType MetricType = iota
	TimerMetricType
)

type Desc struct {
	Name string
	mt   MetricType
	gid  GID
	id   ID
}

type descOption func(*Desc)

// WithGroup assigns the associated measurement to group identified by gid.
func WithGroup(gid GID) descOption {
	return func(d *Desc) {
		d.gid = gid
	}
}

func newDesc(name string, opts ...descOption) *Desc {
	desc := &Desc{Name: name}
	for _, o := range opts {
		o(desc)
	}
	return desc
}

const (
	idMask   = (1 << 32) - 1
	gidShift = 32
)

type (
	GID uint32
	ID  uint64
)

func newID(id int, gid GID) ID {
	return ID(gid)<<gidShift + (ID(id) & idMask)
}

func (id ID) id() int {
	return int(id & idMask)
}

func (id ID) gid() int {
	return int(id >> gidShift)
}

func (id *ID) setGID(gid GID) {
	*id |= ID(gid) << gidShift
}
