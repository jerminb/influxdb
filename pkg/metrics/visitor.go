package metrics

import (
	"bytes"
	"fmt"
	"strings"
)

type Node interface {
	Name() string
}

// A Visitor's Visit method is invoked for each node encountered by Walk.
// If the result of Visit is not nil, Walk visits each of the children.
type Visitor interface {
	Visit(Node) Visitor
}

// Walk traverses a Scope, Group and each measurement in depth-first order.
func Walk(v Visitor, node Node) {
	if v = v.Visit(node); v == nil {
		return
	}

	switch n := node.(type) {
	case *Scope:
		for _, coll := range n.c {
			Walk(v, coll)
		}

		for _, child := range n.children {
			Walk(v, child)
		}

	case *Group:
		for i := range n.counters {
			Walk(v, &n.counters[i])
		}
		for i := range n.timers {
			Walk(v, &n.timers[i])
		}
	}
}

type walkFuncVisitor func(Node) bool

func (fn walkFuncVisitor) Visit(n Node) Visitor {
	if fn(n) {
		return fn
	} else {
		return nil
	}
}

func WalkFunc(node Node, fn func(Node) bool) {
	Walk(walkFuncVisitor(fn), node)
}

type printer struct {
	level  int
	prefix string
	buf    bytes.Buffer
}

func (p *printer) String() string { return p.buf.String() }

func (p *printer) Visit(node Node) Visitor {
	switch n := node.(type) {
	case *Scope:
		prefix := strings.Repeat(" ", p.level*2)
		fmt.Fprint(&p.buf, prefix, "→ ", n.Name(), "\n")
		p.level++
		old := p.prefix
		p.prefix = strings.Repeat(" ", p.level*2+2)

		if labels := n.Labels(); labels.Len() > 0 {
			fmt.Fprint(&p.buf, p.prefix, "labels", "\n")
			prefix = p.prefix + "  "
			n.Labels().ForEach(func(k, v string) {
				fmt.Fprint(&p.buf, prefix, k, ": ", v, "\n")
			})
		}

		for _, cn := range n.c {
			Walk(p, cn)
		}
		for _, cn := range n.children {
			Walk(p, cn)
		}

		p.level--
		p.prefix = old
		return nil

	case *Counter, *Timer:
		fmt.Fprint(&p.buf, p.prefix, n, "\n")
		return nil
	}

	return p
}
