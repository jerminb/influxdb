package metrics

import "context"

type key int

const (
	groupKey key = iota
	scopeKey
)

// NewContextWithGroup returns a new context with the given Group added.
func NewContextWithGroup(ctx context.Context, c *Group) context.Context {
	return context.WithValue(ctx, groupKey, c)
}

// GroupFromContext returns the Group associated with ctx or nil if no Group has been assigned.
func GroupFromContext(ctx context.Context) *Group {
	c, _ := ctx.Value(groupKey).(*Group)
	return c
}

// NewContextWithScope returns a new context with the given Collector added.
func NewContextWithScope(ctx context.Context, c *Scope) context.Context {
	return context.WithValue(ctx, scopeKey, c)
}

// ScopeFromContext returns the Scope associated with ctx or nil if no Scope has been assigned.
func ScopeFromContext(ctx context.Context) *Scope {
	c, _ := ctx.Value(scopeKey).(*Scope)
	return c
}
