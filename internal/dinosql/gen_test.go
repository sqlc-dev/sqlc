package dinosql

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kyleconroy/sqlc/internal/pg"
)

func TestColumnsToStruct(t *testing.T) {
	cols := []pg.Column{
		{
			Name:     "other",
			DataType: "text",
			NotNull:  true,
		},
		{
			Name:     "count",
			DataType: "bigint",
			NotNull:  true,
		},
		{
			Name:     "count",
			DataType: "bigint",
			NotNull:  true,
		},
		{
			Name:     "tags",
			DataType: "text",
			NotNull:  true,
			IsArray:  true,
		},
		{
			Name:     "byte_seq",
			DataType: "bytea",
			NotNull:  true,
		},
		{
			Name:     "retyped",
			DataType: "text",
			NotNull:  true,
		},
	}

	// all of the columns are on the 'foo' table
	for i := range cols {
		cols[i].Table = pg.FQN{Schema: "public", Rel: "foo"}
	}

	r := Result{}

	// set up column-based override test
	o := Override{
		GoType: "example.com/pkg.CustomType",
		Column: "foo.retyped",
	}
	o.Parse()
	r.packageSettings = PackageSettings{
		Overrides: []Override{o},
	}

	actual := r.columnsToStruct("Foo", cols)
	expected := &GoStruct{
		Name: "Foo",
		Fields: []GoField{
			{Name: "Other", Type: "string", Tags: map[string]string{"json:": "other"}},
			{Name: "Count", Type: "int64", Tags: map[string]string{"json:": "count"}},
			{Name: "Count_2", Type: "int64", Tags: map[string]string{"json:": "count_2"}},
			{Name: "Tags", Type: "[]string", Tags: map[string]string{"json:": "tags"}},
			{Name: "ByteSeq", Type: "[]byte", Tags: map[string]string{"json:": "byte_seq"}},
			{Name: "Retyped", Type: "pkg.CustomType", Tags: map[string]string{"json:": "retyped"}},
		},
	}
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("struct mismatch: \n%s", diff)
	}
}

func TestInnerType(t *testing.T) {
	r := Result{}
	types := map[string]string{
		"timestamptz":     "time.Time",
		"integer":         "int32",
		"int":             "int32",
		"pg_catalog.int4": "int32",
	}
	for k, v := range types {
		dbType := k
		goType := v
		t.Run(k+"-"+v, func(t *testing.T) {
			col := pg.Column{DataType: dbType, NotNull: true}
			if goType != r.goType(col) {
				t.Errorf("expected Go type for %s to be %s, not %s", dbType, goType, r.goType(col))
			}
		})
	}
}

func TestNullInnerType(t *testing.T) {
	r := Result{}
	types := map[string]string{
		"timestamptz":     "sql.NullTime",
		"integer":         "sql.NullInt32",
		"int":             "sql.NullInt32",
		"pg_catalog.int4": "sql.NullInt32",
	}
	for k, v := range types {
		dbType := k
		goType := v
		t.Run(k+"-"+v, func(t *testing.T) {
			col := pg.Column{DataType: dbType, NotNull: false}
			if goType != r.goType(col) {
				t.Errorf("expected Go type for %s to be %s, not %s", dbType, goType, r.goType(col))
			}
		})
	}
}
