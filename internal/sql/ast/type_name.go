package ast

type TypeName struct {
	Catalog string
	Schema  string
	Name    string

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

// mapTypeName converts internal PostgreSQL type names to their SQL equivalents
func mapTypeName(name string) string {
	switch name {
	case "int4":
		return "integer"
	case "int8":
		return "bigint"
	case "int2":
		return "smallint"
	case "float4":
		return "real"
	case "float8":
		return "double precision"
	case "bool":
		return "boolean"
	case "bpchar":
		return "character"
	case "timestamptz":
		return "timestamp with time zone"
	case "timetz":
		return "time with time zone"
	default:
		return name
	}
}

func (n *TypeName) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	if items(n.Names) {
		// Check if this is a pg_catalog type that should be expanded
		if len(n.Names.Items) == 2 {
			first, _ := n.Names.Items[0].(*String)
			second, _ := n.Names.Items[1].(*String)
			if first != nil && second != nil && first.Str == "pg_catalog" {
				// pg_catalog.timestamptz -> timestamp with time zone
				// pg_catalog.timetz -> time with time zone
				// etc.
				buf.WriteString(mapTypeName(second.Str))
				goto addMods
			}
		}
		// For single name types, just output as-is (don't expand)
		if len(n.Names.Items) == 1 {
			if s, ok := n.Names.Items[0].(*String); ok {
				buf.WriteString(s.Str)
				goto addMods
			}
		}
		buf.join(n.Names, ".")
	} else if n.Schema == "pg_catalog" {
		// pg_catalog.typename -> expanded form (via Schema/Name fields)
		buf.WriteString(mapTypeName(n.Name))
	} else if n.Schema != "" {
		// schema.typename
		buf.WriteString(n.Schema)
		buf.WriteString(".")
		buf.WriteString(n.Name)
	} else {
		// Simple type name - don't expand aliases
		buf.WriteString(n.Name)
	}
addMods:
	// Add type modifiers (e.g., varchar(255))
	if items(n.Typmods) {
		buf.WriteString("(")
		buf.join(n.Typmods, ", ")
		buf.WriteString(")")
	}
	if items(n.ArrayBounds) {
		buf.WriteString("[]")
	}
}
