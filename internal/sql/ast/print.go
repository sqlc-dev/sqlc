package ast

import (
	"strings"

	"github.com/sqlc-dev/sqlc/internal/debug"
	"github.com/sqlc-dev/sqlc/internal/sql/format"
)

type nodeFormatter interface {
	Format(*TrackedBuffer, format.Dialect)
}

type TrackedBuffer struct {
	*strings.Builder
}

// NewTrackedBuffer creates a new TrackedBuffer.
func NewTrackedBuffer() *TrackedBuffer {
	return &TrackedBuffer{
		Builder: new(strings.Builder),
	}
}

func (t *TrackedBuffer) astFormat(n Node, d format.Dialect) {
	if ft, ok := n.(nodeFormatter); ok {
		ft.Format(t, d)
	} else {
		debug.Dump(n)
	}
}

func (t *TrackedBuffer) join(n *List, d format.Dialect, sep string) {
	if n == nil {
		return
	}
	for i, item := range n.Items {
		if _, ok := item.(*TODO); ok {
			continue
		}
		if i > 0 {
			t.WriteString(sep)
		}
		t.astFormat(item, d)
	}
}

func Format(n Node, d format.Dialect) string {
	tb := NewTrackedBuffer()
	if ft, ok := n.(nodeFormatter); ok {
		ft.Format(tb, d)
	}
	return tb.String()
}

func set(n Node) bool {
	if n == nil {
		return false
	}
	_, ok := n.(*TODO)
	if ok {
		return false
	}
	return true
}

func items(n *List) bool {
	if n == nil {
		return false
	}
	return len(n.Items) > 0
}

func todo(n *List) bool {
	for _, item := range n.Items {
		if _, ok := item.(*TODO); !ok {
			return false
		}
	}
	return true
}
