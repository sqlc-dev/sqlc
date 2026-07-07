package core

// Command identifies the kind of statement that produced a PrepareResult.
// Only the four DML statements that can have a prepare-able shape are
// emitted; DDL and TCL produce an empty result with Command == "".
type Command string

const (
	CommandSelect Command = "SELECT"
	CommandInsert Command = "INSERT"
	CommandUpdate Command = "UPDATE"
	CommandDelete Command = "DELETE"
)

// PrepareResult describes the output of preparing a SQL statement.
type PrepareResult struct {
	Command    Command     `json:"command,omitempty"`
	Columns    []Column    `json:"columns"`
	Parameters []Parameter `json:"parameters"`
}

// ColumnSource identifies the table column a result column or bind
// parameter is sourced from. All fields are optional; the struct is
// emitted only when at least one is populated.
//
// Schema / Table / Column are the *origin* identifiers (pre-alias) —
// they correspond to sqlite's `sqlite3_column_origin_name` and mysql's
// `org_table` / `org_name`. TableAlias is the name the query used to
// refer to the table; it lets codegen distinguish `t1` from `t2` in
// `SELECT t1.x, t2.x FROM t t1 JOIN t t2 ...`.
type ColumnSource struct {
	Schema     string `json:"schema,omitempty"`
	Table      string `json:"table,omitempty"`
	TableAlias string `json:"table_alias,omitempty"`
	Column     string `json:"column,omitempty"`
}

// Column describes a single output column from a prepared statement.
//
// SourceClassOID and SourceAttributeOID are populated when the column
// is a direct reference to a table column (i.e. it appears in
// sql_attribute); they are zero for computed/derived expressions like
// aggregates or arithmetic. Source is the human-readable resolution of
// those OIDs and is populated under the same conditions.
//
// DeclType, TypeLength, TypeScale, IsPrimaryKey, IsUnique, and
// IsAutoIncrement come from the resolved source attribute and are zero
// for computed expressions.
type Column struct {
	Name               string        `json:"name"`
	DataType           string        `json:"data_type"`
	TypeOID            int64         `json:"type_oid,omitempty"`
	NotNull            bool          `json:"not_null"`
	SourceClassOID     int64         `json:"source_class_oid,omitempty"`
	SourceAttributeOID int64         `json:"source_attribute_oid,omitempty"`
	Source             *ColumnSource `json:"source,omitempty"`
	DeclType           string        `json:"decl_type,omitempty"`
	TypeLength         int           `json:"type_length,omitempty"`
	TypeScale          int           `json:"type_scale,omitempty"`
	IsPrimaryKey       bool          `json:"is_primary_key,omitempty"`
	IsUnique           bool          `json:"is_unique,omitempty"`
	IsAutoIncrement    bool          `json:"is_auto_increment,omitempty"`
}

// Parameter describes a bind parameter in a prepared statement.
//
// Number is the 1-based position of the parameter as it appeared in the
// source ($1, $2, ...). Name is populated for named-parameter dialects
// or when a sqlc-style "-- name:" annotation gave the param a name; it
// is empty otherwise. DataType / TypeOID / NotNull come from the
// resolved usage site (an operator overload, function argument, etc.).
// Source identifies the column the parameter binds against (e.g.
// `users.age` for `WHERE age > $1`), when one can be inferred.
type Parameter struct {
	Number   int           `json:"number"`
	Name     string        `json:"name,omitempty"`
	DataType string        `json:"data_type,omitempty"`
	TypeOID  int64         `json:"type_oid,omitempty"`
	NotNull  bool          `json:"not_null"`
	Source   *ColumnSource `json:"source,omitempty"`
}
