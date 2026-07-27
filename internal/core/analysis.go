package core

type Command string

const (
	CommandSelect Command = "SELECT"
	CommandInsert Command = "INSERT"
	CommandUpdate Command = "UPDATE"
	CommandDelete Command = "DELETE"
)

type PrepareResult struct {
	Command    Command     `json:"command,omitempty"`
	Columns    []Column    `json:"columns"`
	Parameters []Parameter `json:"parameters"`
}

type ColumnSource struct {
	Schema     string `json:"schema,omitempty"`
	Table      string `json:"table,omitempty"`
	TableAlias string `json:"table_alias,omitempty"`
	Column     string `json:"column,omitempty"`
}

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

type Parameter struct {
	Number   int           `json:"number"`
	Name     string        `json:"name,omitempty"`
	DataType string        `json:"data_type,omitempty"`
	TypeOID  int64         `json:"type_oid,omitempty"`
	NotNull  bool          `json:"not_null"`
	Source   *ColumnSource `json:"source,omitempty"`
}
