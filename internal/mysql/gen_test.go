package mysql

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kyleconroy/sqlc/internal/dinosql"
	"github.com/kyleconroy/sqlc/internal/pg"
	"vitess.io/vitess/go/vt/sqlparser"
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

	for _, tc := range tcase {
		enumName := enumNameFromColDef(&tc.input, mockSettings)
		if diff := cmp.Diff(enumName, tc.output); diff != "" {
			t.Errorf(diff)
		}
	}
}

func TestEnums(t *testing.T) {
	tcase := [...]struct {
		input  Result
		output []dinosql.GoEnum
	}{
		{
			input: Result{Schema: mockSchema},
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
	for _, tc := range tcase {
		enums := tc.input.Enums(mockSettings)
		if diff := cmp.Diff(enums, tc.output); diff != "" {
			t.Errorf(diff)
		}
	}
}

func TestStructs(t *testing.T) {
	tcase := [...]struct {
		input  Result
		output []dinosql.GoStruct
	}{
		{
			input: Result{Schema: mockSchema},
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
		structs := tc.input.Structs(mockSettings)
		if diff := cmp.Diff(structs, tc.output); diff != "" {
			t.Errorf(diff)
		}
	}
}
