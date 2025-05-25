package selector

// Selector is an interface used by a compiler for generating expressions for
// output columns in a `SELECT ...` or `RETURNING ...` statement.
//
// This interface is exclusively needed at the moment for SQLite, which must
// wrap output `jsonb` columns with a `json(column_name)` invocation so that a
// publicly consumable format (i.e. not jsonb) is returned.
type Selector interface {
	// ColumnExpr generates output to be used in a `SELECT ...` or `RETURNING
	// ...` statement based on input column name and metadata.
	ColumnExpr(name string, dataType string) string
}

// DefaultSelector is a Selector implementation that does the simpliest possible
// pass through when generating column expressions. Its use is suitable for all
// database engines not requiring additional customization.
type DefaultSelector struct{}

func NewDefaultSelector() *DefaultSelector {
	return &DefaultSelector{}
}

func (s *DefaultSelector) ColumnExpr(name string, dataType string) string {
	return name
}
