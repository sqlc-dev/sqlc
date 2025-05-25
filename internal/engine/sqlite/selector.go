package sqlite

type Selector struct{}

func NewSelector() *Selector {
	return &Selector{}
}

func (s *Selector) ColumnExpr(name string, dataType string) string {
	// Under SQLite, neither json nor jsonb are real data types, and rather just
	// of type blob, so database drivers just return whatever raw binary is
	// stored as values. This is a problem for jsonb, which is considered an
	// internal format to SQLite and no attempt should be made to parse it
	// outside of the database itself. For jsonb columns in SQLite, wrap values
	// in `json(col)` to coerce the internal binary format to JSON parsable by
	// the user-space application.
	if dataType == "jsonb" {
		return "json(" + name + ")"
	}
	return name
}
