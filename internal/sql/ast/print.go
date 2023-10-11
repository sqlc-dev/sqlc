package ast

import (
	"strings"

	"github.com/sqlc-dev/sqlc/internal/debug"
)

type formatter interface {
	Format(*TrackedBuffer)
}

type TrackedBuffer struct {
	*strings.Builder
}

// NewTrackedBuffer creates a new TrackedBuffer.
func NewTrackedBuffer() *TrackedBuffer {
	buf := &TrackedBuffer{
		Builder: new(strings.Builder),
	}
	return buf
}

func (t *TrackedBuffer) astFormat(n Node) {
	if ft, ok := n.(formatter); ok {
		ft.Format(t)
	} else {
		debug.Dump(n)
	}
}

func (t *TrackedBuffer) join(n *List, sep string) {
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
		t.astFormat(item)
	}
}

func Format(n Node) string {
	tb := NewTrackedBuffer()
	if ft, ok := n.(formatter); ok {
		ft.Format(tb)
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
