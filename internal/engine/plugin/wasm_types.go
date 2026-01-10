// Package plugin provides JSON types for WASM engine plugins.
// WASM plugins use JSON instead of Protobuf because they can be written in any language.
package plugin

// WASMParseRequest is sent to the WASM plugin to parse SQL.
type WASMParseRequest struct {
	SQL string `json:"sql"`
}

// WASMParseResponse contains the parsed statements.
type WASMParseResponse struct {
	Statements []WASMStatement `json:"statements"`
}

// WASMStatement represents a parsed SQL statement.
type WASMStatement struct {
	RawSQL       string `json:"raw_sql"`
	StmtLocation int    `json:"stmt_location"`
	StmtLen      int    `json:"stmt_len"`
	ASTJSON      []byte `json:"ast_json"`
}

// WASMGetCatalogRequest is sent to get the initial catalog.
type WASMGetCatalogRequest struct{}

// WASMGetCatalogResponse contains the initial catalog.
type WASMGetCatalogResponse struct {
	Catalog WASMCatalog `json:"catalog"`
}

// WASMCatalog represents the database catalog.
type WASMCatalog struct {
	DefaultSchema string       `json:"default_schema"`
	Name          string       `json:"name"`
	Comment       string       `json:"comment"`
	Schemas       []WASMSchema `json:"schemas"`
	SearchPath    []string     `json:"search_path"`
}

// WASMSchema represents a database schema.
type WASMSchema struct {
	Name      string         `json:"name"`
	Comment   string         `json:"comment"`
	Tables    []WASMTable    `json:"tables"`
	Enums     []WASMEnum     `json:"enums"`
	Functions []WASMFunction `json:"functions"`
	Types     []WASMType     `json:"types"`
}

// WASMTable represents a database table.
type WASMTable struct {
	Catalog string       `json:"catalog"`
	Schema  string       `json:"schema"`
	Name    string       `json:"name"`
	Columns []WASMColumn `json:"columns"`
	Comment string       `json:"comment"`
}

// WASMColumn represents a table column.
type WASMColumn struct {
	Name       string `json:"name"`
	DataType   string `json:"data_type"`
	NotNull    bool   `json:"not_null"`
	IsArray    bool   `json:"is_array"`
	ArrayDims  int    `json:"array_dims"`
	Comment    string `json:"comment"`
	Length     int    `json:"length"`
	IsUnsigned bool   `json:"is_unsigned"`
}

// WASMEnum represents an enum type.
type WASMEnum struct {
	Schema  string   `json:"schema"`
	Name    string   `json:"name"`
	Values  []string `json:"values"`
	Comment string   `json:"comment"`
}

// WASMFunction represents a database function.
type WASMFunction struct {
	Schema     string            `json:"schema"`
	Name       string            `json:"name"`
	Args       []WASMFunctionArg `json:"args"`
	ReturnType WASMDataType      `json:"return_type"`
	Comment    string            `json:"comment"`
}

// WASMFunctionArg represents a function argument.
type WASMFunctionArg struct {
	Name       string       `json:"name"`
	Type       WASMDataType `json:"type"`
	HasDefault bool         `json:"has_default"`
}

// WASMDataType represents a SQL data type.
type WASMDataType struct {
	Catalog string `json:"catalog"`
	Schema  string `json:"schema"`
	Name    string `json:"name"`
}

// WASMType represents a composite or custom type.
type WASMType struct {
	Schema  string `json:"schema"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

// WASMIsReservedKeywordRequest is sent to check if a keyword is reserved.
type WASMIsReservedKeywordRequest struct {
	Keyword string `json:"keyword"`
}

// WASMIsReservedKeywordResponse contains the result.
type WASMIsReservedKeywordResponse struct {
	IsReserved bool `json:"is_reserved"`
}

// WASMGetCommentSyntaxRequest is sent to get supported comment syntax.
type WASMGetCommentSyntaxRequest struct{}

// WASMGetCommentSyntaxResponse contains supported comment syntax.
type WASMGetCommentSyntaxResponse struct {
	Dash      bool `json:"dash"`
	SlashStar bool `json:"slash_star"`
	Hash      bool `json:"hash"`
}

// WASMGetDialectRequest is sent to get dialect information.
type WASMGetDialectRequest struct{}

// WASMGetDialectResponse contains dialect information.
type WASMGetDialectResponse struct {
	QuoteChar   string `json:"quote_char"`
	ParamStyle  string `json:"param_style"`
	ParamPrefix string `json:"param_prefix"`
	CastSyntax  string `json:"cast_syntax"`
}
