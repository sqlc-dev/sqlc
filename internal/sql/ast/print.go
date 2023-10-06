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

func Format(n Node) string {
	tb := NewTrackedBuffer()
	if ft, ok := n.(formatter); ok {
		ft.Format(tb)
	}
	return tb.String()
}
