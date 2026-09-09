// Package analysis is the shape of an engine's answer about a query: the
// JSON `sqlc analyze` prints, so that a case's committed stdout.json can be
// compared with what the database itself reports, byte for byte.
package analysis

import (
	"bytes"
	"encoding/json"
)

// Query is what was found out about one query.
type Query struct {
	Name    string   `json:"name"`
	Cmd     string   `json:"cmd"`
	Columns []Column `json:"columns"`
	Params  []Param  `json:"params"`
}

// Column describes a result column, or the column a parameter stands in
// for.
type Column struct {
	Name  string    `json:"name"`
	Type  *TypeExpr `json:"type,omitempty"`
	Table string    `json:"table,omitempty"`
}

// Param is one parameter and what it is compared with or assigned to.
type Param struct {
	Number int    `json:"number"`
	Column Column `json:"column"`
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
func Encode(queries []Query) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(queries); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
