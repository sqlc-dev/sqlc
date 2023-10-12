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

func (n *TypeName) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	if items(n.Names) {
		buf.join(n.Names, ".")
	} else {
		if n.Name == "int4" {
			buf.WriteString("INTEGER")
		} else {
			buf.WriteString(n.Name)
		}
	}
	if items(n.ArrayBounds) {
		buf.WriteString("[]")
	}
}
