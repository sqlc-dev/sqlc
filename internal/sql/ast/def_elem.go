package ast

type DefElem struct {
	Defnamespace *string
	Defname      *string
	Arg          Node
	Defaction    DefElemAction
	Location     int
}

func (n *DefElem) Pos() int {
	return n.Location
}

func (n *DefElem) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	if n.Defname != nil {
		switch *n.Defname {
		case "as":
			buf.WriteString("AS ")
			// AS clause contains function body which needs quoting
			if l, ok := n.Arg.(*List); ok {
				for i, item := range l.Items {
					if i > 0 {
						buf.WriteString(", ")
					}
					if s, ok := item.(*String); ok {
						buf.WriteString("'")
						buf.WriteString(s.Str)
						buf.WriteString("'")
					} else {
						buf.astFormat(item)
					}
				}
			} else {
				buf.astFormat(n.Arg)
			}
		case "language":
			buf.WriteString("LANGUAGE ")
			buf.astFormat(n.Arg)
		case "volatility":
			// VOLATILE, STABLE, IMMUTABLE
			buf.astFormat(n.Arg)
		case "strict":
			if s, ok := n.Arg.(*Boolean); ok && s.Boolval {
				buf.WriteString("STRICT")
			} else {
				buf.WriteString("CALLED ON NULL INPUT")
			}
		case "security":
			if s, ok := n.Arg.(*Boolean); ok && s.Boolval {
				buf.WriteString("SECURITY DEFINER")
			} else {
				buf.WriteString("SECURITY INVOKER")
			}
		default:
			buf.WriteString(*n.Defname)
			if n.Arg != nil {
				buf.WriteString(" ")
				buf.astFormat(n.Arg)
			}
		}
	}
}
