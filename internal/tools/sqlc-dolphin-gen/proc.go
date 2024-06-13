package main

import (
	"context"
	"database/sql"
	"strings"
)

const catalogFuncs = `
SELECT
    r.ROUTINE_SCHEMA AS "schema",
    r.SPECIFIC_NAME AS "name",
    r.ROUTINE_TYPE AS "routine_type", -- FUNCTION or PROCEDURE
    r.DATA_TYPE AS "return_type",
    p.ORDINAL_POSITION AS "param_position",
    p.PARAMETER_MODE AS "param_direction", -- IN or OUT
    p.PARAMETER_NAME AS "param_name",
    p.DATA_TYPE AS "param_type"
FROM information_schema.ROUTINES r
LEFT OUTER JOIN information_schema.PARAMETERS p ON (r.ROUTINE_SCHEMA = p.SPECIFIC_SCHEMA AND r.SPECIFIC_NAME = p.SPECIFIC_NAME)

WHERE ROUTINE_SCHEMA = ?
AND p.PARAMETER_MODE IS NOT NULL -- this looks like parameter 0 of a function (not procedure) is a copy of the return type from ROUTINES
-- simply order all columns to keep subsequent runs stable
ORDER BY 1,2,3,4,5,6,7,8;
`

type routineRow struct {
	Schema         string
	Name           string
	Type           string
	ReturnType     *string
	ParamPosition  *int
	ParamDirection *string
	ParamName      *string
	ParamType      *string
}

type Proc struct {
	Name       string
	ReturnType string
	ArgTypes   []string
	ArgNames   []string
	ArgModes   []string
}

func (p *Proc) Args() []Arg {
	var args []Arg

	for i, argType := range p.ArgTypes {
		mode := "i"
		name := ""
		if i < len(p.ArgModes) {
			mode = p.ArgModes[i]
		}
		if i < len(p.ArgNames) {
			name = p.ArgNames[i]
		}

		args = append(args, Arg{
			Name: name,
			Type: argType,
			Mode: mode,
		})
	}

	return args
}

func (a *Arg) GoMode() string {
	switch a.Mode {
	case "", "IN":
		return "ast.FuncParamIn"
	case "OUT":
		return "ast.FuncParamOut"
	}

	return ""
}

type Arg struct {
	Name string
	Mode string
	Type string
}

func scanProcs(rows *sql.Rows) ([]Proc, error) {
	defer rows.Close()

	var procs []Proc
	var prevProc *Proc

	for rows.Next() {
		// Read all rows for a given proc, make up the function
		var rowData routineRow
		err := rows.Scan(
			&rowData.Schema,
			&rowData.Name,
			&rowData.Type,
			&rowData.ReturnType,
			&rowData.ParamPosition,
			&rowData.ParamDirection,
			&rowData.ParamName,
			&rowData.ParamType,
		)
		if err != nil {
			return nil, err
		}

		if prevProc == nil || rowData.Name != prevProc.Name {
			p := Proc{
				Name: strings.ToLower(rowData.Name),
			}
			if rowData.ReturnType != nil && *rowData.ReturnType != "" {
				p.ReturnType = *rowData.ReturnType
			}

			procs = append(procs, p)
			prevProc = &procs[len(procs)-1]
		}

		if rowData.ParamPosition != nil {
			prevProc.ArgNames = append(prevProc.ArgNames, *rowData.ParamName)
			prevProc.ArgTypes = append(prevProc.ArgTypes, *rowData.ParamType)
			prevProc.ArgModes = append(prevProc.ArgModes, *rowData.ParamDirection)
		}
	}

	return procs, rows.Err()
}

func readProcs(ctx context.Context, conn *sql.DB, schemaName string) ([]Proc, error) {
	rows, err := conn.Query(catalogFuncs, schemaName)
	if err != nil {
		return nil, err
	}

	return scanProcs(rows)
}
