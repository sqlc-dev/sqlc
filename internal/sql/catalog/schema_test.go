package catalog

import (
	"strings"
	"testing"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

func TestSchemaGetFunc(t *testing.T) {

	testCases := []struct {
		name    string
		rel     *ast.FuncName
		tns     []*ast.TypeName
		outStub func(t *testing.T, schema *Schema, schemaFunc *Function, funcIndex int, err error)
	}{
		{
			name: "FunctionFound",
			rel: &ast.FuncName{
				Name:   "func1",
				Schema: "schema1",
			},
			tns: []*ast.TypeName{
				{
					Name:    "type1",
					Catalog: "catalog1",
					Schema:  "schema1",
				},
			},
			outStub: func(t *testing.T, schema *Schema, schemaFunc *Function, funcIndex int, err error) {

				if funcIndex < 0 {
					t.Errorf("invalid funcIndex want greater than 0 got: %d", funcIndex)
				}

				if schema.Funcs[funcIndex].Name != schemaFunc.Name {
					t.Errorf("schema function want: %s, got %s", schema.Funcs[funcIndex].Name, schemaFunc.Name)
				}

				if err != nil {
					t.Errorf("err should be nil got %v", err)
				}
			},
		},
		{
			name: "FunctionAgumentMismatch",
			rel: &ast.FuncName{
				Name:   "func1",
				Schema: "schema1",
			},
			tns: []*ast.TypeName{
				{
					Name:    "type1",
					Catalog: "catalog1",
					Schema:  "schema1",
				},
				{
					Name:    "type2",
					Catalog: "catalog1",
					Schema:  "schema1",
				},
			},
			outStub: func(t *testing.T, schema *Schema, schemaFunc *Function, funcIndex int, err error) {

				if funcIndex > -1 {
					t.Errorf("invalid funcIndex want less than 0 got: %d", funcIndex)
				}

				if schemaFunc != nil {
					t.Error("schema function should be nil")
				}

				if err == nil {
					t.Error("err should not be nil")
				}
			},
		},
		{
			name: "FunctionNameMismatch",
			rel: &ast.FuncName{
				Name:   "func",
				Schema: "schema1",
			},
			tns: []*ast.TypeName{
				{
					Name:    "type1",
					Catalog: "catalog1",
					Schema:  "schema1",
				},
			},
			outStub: func(t *testing.T, schema *Schema, schemaFunc *Function, funcIndex int, err error) {

				if funcIndex > -1 {
					t.Errorf("invalid funcIndex want less than 0 got: %d", funcIndex)
				}

				if schemaFunc != nil {
					t.Error("schema function should be nil")
				}

				if err == nil {
					t.Error("err should not be nil")
				}
			},
		},
		{
			name: "SchemaNameMismatch",
			rel: &ast.FuncName{
				Name:   "func1",
				Schema: "schema2",
			},
			tns: []*ast.TypeName{
				{
					Name:    "type1",
					Catalog: "catalog1",
					Schema:  "schema2",
				},
			},
			outStub: func(t *testing.T, schema *Schema, schemaFunc *Function, funcIndex int, err error) {

				if funcIndex > -1 {
					t.Errorf("invalid funcIndex want less than 0 got: %d", funcIndex)
				}

				if schemaFunc != nil {
					t.Error("schema function should be nil")
				}

				if err == nil {
					t.Error("err should not be nil")
				}
			},
		},
	}

	schema := &Schema{
		Name: "schema1",
		Funcs: []*Function{
			{
				Name: "func1",
				Args: []*Argument{
					{
						Name: "arg1",
						Type: &ast.TypeName{
							Name:    "type1",
							Catalog: "catalog1",
							Schema:  "schema1",
						},
						Mode: ast.FuncParamIn,
					},
				},
			},
			{
				Name: "func2",
				Args: []*Argument{
					{
						Name: "arg1",
						Type: &ast.TypeName{
							Name:    "type1",
							Catalog: "catalog1",
							Schema:  "schema1",
						},
						Mode: ast.FuncParamIn,
					},
				},
			},
		},
	}

	for i := range testCases {

		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			outFunc, outIndex, err := schema.getFunc(tc.rel, tc.tns)
			tc.outStub(t, schema, outFunc, outIndex, err)
		})
	}
}

func TestSchemaGetFuncByName(t *testing.T) {

	schema := &Schema{
		Name: "schema1",
		Funcs: []*Function{
			{
				Name: "func1",
				Args: []*Argument{
					{
						Name: "arg1",
						Type: &ast.TypeName{
							Name:    "type1",
							Catalog: "catalog1",
							Schema:  "schema1",
						},
						Mode: ast.FuncParamIn,
					},
				},
			},
			{
				Name: "func2",
				Args: []*Argument{
					{
						Name: "arg1",
						Type: &ast.TypeName{
							Name:    "type1",
							Catalog: "catalog1",
							Schema:  "schema1",
						},
						Mode: ast.FuncParamIn,
					},
				},
			},
		},
	}

	testCases := []struct {
		name    string
		schema  *Schema
		rel     *ast.FuncName
		outStub func(t *testing.T, schema *Schema, schemaFunc *Function, funcIndex int, err error)
	}{
		{
			name:   "FunctionFound",
			schema: schema,
			rel: &ast.FuncName{
				Name:   "func1",
				Schema: "schema1",
			},
			outStub: func(t *testing.T, schema *Schema, schemaFunc *Function, funcIndex int, err error) {

				if funcIndex < 0 {
					t.Errorf("invalid funcIndex want greater than 0 got: %d", funcIndex)
				}

				if schema.Funcs[funcIndex].Name != schemaFunc.Name {
					t.Errorf("schema function want: %s, got %s", schema.Funcs[funcIndex].Name, schemaFunc.Name)
				}

				if err != nil {
					t.Errorf("err should be nil got %v", err)
				}
			},
		},
		{
			name: "FunctionNotUnique",
			schema: func() *Schema {

				// Create a new schema

				newSchema := Schema{}
				newSchema = *schema

				// Add a duplicate func for test
				newSchema.Funcs = append(newSchema.Funcs, &Function{
					Name: "func1",
					Args: []*Argument{
						{
							Name: "arg1",
							Type: &ast.TypeName{
								Name:    "type1",
								Catalog: "catalog1",
								Schema:  "schema1",
							},
							Mode: ast.FuncParamIn,
						},
					},
				})

				return &newSchema
			}(),
			rel: &ast.FuncName{
				Name:   "func1",
				Schema: "schema1",
			},
			outStub: func(t *testing.T, schema *Schema, schemaFunc *Function, funcIndex int, err error) {

				if funcIndex > -1 {
					t.Errorf("invalid funcIndex want less than 0 got: %d", funcIndex)
				}

				if schemaFunc != nil {
					t.Error("schema function should be nil")
				}

				if err == nil {
					t.Error("err should not be nil")
				}
			},
		},
		{
			name:   "FunctionNotFound",
			schema: schema,
			rel: &ast.FuncName{
				Name:   "func",
				Schema: "schema1",
			},
			outStub: func(t *testing.T, schema *Schema, schemaFunc *Function, funcIndex int, err error) {

				if funcIndex > -1 {
					t.Errorf("invalid funcIndex want less than 0 got: %d", funcIndex)
				}

				if schemaFunc != nil {
					t.Error("schema function should be nil")
				}

				if err == nil {
					t.Error("err should not be nil")
				}
			},
		},
	}

	for i := range testCases {

		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			outFunc, outIndex, err := tc.schema.getFuncByName(tc.rel)
			tc.outStub(t, tc.schema, outFunc, outIndex, err)
		})
	}
}

func TestSchemaGetTable(t *testing.T) {

	schema := &Schema{
		Name: "schema1",
		Tables: []*Table{
			{
				Rel: &ast.TableName{
					Name: "table1",
				},
			},
			{
				Rel: &ast.TableName{
					Name: "table2",
				},
			},
		},
	}

	testCases := []struct {
		name    string
		rel     *ast.TableName
		outStub func(t *testing.T, schema *Schema, schemaFunc *Table, tableIndex int, err error)
	}{
		{
			name: "TableFound",
			rel: &ast.TableName{
				Name: "table2",
			},
			outStub: func(t *testing.T, schema *Schema, schemaTable *Table, tableIndex int, err error) {

				if tableIndex < 0 {
					t.Errorf("invalid tableIndex want greater than 0 got: %d", tableIndex)
				}

				if schema.Tables[tableIndex].Rel.Name != schemaTable.Rel.Name {
					t.Errorf("schema table want: %s, got %s", schema.Tables[tableIndex].Rel.Name, schemaTable.Rel.Name)
				}

				if err != nil {
					t.Errorf("err should be nil got %v", err)
				}
			},
		},
		{
			name: "TableNotFound",
			rel: &ast.TableName{
				Name: "table",
			},
			outStub: func(t *testing.T, schema *Schema, schemaTable *Table, tableIndex int, err error) {

				if tableIndex > -1 {
					t.Errorf("invalid tableIndex want less than 0 got: %d", tableIndex)
				}

				if schemaTable != nil {
					t.Error("schema table should be nil")
				}

				if err == nil {
					t.Error("err should not be nil")
				}
			},
		},
	}

	for i := range testCases {

		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			outTable, outIndex, err := schema.getTable(tc.rel)
			tc.outStub(t, schema, outTable, outIndex, err)
		})
	}
}

func TestSchemaGetType(t *testing.T) {

	schema := &Schema{
		Name: "schema1",
		Types: []Type{
			&Enum{
				Name: "emum1",
			},
			&Enum{
				Name: "emum2",
			},
			&CompositeType{
				Name: "compositeType1",
			},
		},
	}

	testCases := []struct {
		name    string
		rel     *ast.TypeName
		outStub func(t *testing.T, schema *Schema, schemaEnum Type, typeIndex int, err error)
	}{
		{
			name: "TypeFound",
			rel: &ast.TypeName{
				Name: "emum1",
			},
			outStub: func(t *testing.T, schema *Schema, schemaType Type, typeIndex int, err error) {

				if typeIndex < 0 {
					t.Errorf("invalid typeIndex want greater than 0 got: %d", typeIndex)
				}

				if schemaType == nil {
					t.Error("schema type should not be nil")
				}

				if err != nil {
					t.Errorf("err should be nil got %v", err)
				}
			},
		},
		{
			name: "TypeNotFound",
			rel: &ast.TypeName{
				Name: "enum",
			},
			outStub: func(t *testing.T, schema *Schema, schemaType Type, typeIndex int, err error) {

				if typeIndex > -1 {
					t.Errorf("invalid typeIndex want greater than 0 got: %d", typeIndex)
				}

				if schemaType != nil {
					t.Error("schema type should not be nil")
				}

				if err == nil {
					t.Error("err should not be nil")
				}
			},
		},
		{
			name: "TypeInvalid",
			rel: &ast.TypeName{
				Name: "compositeType1",
			},
			outStub: func(t *testing.T, schema *Schema, schemaType Type, typeIndex int, err error) {

				if typeIndex > -1 {
					t.Errorf("invalid typeIndex want greater than 0 got: %d", typeIndex)
				}

				if schemaType != nil {
					t.Error("schema type should not be nil")
				}

				if err == nil {
					t.Error("err should not be nil")
				}
			},
		},
	}

	for i := range testCases {

		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			outType, outIndex, err := schema.getType(tc.rel)
			tc.outStub(t, schema, outType, outIndex, err)
		})
	}
}

func TestGetSchema(t *testing.T) {

	testCases := []struct {
		name       string
		schemaName string
		outStub    func(t *testing.T, outSchema *Schema, err error)
	}{
		{
			name:       "Schema Found",
			schemaName: "default",
			outStub: func(t *testing.T, schema *Schema, err error) {

				if schema.Name != "default" {
					t.Errorf("schema name want default got %s", schema.Name)
				}

				if err != nil {
					t.Errorf("err should be nil got %v", err)
				}
			},
		},
		{
			name:       "Schema Not Found",
			schemaName: "wrongSchema",
			outStub: func(t *testing.T, schema *Schema, err error) {

				if schema != nil {
					t.Error("should be nil")
				}

				if err == nil {
					t.Error("err should not be nil")
				}
			},
		},
	}

	newCatalog := New("default")

	for i := range testCases {

		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			outSchema, err := newCatalog.getSchema(tc.schemaName)
			tc.outStub(t, outSchema, err)
		})
	}
}

func TestCreteSchema(t *testing.T) {

	testCases := []struct {
		name    string
		stmt    *ast.CreateSchemaStmt
		outStub func(t *testing.T, newCatalog *Catalog, err error)
	}{
		{
			name: "No Schema Name",
			stmt: func() *ast.CreateSchemaStmt {
				return &ast.CreateSchemaStmt{
					Name: nil,
				}
			}(),
			outStub: func(t *testing.T, newCatalog *Catalog, err error) {

				if !strings.Contains(err.Error(), "empty name") {
					t.Errorf("should contain phrase: empty name, got: %s", err.Error())
				}

				if len(newCatalog.Schemas) != 1 {
					t.Errorf("catalog schema length: want 1, got: %d", len(newCatalog.Schemas))
				}
			},
		},
		{
			name: "Schema Name Exists",
			stmt: func() *ast.CreateSchemaStmt {

				schemaName := "default"

				return &ast.CreateSchemaStmt{
					Name: &schemaName,
				}
			}(),
			outStub: func(t *testing.T, newCatalog *Catalog, err error) {

				if err == nil {
					t.Error("should not be nil")
				}

				if len(newCatalog.Schemas) != 1 {
					t.Errorf("catalog schema length: want 1, got: %d", len(newCatalog.Schemas))
				}
			},
		},
		{
			name: "Schema Created",
			stmt: func() *ast.CreateSchemaStmt {

				schemaName := "new_schema"

				return &ast.CreateSchemaStmt{
					Name: &schemaName,
				}
			}(),
			outStub: func(t *testing.T, newCatalog *Catalog, err error) {

				if err != nil {
					t.Errorf("should be nil got: %v", err)
				}

				if len(newCatalog.Schemas) != 2 {
					t.Errorf("catalog schema length: want 2, got: %d", len(newCatalog.Schemas))
				}
			},
		},
	}

	newCatalog := New("default")

	for i := range testCases {

		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			err := newCatalog.createSchema(tc.stmt)
			tc.outStub(t, newCatalog, err)
		})
	}
}

func TestDropSchema(t *testing.T) {

	testCases := []struct {
		name    string
		stmt    *ast.DropSchemaStmt
		outStub func(t *testing.T, newCatalog *Catalog, err error)
	}{
		{
			name: "NoSchemaProvided",
			stmt: &ast.DropSchemaStmt{
				Schemas: nil,
			},
			outStub: func(t *testing.T, newCatalog *Catalog, err error) {

				if err != nil {
					t.Errorf("err should be nil, got: %v", err)
				}

				if len(newCatalog.Schemas) != 5 {
					t.Errorf("catalog schema length: want 5, got: %d", len(newCatalog.Schemas))
				}
			},
		},
		{
			name: "DeleteOneSchema",
			stmt: &ast.DropSchemaStmt{
				Schemas: []*ast.String{
					{
						Str: "schema1",
					},
				},
			},
			outStub: func(t *testing.T, newCatalog *Catalog, err error) {

				if err != nil {
					t.Errorf("err should be nil, got: %v", err)
				}

				if len(newCatalog.Schemas) != 4 {
					t.Errorf("catalog schema length: want 4, got: %d", len(newCatalog.Schemas))
				}
			},
		},
		{
			name: "DeleteMultipleSchema",
			stmt: &ast.DropSchemaStmt{
				Schemas: []*ast.String{
					{
						Str: "schema1",
					},
					{
						Str: "schema3",
					},
				},
			},
			outStub: func(t *testing.T, newCatalog *Catalog, err error) {

				if err != nil {
					t.Errorf("err should be nil, got: %v", err)
				}

				if len(newCatalog.Schemas) != 3 {
					t.Errorf("catalog schema length: want 3, got: %d", len(newCatalog.Schemas))
				}
			},
		},
		{
			name: "AllowedMissingSchema",
			stmt: &ast.DropSchemaStmt{
				Schemas: []*ast.String{
					{
						Str: "schema10",
					},
				},
				MissingOk: true,
			},
			outStub: func(t *testing.T, newCatalog *Catalog, err error) {

				if err != nil {
					t.Errorf("err should be nil, got: %v", err)
				}

				if len(newCatalog.Schemas) != 5 {
					t.Errorf("catalog schema length: want 5, got: %d", len(newCatalog.Schemas))
				}
			},
		},
		{
			name: "SchemaNotFound",
			stmt: &ast.DropSchemaStmt{
				Schemas: []*ast.String{
					{
						Str: "missing_schema",
					},
				},
			},
			outStub: func(t *testing.T, newCatalog *Catalog, err error) {

				if err == nil {
					t.Error("err should not nil")
				}

				if len(newCatalog.Schemas) != 5 {
					t.Errorf("catalog schema length: want 5, got: %d", len(newCatalog.Schemas))
				}
			},
		},
	}

	for i := range testCases {

		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			newCatalog := New("default")

			newCatalog.Schemas = append(
				newCatalog.Schemas,
				&Schema{Name: "schema1"},
				&Schema{Name: "schema2"},
				&Schema{Name: "schema3"},
				&Schema{Name: "schema4"})

			err := newCatalog.dropSchema(tc.stmt)
			tc.outStub(t, newCatalog, err)
		})
	}
}
