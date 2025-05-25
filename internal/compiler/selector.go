package compiler

// selector is an interface used by a compiler for generating expressions for
// output columns in a `SELECT ...` or `RETURNING ...` statement.
//
// This interface is exclusively needed at the moment for SQLite, which must
// wrap output `jsonb` columns with a `json(column_name)` invocation so that a
// publicly consumable format (i.e. not jsonb) is returned.
type selector interface {
	// ColumnExpr generates output to be used in a `SELECT ...` or `RETURNING
	// ...` statement based on input column name and metadata.
	ColumnExpr(name string, column *Column) string
}

// defaultSelector is a selector implementation that does the simpliest possible
// pass through when generating column expressions. Its use is suitable for all
// database engines not requiring additional customization.
type defaultSelector struct{}

func newDefaultSelector() *defaultSelector {
	return &defaultSelector{}
}

func (s *defaultSelector) ColumnExpr(name string, column *Column) string {
	return name
}

type sqliteSelector struct{}

func newSQLiteSelector() *sqliteSelector {
	return &sqliteSelector{}
}

func (s *sqliteSelector) ColumnExpr(name string, column *Column) string {
	// Under SQLite, neither json nor jsonb are real data types, and rather just
	// of type blob, so database drivers just return whatever raw binary is
	// stored as values. This is a problem for jsonb, which is considered an
	// internal format to SQLite and no attempt should be made to parse it
	// outside of the database itself. For jsonb columns in SQLite, wrap values
	// in `json(col)` to coerce the internal binary format to JSON parsable by
	// the user-space application.
	if column.DataType == "jsonb" {
		return "json(" + name + ")"
	}
	return name
}
