package golang

import (
	"strings"

	"github.com/kyleconroy/sqlc/internal/plugin"
)

type Field struct {
	Name    string // CamelCased name for Go
	DBName  string // Name as used in the DB
	Type    string
	Tags    map[string]string
	Comment string
	Column  *plugin.Column
	// EmbedFields contains the embedded fields that require scanning.
	EmbedFields []string
}

func (gf Field) Tag() string {
	return TagsToString(gf.Tags)
}

func (gf Field) HasSqlcSlice() bool {
	return gf.Column.IsSqlcSlice
}

func toLowerCase(str string) string {
	if str == "" {
		return ""
	}

	return strings.ToLower(str[:1]) + str[1:]
}
