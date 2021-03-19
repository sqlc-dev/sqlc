// +build !windows

package postgresql

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	nodes "github.com/pganalyze/pg_query_go/v2"

	"github.com/kyleconroy/sqlc/internal/metadata"
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

func stringSlice(list *nodes.List) []string {
	items := []string{}
	for _, item := range list.Items {
		if n, ok := item.Node.(*nodes.Node_String_); ok {
			items = append(items, n.String_.Str)
		}
	}
	return items
}

func stringSliceFromNodes(s []*nodes.Node) []string {
	var items []string
	for _, item := range s {
		if n, ok := item.Node.(*nodes.Node_String_); ok {
			items = append(items, n.String_.Str)
		}
	}
	return items
}

type relation struct {
	Catalog string
	Schema  string
	Name    string
}

func (r relation) TableName() *ast.TableName {
	return &ast.TableName{
		Catalog: r.Catalog,
		Schema:  r.Schema,
		Name:    r.Name,
	}
}

func (r relation) TypeName() *ast.TypeName {
	return &ast.TypeName{
		Catalog: r.Catalog,
		Schema:  r.Schema,
		Name:    r.Name,
	}
}

func (r relation) FuncName() *ast.FuncName {
	return &ast.FuncName{
		Catalog: r.Catalog,
		Schema:  r.Schema,
		Name:    r.Name,
	}
}

func parseFuncParamMode(m nodes.FunctionParameterMode) (ast.FuncParamMode, error) {
	switch m {
	case 'i':
		return ast.FuncParamIn, nil
	case 'o':
		return ast.FuncParamOut, nil
	case 'b':
		return ast.FuncParamInOut, nil
	case 'v':
		return ast.FuncParamVariadic, nil
	case 't':
		return ast.FuncParamTable, nil
	default:
		return -1, fmt.Errorf("parse func param: invalid mode %v", m)
	}
}

func parseRelationFromNodes(list []*nodes.Node) (*relation, error) {
	parts := stringSliceFromNodes(list)
	switch len(parts) {
	case 1:
		return &relation{
			Name: parts[0],
		}, nil
	case 2:
		return &relation{
			Schema: parts[0],
			Name:   parts[1],
		}, nil
	case 3:
		return &relation{
			Catalog: parts[0],
			Schema:  parts[1],
			Name:    parts[2],
		}, nil
	default:
		return nil, fmt.Errorf("invalid name: %s", joinNodes(list, "."))
	}
}

func parseRelationFromRangeVar(rv *nodes.RangeVar) *relation {
	return &relation{
		Catalog: rv.Catalogname,
		Schema:  rv.Schemaname,
		Name:    rv.Relname,
	}
}

func parseColName(node *nodes.Node) (*ast.ColumnRef, *ast.TableName, error) {
	switch n := node.Node.(type) {
	case *nodes.Node_List:
		parts := stringSlice(n.List)
		var tbl *ast.TableName
		var ref *ast.ColumnRef
		switch len(parts) {
		case 2:
			tbl = &ast.TableName{Name: parts[0]}
			ref = &ast.ColumnRef{Name: parts[1]}
		case 3:
			tbl = &ast.TableName{Schema: parts[0], Name: parts[1]}
			ref = &ast.ColumnRef{Name: parts[2]}
		case 4:
			tbl = &ast.TableName{Catalog: parts[0], Schema: parts[1], Name: parts[2]}
			ref = &ast.ColumnRef{Name: parts[3]}
		default:
			return nil, nil, fmt.Errorf("column specifier %q is not the proper format, expected '[catalog.][schema.]colname.tablename'", strings.Join(parts, "."))
		}
		return ref, tbl, nil
	default:
		return nil, nil, fmt.Errorf("parseColName: unexpected node type: %T", n)
	}
}

func join(list *nodes.List, sep string) string {
	return strings.Join(stringSlice(list), sep)
}

func joinNodes(list []*nodes.Node, sep string) string {
	return strings.Join(stringSliceFromNodes(list), sep)
}

func NewParser() *Parser {
	return &Parser{}
}

type Parser struct {
}

var errSkip = errors.New("skip stmt")

func (p *Parser) Parse(r io.Reader) ([]ast.Statement, error) {
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	tree, err := nodes.Parse(string(contents))
	if err != nil {
		return nil, err
	}

	var stmts []ast.Statement
	for _, raw := range tree.Stmts {
		n, err := translate(raw.Stmt)
		if err == errSkip {
			continue
		}
		if err != nil {
			return nil, err
		}
		if n == nil {
			return nil, fmt.Errorf("unexpected nil node")
		}
		stmts = append(stmts, ast.Statement{
			Raw: &ast.RawStmt{
				Stmt:         n,
				StmtLocation: int(raw.StmtLocation),
				StmtLen:      int(raw.StmtLen),
			},
		})
	}
	return stmts, nil
}

// https://www.postgresql.org/docs/current/sql-syntax-lexical.html#SQL-SYNTAX-COMMENTS
func (p *Parser) CommentSyntax() metadata.CommentSyntax {
	return metadata.CommentSyntax{
		Dash:      true,
		SlashStar: true,
	}
}

func translate(node *nodes.Node) (ast.Node, error) {
	switch inner := node.Node.(type) {

	case *nodes.Node_AlterEnumStmt:
		n := inner.AlterEnumStmt
		rel, err := parseRelationFromNodes(n.TypeName)
		if err != nil {
			return nil, err
		}
		if n.OldVal != "" {
			return &ast.AlterTypeRenameValueStmt{
				Type:     rel.TypeName(),
				OldValue: &n.OldVal,
				NewValue: &n.NewVal,
			}, nil
		} else {
			return &ast.AlterTypeAddValueStmt{
				Type:               rel.TypeName(),
				NewValue:           &n.NewVal,
				SkipIfNewValExists: n.SkipIfNewValExists,
			}, nil
		}

	case *nodes.Node_AlterObjectSchemaStmt:
		n := inner.AlterObjectSchemaStmt
		switch n.ObjectType {

		case nodes.ObjectType_OBJECT_TABLE:
			rel := parseRelationFromRangeVar(n.Relation)
			return &ast.AlterTableSetSchemaStmt{
				Table:     rel.TableName(),
				NewSchema: &n.Newschema,
			}, nil
		}
		return nil, errSkip

	case *nodes.Node_AlterTableStmt:
		n := inner.AlterTableStmt
		rel := parseRelationFromRangeVar(n.Relation)
		at := &ast.AlterTableStmt{
			Table: rel.TableName(),
			Cmds:  &ast.List{},
		}
		for _, cmd := range n.Cmds {
			switch cmdOneOf := cmd.Node.(type) {
			case *nodes.Node_AlterTableCmd:
				altercmd := cmdOneOf.AlterTableCmd
				item := &ast.AlterTableCmd{Name: &altercmd.Name, MissingOk: altercmd.MissingOk}

				switch altercmd.Subtype {
				case nodes.AlterTableType_AT_AddColumn:
					d := altercmd.Def.(nodes.ColumnDef)
					rel, err := parseRelationFromNodes(d.TypeName)
					if err != nil {
						return nil, err
					}
					item.Subtype = ast.AT_AddColumn
					item.Def = &ast.ColumnDef{
						Colname:   *d.Colname,
						TypeName:  rel.TypeName(),
						IsNotNull: isNotNull(d),
						IsArray:   isArray(d.TypeName),
					}

				case nodes.AlterTableType_AT_AlterColumnType:
					d := altercmd.Def.(nodes.ColumnDef)
					col := ""
					if cmd.Name != nil {
						col = *altercmd.Name
					} else if d.Colname != nil {
						col = *d.Colname
					} else {
						return nil, fmt.Errorf("unknown name for alter column type")
					}
					rel, err := parseRelationFromNodes(d.TypeName)
					if err != nil {
						return nil, err
					}
					item.Subtype = ast.AT_AlterColumnType
					item.Def = &ast.ColumnDef{
						Colname:   col,
						TypeName:  rel.TypeName(),
						IsNotNull: isNotNull(d),
						IsArray:   isArray(d.TypeName),
					}

				case nodes.AlterTableType_AT_DropColumn:
					item.Subtype = ast.AT_DropColumn

				case nodes.AlterTableType_AT_DropNotNull:
					item.Subtype = ast.AT_DropNotNull

				case nodes.AlterTableType_AT_SetNotNull:
					item.Subtype = ast.AT_SetNotNull

				default:
					continue
				}

				at.Cmds.Items = append(at.Cmds.Items, item)
			}
		}
		return at, nil

	case *nodes.Node_CommentStmt:
		n := inner.CommentStmt
		switch n.Objtype {

		case nodes.ObjectType_OBJECT_COLUMN:
			col, tbl, err := parseColName(n.Object)
			if err != nil {
				return nil, fmt.Errorf("COMMENT ON COLUMN: %w", err)
			}
			return &ast.CommentOnColumnStmt{
				Col:     col,
				Table:   tbl,
				Comment: n.Comment,
			}, nil

		case nodes.ObjectType_OBJECT_SCHEMA:
			o, ok := n.Object.(nodes.String)
			if !ok {
				return nil, fmt.Errorf("COMMENT ON SCHEMA: unexpected node type: %T", n.Object)
			}
			return &ast.CommentOnSchemaStmt{
				Schema:  &ast.String{Str: o.Str},
				Comment: n.Comment,
			}, nil

		case nodes.ObjectType_OBJECT_TABLE:
			name, err := parseTableName(n.Object)
			if err != nil {
				return nil, fmt.Errorf("COMMENT ON TABLE: %w", err)
			}
			return &ast.CommentOnTableStmt{
				Table:   name,
				Comment: n.Comment,
			}, nil

		case nodes.ObjectType_OBJECT_TYPE:
			name, err := parseTypeName(n.Object)
			if err != nil {
				return nil, err
			}
			return &ast.CommentOnTypeStmt{
				Type:    name,
				Comment: n.Comment,
			}, nil

		}
		return nil, errSkip

	case *nodes.Node_CompositeTypeStmt:
		n := inner.CompositeTypeStmt
		name, err := parseTypeName(n.Typevar)
		if err != nil {
			return nil, err
		}
		return &ast.CompositeTypeStmt{
			TypeName: name,
		}, nil

	case *nodes.Node_CreateStmt:
		n := inner.CreateStmt
		name, err := parseTableName(*n.Relation)
		if err != nil {
			return nil, err
		}
		create := &ast.CreateTableStmt{
			Name:        name,
			IfNotExists: n.IfNotExists,
		}
		primaryKey := make(map[string]bool)
		for _, elt := range n.TableElts.Items {
			switch n := elt.(type) {
			case nodes.Constraint:
				if n.Contype == nodes.CONSTR_PRIMARY {
					for _, item := range n.Keys.Items {
						primaryKey[item.(nodes.String).Str] = true
					}
				}
			}
		}
		for _, elt := range n.TableElts.Items {
			switch n := elt.(type) {
			case nodes.ColumnDef:
				tn, err := parseTypeName(n.TypeName)
				if err != nil {
					return nil, err
				}
				create.Cols = append(create.Cols, &ast.ColumnDef{
					Colname:   *n.Colname,
					TypeName:  tn,
					IsNotNull: isNotNull(n) || primaryKey[*n.Colname],
					IsArray:   isArray(n.TypeName),
				})
			}
		}
		return create, nil

	case *nodes.Node_CreateEnumStmt:
		n := inner.CreateEnumStmt
		name, err := parseTypeName(n.TypeName)
		if err != nil {
			return nil, err
		}
		stmt := &ast.CreateEnumStmt{
			TypeName: name,
			Vals:     &ast.List{},
		}
		for _, val := range n.Vals.Items {
			switch v := val.(type) {
			case nodes.String:
				stmt.Vals.Items = append(stmt.Vals.Items, &ast.String{
					Str: v.Str,
				})
			}
		}
		return stmt, nil

	case *nodes.Node_CreateFunctionStmt:
		n := inner.CreateFunctionStmt
		fn, err := parseFuncName(n.Funcname)
		if err != nil {
			return nil, err
		}
		rt, err := parseTypeName(n.ReturnType)
		if err != nil {
			return nil, err
		}
		stmt := &ast.CreateFunctionStmt{
			Func:       fn,
			ReturnType: rt,
			Replace:    n.Replace,
			Params:     &ast.List{},
		}
		for _, item := range n.Parameters.Items {
			arg := item.(nodes.FunctionParameter)
			tn, err := parseTypeName(arg.ArgType)
			if err != nil {
				return nil, err
			}
			mode, err := parseFuncParamMode(arg.Mode)
			if err != nil {
				return nil, err
			}
			fp := &ast.FuncParam{
				Name: arg.Name,
				Type: tn,
				Mode: mode,
			}
			if arg.Defexpr != nil {
				fp.DefExpr = &ast.TODO{}
			}
			stmt.Params.Items = append(stmt.Params.Items, fp)
		}
		return stmt, nil

	case *nodes.Node_CreateSchemaStmt:
		n := inner.CreateSchemaStmt
		return &ast.CreateSchemaStmt{
			Name:        n.Schemaname,
			IfNotExists: n.IfNotExists,
		}, nil

	case *nodes.Node_DropStmt:
		n := inner.DropStmt
		switch n.RemoveType {

		case nodes.OBJECT_FUNCTION:
			drop := &ast.DropFunctionStmt{
				MissingOk: n.MissingOk,
			}
			for _, obj := range n.Objects.Items {
				owa, ok := obj.(nodes.ObjectWithArgs)
				if !ok {
					return nil, fmt.Errorf("nodes.DropStmt: FUNCTION: unknown type in objects list: %T", obj)
				}
				fn, err := parseFuncName(owa.Objname)
				if err != nil {
					return nil, fmt.Errorf("nodes.DropStmt: FUNCTION: %w", err)
				}
				args := make([]*ast.TypeName, len(owa.Objargs.Items))
				for i, objarg := range owa.Objargs.Items {
					tn, ok := objarg.(nodes.TypeName)
					if !ok {
						return nil, fmt.Errorf("nodes.DropStmt: FUNCTION: unknown type in objargs list: %T", objarg)
					}
					at, err := parseTypeName(tn)
					if err != nil {
						return nil, fmt.Errorf("nodes.DropStmt: FUNCTION: %w", err)
					}
					args[i] = at
				}
				drop.Funcs = append(drop.Funcs, &ast.FuncSpec{
					Name:    fn,
					Args:    args,
					HasArgs: !owa.ArgsUnspecified,
				})
			}
			return drop, nil

		case nodes.ObjectType_OBJECT_SCHEMA:
			drop := &ast.DropSchemaStmt{
				MissingOk: n.MissingOk,
			}
			for _, obj := range n.Objects.Items {
				val, ok := obj.(nodes.String)
				if !ok {
					return nil, fmt.Errorf("nodes.DropStmt: SCHEMA: unknown type in objects list: %T", obj)
				}
				drop.Schemas = append(drop.Schemas, &ast.String{Str: val.Str})
			}
			return drop, nil

		case nodes.ObjectType_OBJECT_TABLE:
			drop := &ast.DropTableStmt{
				IfExists: n.MissingOk,
			}
			for _, obj := range n.Objects.Items {
				name, err := parseTableName(obj)
				if err != nil {
					return nil, fmt.Errorf("nodes.DropStmt: TABLE: %w", err)
				}
				drop.Tables = append(drop.Tables, name)
			}
			return drop, nil

		case nodes.ObjectType_OBJECT_TYPE:
			drop := &ast.DropTypeStmt{
				IfExists: n.MissingOk,
			}
			for _, obj := range n.Objects.Items {
				name, err := parseTypeName(obj)
				if err != nil {
					return nil, fmt.Errorf("nodes.DropStmt: TYPE: %w", err)
				}
				drop.Types = append(drop.Types, name)
			}
			return drop, nil

		}
		return nil, errSkip

	case *nodes.Node_RenameStmt:
		n := inner.RenameStmt
		switch n.RenameType {

		case nodes.ObjectType_OBJECT_COLUMN:
			tbl, err := parseTableName(*n.Relation)
			if err != nil {
				return nil, fmt.Errorf("nodes.RenameType: COLUMN: %w", err)
			}
			return &ast.RenameColumnStmt{
				Table:   tbl,
				Col:     &ast.ColumnRef{Name: *n.Subname},
				NewName: n.Newname,
			}, nil

		case nodes.ObjectType_OBJECT_TABLE:
			tbl, err := parseTableName(*n.Relation)
			if err != nil {
				return nil, fmt.Errorf("nodes.RenameType: TABLE: %w", err)
			}
			return &ast.RenameTableStmt{
				Table:   tbl,
				NewName: n.Newname,
			}, nil

		}
		return nil, errSkip

	default:
		return convert(n)
	}
}
