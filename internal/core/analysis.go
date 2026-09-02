package core

type Command string

const (
	CommandSelect Command = "SELECT"
	CommandInsert Command = "INSERT"
	CommandUpdate Command = "UPDATE"
	CommandDelete Command = "DELETE"
)

type PrepareResult struct {
	Command    Command         `json:"command,omitempty"`
	Columns    []Column        `json:"columns"`
	Parameters []Parameter     `json:"parameters"`
	Stars      []StarExpansion `json:"stars,omitempty"`
}

// StarExpansion is what a star in a target list stands for. The analyzer
// resolves the reference against the query's scope and reports the columns it
// covers; rewriting the query text with them is the caller's to do, since only
// it knows how the engine quotes an identifier.
type StarExpansion struct {
	// Location is where the target the star belongs to starts, measured the
	// way the AST measures a node: from the beginning of the file the
	// statement was parsed from.
	Location int `json:"location"`

	// Fields is the reference as it was written, with the star as its last
	// element: ["*"] for a bare star and ["foo", "*"] for a qualified one.
	Fields []string `json:"fields"`

	// Alias is the output name the target was given, if any.
	Alias string `json:"alias,omitempty"`

	Columns []StarColumn `json:"columns"`
}

// StarColumn is a single column a star expanded to.
type StarColumn struct {
	// Relation is the name the column's relation goes by in the query, which
	// is its alias when it was given one.
	Relation string `json:"relation,omitempty"`
	Name     string `json:"name"`
	DataType string `json:"data_type,omitempty"`
}

type ColumnSource struct {
	Schema     string `json:"schema,omitempty"`
	Table      string `json:"table,omitempty"`
	TableAlias string `json:"table_alias,omitempty"`
	Column     string `json:"column,omitempty"`
}

type Column struct {
	Name     string `json:"name"`
	DataType string `json:"data_type"`
	// Type is the column's type as an expression, carrying what DataType
	// and IsArray flatten away: arguments, nesting and inner nullability.
	Type               *TypeExpr     `json:"type,omitempty"`
	TypeOID            int64         `json:"type_oid,omitempty"`
	NotNull            bool          `json:"not_null"`
	IsArray            bool          `json:"is_array,omitempty"`
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

type Parameter struct {
	Number   int           `json:"number"`
	Name     string        `json:"name,omitempty"`
	DataType string        `json:"data_type,omitempty"`
	Type     *TypeExpr     `json:"type,omitempty"`
	TypeOID  int64         `json:"type_oid,omitempty"`
	NotNull  bool          `json:"not_null"`
	IsArray  bool          `json:"is_array,omitempty"`
	Source   *ColumnSource `json:"source,omitempty"`
}
