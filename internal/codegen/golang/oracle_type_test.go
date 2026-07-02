package golang

import (
	"testing"

	"github.com/sqlc-dev/sqlc/internal/codegen/golang/opts"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

func TestOracleType(t *testing.T) {
	req := &plugin.GenerateRequest{Settings: &plugin.Settings{Engine: "oracle"}}

	cases := []struct {
		name     string
		dataType string
		notNull  bool
		emitPtr  bool
		want     string
	}{
		{"number nullable", "number", false, false, "sql.NullFloat64"},
		{"number not null", "number", true, false, "float64"},
		{"number ptr", "number", false, true, "*float64"},
		{"integer not null", "integer", true, false, "int64"},
		{"int nullable", "int", false, false, "sql.NullInt64"},
		{"varchar2 not null", "varchar2", true, false, "string"},
		{"varchar2(100) not null", "varchar2(100)", true, false, "string"},
		{"nvarchar2 nullable", "nvarchar2", false, false, "sql.NullString"},
		{"char not null", "char", true, false, "string"},
		{"clob not null", "clob", true, false, "string"},
		{"date not null", "date", true, false, "time.Time"},
		{"timestamp nullable", "timestamp", false, false, "sql.NullTime"},
		{"timestamp(6) not null", "timestamp(6)", true, false, "time.Time"},
		{"raw", "raw", true, false, "[]byte"},
		{"blob", "blob", true, false, "[]byte"},
		{"binary_double not null", "binary_double", true, false, "float64"},
		{"boolean not null", "boolean", true, false, "bool"},
		{"boolean nullable", "boolean", false, false, "sql.NullBool"},
		{"interval", "interval", true, false, "string"},
		{"rowid", "rowid", true, false, "string"},
		{"unknown", "sdo_geometry_xyz", true, false, "interface{}"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			options := &opts.Options{EmitPointersForNullTypes: tc.emitPtr}
			col := &plugin.Column{
				Type:    &plugin.Identifier{Name: tc.dataType},
				NotNull: tc.notNull,
			}
			got := oracleType(req, options, col)
			if got != tc.want {
				t.Errorf("oracleType(%q, notNull=%v, emitPtr=%v) = %q, want %q",
					tc.dataType, tc.notNull, tc.emitPtr, got, tc.want)
			}
		})
	}
}
