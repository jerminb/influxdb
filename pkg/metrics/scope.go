package metrics

import (
	"fmt"
	"sync"

	"github.com/xlab/treeprint"
)

type Scope struct {
	name string
	r    *Registry
	p    *Scope

	mu       sync.Mutex
	c        []*Group
	children []*Scope
	labels   LabelSet
}

func (s *Scope) Name() string       { return s.name }
func (s *Scope) Labels() LabelSet   { return s.labels }
func (s *Scope) Groups() []*Group   { return s.c }
func (s *Scope) Children() []*Scope { return s.children }

func (s *Scope) String() string {
	p := &printer{}
	Walk(p, s)
	return p.String()
}

func (s *Scope) NewGroup(gid GID) *Group {
	c := s.r.mustGetGroupRegistry(gid).newGroup()

	s.mu.Lock()
	s.c = append(s.c, c)
	s.mu.Unlock()

	return c
}

// NewScope creates a child scope for collecting metrics.
func (s *Scope) NewScope(name string, opts ...scopeOption) *Scope {
	cs := &Scope{name: name, r: s.r, p: s}
	for _, o := range opts {
		o(cs)
	}

	s.mu.Lock()
	s.children = append(s.children, cs)
	s.mu.Unlock()

	return cs
}

func (s *Scope) AddLabels(labels LabelSet) {
	s.mu.Lock()
	s.labels.Merge(labels)
	s.mu.Unlock()
}

func (s *Scope) Tree() treeprint.Tree {
	t := newTreeVisitor()
	Walk(t, s)
	return t.root
}

type scopeOption func(*Scope)

func ScopeLabels(labels LabelSet) scopeOption {
	return func(scope *Scope) {
		scope.labels = labels
	}
}

type treeVisitor struct {
	root  treeprint.Tree
	trees []treeprint.Tree
}

func newTreeVisitor() *treeVisitor {
	t := treeprint.New()
	return &treeVisitor{root: t, trees: []treeprint.Tree{t}}
}

func (v *treeVisitor) Visit(node Node) Visitor {
	switch n := node.(type) {
	case *Scope:
		t := v.trees[len(v.trees)-1].AddBranch(n.Name())
		v.trees = append(v.trees, t)

		if labels := n.Labels(); labels.Len() > 0 {
			l := t.AddBranch("labels")
			n.Labels().ForEach(func(k, v string) {
				l.AddNode(k + ": " + v)
			})
		}

		for _, cn := range n.c {
			Walk(v, cn)
		}
		for _, cn := range n.children {
			Walk(v, cn)
		}

		v.trees[len(v.trees)-1] = nil
		v.trees = v.trees[:len(v.trees)-1]

		return nil

	case *Counter, *Timer:
		v.trees[len(v.trees)-1].AddNode(n.(fmt.Stringer).String())
		return nil
	}

	return v
}
