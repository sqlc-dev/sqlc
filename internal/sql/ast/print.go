package ast

import (
	"fmt"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/debug"
	"github.com/sqlc-dev/sqlc/internal/sql/format"
)

type nodeFormatter interface {
	Format(*TrackedBuffer)
}

type TrackedBuffer struct {
	*strings.Builder
	formatter format.Formatter
}

// NewTrackedBuffer creates a new TrackedBuffer with the given formatter.
func NewTrackedBuffer(f format.Formatter) *TrackedBuffer {
	buf := &TrackedBuffer{
		Builder:   new(strings.Builder),
		formatter: f,
	}
	return buf
}

// QuoteIdent returns a quoted identifier if it needs quoting.
// If no formatter is set, it returns the identifier unchanged.
func (t *TrackedBuffer) QuoteIdent(s string) string {
	if t.formatter != nil {
		return t.formatter.QuoteIdent(s)
	}
	return s
}

// TypeName returns the SQL type name for the given namespace and name.
// If no formatter is set, it returns "ns.name" or just "name".
func (t *TrackedBuffer) TypeName(ns, name string) string {
	if t.formatter != nil {
		return t.formatter.TypeName(ns, name)
	}
	if ns != "" {
		return ns + "." + name
	}
	return name
}

// Param returns the parameter placeholder for the given number.
// If no formatter is set, it returns PostgreSQL-style $n.
func (t *TrackedBuffer) Param(n int) string {
	if t.formatter != nil {
		return t.formatter.Param(n)
	}
	return fmt.Sprintf("$%d", n)
}

// Cast returns a type cast expression.
// If no formatter is set, it returns PostgreSQL-style expr::type.
func (t *TrackedBuffer) Cast(arg, typeName string) string {
	if t.formatter != nil {
		return t.formatter.Cast(arg, typeName)
	}
	return arg + "::" + typeName
}

// NamedParam returns the named parameter placeholder for the given name.
// If no formatter is set, it returns PostgreSQL-style @name.
func (t *TrackedBuffer) NamedParam(name string) string {
	if t.formatter != nil {
		return t.formatter.NamedParam(name)
	}
	return "@" + name
}

func (t *TrackedBuffer) astFormat(n Node) {
	if ft, ok := n.(nodeFormatter); ok {
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

func Format(n Node, f format.Formatter) string {
	tb := NewTrackedBuffer(f)
	if ft, ok := n.(nodeFormatter); ok {
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
