package endtoend

import (
	"bytes"
	"encoding/json"
)

// An engine's answer is written in the JSON `sqlc analyze` prints, so that
// a case's committed stdout.json can be compared with it byte for byte.

// AnalyzedQuery is what was found out about one query.
type AnalyzedQuery struct {
	Name    string           `json:"name"`
	Cmd     string           `json:"cmd"`
	Columns []AnalyzedColumn `json:"columns"`
	Params  []AnalyzedParam  `json:"params"`
}

// AnalyzedColumn describes a result column, or the column a parameter
// stands in for.
type AnalyzedColumn struct {
	Name  string    `json:"name"`
	Type  *TypeExpr `json:"type,omitempty"`
	Table string    `json:"table,omitempty"`
}

// AnalyzedParam is one parameter and what it is compared with or assigned
// to.
type AnalyzedParam struct {
	Number int            `json:"number"`
	Column AnalyzedColumn `json:"column"`
}

// TypeExpr is a type as a call expression: a lowercased name applied to an
// ordered argument list, each argument another type, an integer, a boolean
// or a quoted string, optionally labelled. Nullability is an attribute of
// the type rather than a wrapper.
type TypeExpr struct {
	Name     string    `json:"name"`
	Nullable bool      `json:"nullable,omitempty"`
	Args     []TypeArg `json:"args,omitempty"`
}

// TypeArg is one argument of a TypeExpr.
type TypeArg struct {
	Label  string    `json:"label,omitempty"`
	Type   *TypeExpr `json:"type,omitempty"`
	Int    *int64    `json:"int,omitempty"`
	Bool   *bool     `json:"bool,omitempty"`
	String *string   `json:"string,omitempty"`
}

// Encode prints the answer the way sqlc analyze does.
func Encode(queries []AnalyzedQuery) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(queries); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
