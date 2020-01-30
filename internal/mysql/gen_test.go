package mysql

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"vitess.io/vitess/go/vt/sqlparser"

	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/dinosql"
	"github.com/kyleconroy/sqlc/internal/pg"
)

func TestArgName(t *testing.T) {
	tcase := [...]struct {
		input  string
		output string
	}{
		{
			input:  "get_users",
			output: "getUsers",
		},
		{
			input:  "get_users_by_id",
			output: "getUsersByID",
		},
		{
			input:  "get_all_",
			output: "getAll",
		},
	}

	for _, tc := range tcase {
		name := argName(tc.input)
		if diff := cmp.Diff(name, tc.output); diff != "" {
			t.Errorf(diff)
		}
	}
}
func TestEnumName(t *testing.T) {
	tcase := [...]struct {
		input  sqlparser.ColumnDefinition
		output string
	}{
		{
			input:  sqlparser.ColumnDefinition{Name: sqlparser.NewColIdent("first_name")},
			output: "FirstNameType",
		},
		{
			input:  sqlparser.ColumnDefinition{Name: sqlparser.NewColIdent("user_id")},
			output: "UserIDType",
		},
		{
			input:  sqlparser.ColumnDefinition{Name: sqlparser.NewColIdent("last_name")},
			output: "LastNameType",
		},
	}

	generator := PackageGenerator{mockSchema, config.CombinedSettings{}, ""}
	for _, tc := range tcase {
		enumName := generator.enumNameFromColDef(&tc.input)
		if diff := cmp.Diff(enumName, tc.output); diff != "" {
			t.Errorf(diff)
		}
	}
}

func TestEnums(t *testing.T) {
	generator := PackageGenerator{mockSchema, config.CombinedSettings{}, ""}
	tcase := [...]struct {
		input  Result
		output []dinosql.GoEnum
	}{
		{
			input: Result{PackageGenerator: generator},
			output: []dinosql.GoEnum{
				{
					Name: "JobStatusType",
					Constants: []dinosql.GoConstant{
						{Name: "applied", Type: "JobStatusType", Value: "applied"},
						{Name: "pending", Type: "JobStatusType", Value: "pending"},
						{Name: "accepted", Type: "JobStatusType", Value: "accepted"},
						{Name: "rejected", Type: "JobStatusType", Value: "rejected"},
					},
				},
			},
		},
	}
	settings := config.Combine(config.GenerateSettings{}, config.PackageSettings{})
	for _, tc := range tcase {
		enums := tc.input.Enums(settings)
		if diff := cmp.Diff(enums, tc.output); diff != "" {
			t.Errorf(diff)
		}
	}
}

func TestStructs(t *testing.T) {
	settings := config.Combine(config.GenerateSettings{}, config.PackageSettings{})
	generator := PackageGenerator{mockSchema, settings, "db"}
	tcase := [...]struct {
		input  Result
		output []dinosql.GoStruct
	}{
		{
			input: Result{PackageGenerator: generator},
			output: []dinosql.GoStruct{
				{
					Table: pg.FQN{Catalog: "orders"},
					Name:  "Order",
					Fields: []dinosql.GoField{
						{Name: "ID", Type: "int", Tags: map[string]string{"json:": "id"}},
						{Name: "Price", Type: "float64", Tags: map[string]string{"json:": "price"}},
						{Name: "UserID", Type: "int", Tags: map[string]string{"json:": "user_id"}},
					},
				},
				{
					Table: pg.FQN{Catalog: "users"},
					Name:  "User",
					Fields: []dinosql.GoField{
						{Name: "FirstName", Type: "string", Tags: map[string]string{"json:": "first_name"}},
						{Name: "LastName", Type: "sql.NullString", Tags: map[string]string{"json:": "last_name"}},
						{Name: "ID", Type: "int", Tags: map[string]string{"json:": "id"}},
						{Name: "Age", Type: "int", Tags: map[string]string{"json:": "age"}},
						{Name: "JobStatus", Type: "JobStatusType", Tags: map[string]string{"json:": "job_status"}},
					}},
			},
		},
	}

	for _, tc := range tcase {
		structs := tc.input.Structs(settings)
		if diff := cmp.Diff(structs, tc.output); diff != "" {
			t.Errorf(diff)
		}
	}
}

func TestTypeOverride(t *testing.T) {
	tests := [...]struct {
		overrides      []config.Override
		col            Column
		expectedGoType string
	}{
		{
			overrides: []config.Override{
				{
					DBType:     "uuid",
					GoTypeName: "KSUID", // this is populated by the dinosql.Parse
				},
			},
			col: Column{
				ColumnDefinition: &sqlparser.ColumnDefinition{
					Type: sqlparser.ColumnType{
						Type:    "uuid",
						NotNull: true,
					},
				},
			},
			expectedGoType: "KSUID",
		},
		{
			overrides: []config.Override{
				{
					ColumnName: "user_id", // this is populated by dinosql.Parse
					GoTypeName: "uuid",    // this is populated by dinosql.Parse
				},
			},
			col: Column{
				ColumnDefinition: &sqlparser.ColumnDefinition{
					Name: sqlparser.NewColIdent("user_id"),
					Type: sqlparser.ColumnType{
						Type:    "varchar",
						NotNull: true,
					},
				},
			},
			expectedGoType: "uuid",
		},
	}

	for _, tcase := range tests {
		settings := config.Combine(config.GenerateSettings{}, config.PackageSettings{Overrides: tcase.overrides})
		gen := PackageGenerator{mockSchema, settings, "db"}
		goType := gen.goTypeCol(tcase.col)

		if diff := cmp.Diff(tcase.expectedGoType, goType); diff != "" {
			t.Errorf(diff)
		}
	}
}
