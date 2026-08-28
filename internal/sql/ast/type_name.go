package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type TypeName struct {
	Tag NodeTag[TypeName] `json:"tag"`

	Catalog string
	Schema  string
	Name    string
	// Spelling is the type as the author wrote it, when that differs from
	// Name: SQLite folds the spaces out of multi-word type names ("VARYING
	// CHARACTER" resolves as "VARYINGCHARACTER" in the catalog), so the
	// formatter prints this back instead of the folded form.
	Spelling string

	// From pg.TypeName
	Names       *List
	TypeOid     Oid
	Setof       bool
	PctType     bool
	Typmods     *List
	Typemod     int32
	ArrayBounds *List
	Location    int
}

func (n *TypeName) Pos() int {
	return n.Location
}

func (n *TypeName) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}
	if n.Spelling != "" {
		buf.WriteString(n.Spelling)
		goto addMods
	}
	if items(n.Names) {
		// Check if this is a qualified type (e.g., pg_catalog.int4)
		if len(n.Names.Items) == 2 {
			first, _ := n.Names.Items[0].(*String)
			second, _ := n.Names.Items[1].(*String)
			if first != nil && second != nil {
				buf.WriteString(d.TypeName(first.Str, second.Str))
				goto addMods
			}
		}
		// For single name types, just output as-is
		if len(n.Names.Items) == 1 {
			if s, ok := n.Names.Items[0].(*String); ok {
				buf.WriteString(d.TypeName("", s.Str))
				goto addMods
			}
		}
		buf.join(n.Names, d, ".")
	} else {
		buf.WriteString(d.TypeName(n.Schema, n.Name))
	}
addMods:
	// Add type modifiers (e.g., varchar(255))
	if items(n.Typmods) {
		buf.WriteString("(")
		buf.join(n.Typmods, d, ", ")
		buf.WriteString(")")
	}
	if items(n.ArrayBounds) {
		buf.WriteString("[]")
	}
}
