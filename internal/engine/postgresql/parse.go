package postgresql

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	pg "github.com/lfittl/pg_query_go"
	nodes "github.com/lfittl/pg_query_go/nodes"

	"github.com/kyleconroy/sqlc/internal/metadata"
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

func stringSlice(list nodes.List) []string {
	items := []string{}
	for _, item := range list.Items {
		if n, ok := item.(nodes.String); ok {
			items = append(items, n.Str)
		}
	}
	return items
}

type relation struct {
	Catalog string
	Schema  string
	Name    string
}

func parseFuncName(node nodes.Node) (*ast.FuncName, error) {
	rel, err := parseRelation(node)
	if err != nil {
		return nil, fmt.Errorf("parse func name: %w", err)
	}
	return &ast.FuncName{
		Catalog: rel.Catalog,
		Schema:  rel.Schema,
		Name:    rel.Name,
	}, nil
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

func parseTypeName(node nodes.Node) (*ast.TypeName, error) {
	rel, err := parseRelation(node)
	if err != nil {
		return nil, fmt.Errorf("parse type name: %w", err)
	}
	return &ast.TypeName{
		Catalog: rel.Catalog,
		Schema:  rel.Schema,
		Name:    rel.Name,
	}, nil
}

func parseTableName(node nodes.Node) (*ast.TableName, error) {
	rel, err := parseRelation(node)
	if err != nil {
		return nil, fmt.Errorf("parse table name: %w", err)
	}
	return &ast.TableName{
		Catalog: rel.Catalog,
		Schema:  rel.Schema,
		Name:    rel.Name,
	}, nil
}

func parseRelation(node nodes.Node) (*relation, error) {
	switch n := node.(type) {

	case nodes.List:
		parts := stringSlice(n)
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
			return nil, fmt.Errorf("invalid name: %s", join(n, "."))
		}

	case nodes.RangeVar:
		name := relation{}
		if n.Catalogname != nil {
			name.Catalog = *n.Catalogname
		}
		if n.Schemaname != nil {
			name.Schema = *n.Schemaname
		}
		if n.Relname != nil {
			name.Name = *n.Relname
		}
		return &name, nil

	case *nodes.RangeVar:
		name := relation{}
		if n.Catalogname != nil {
			name.Catalog = *n.Catalogname
		}
		if n.Schemaname != nil {
			name.Schema = *n.Schemaname
		}
		if n.Relname != nil {
			name.Name = *n.Relname
		}
		return &name, nil

	case nodes.TypeName:
		return parseRelation(n.Names)

	case *nodes.TypeName:
		return parseRelation(n.Names)

	default:
		return nil, fmt.Errorf("unexpected node type: %T", n)
	}
}

func parseColName(node nodes.Node) (*ast.ColumnRef, *ast.TableName, error) {
	switch n := node.(type) {
	case nodes.List:
		parts := stringSlice(n)
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

func join(list nodes.List, sep string) string {
	return strings.Join(stringSlice(list), sep)
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
	tree, err := pg.Parse(string(contents))
	if err != nil {
		return nil, err
	}

	var stmts []ast.Statement
	for _, stmt := range tree.Statements {
		raw, ok := stmt.(nodes.RawStmt)
		if !ok {
			return nil, fmt.Errorf("expected RawStmt; got %T", stmt)
		}
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
				StmtLocation: raw.StmtLocation,
				StmtLen:      raw.StmtLen,
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

func translate(node nodes.Node) (ast.Node, error) {
	switch n := node.(type) {

	case nodes.AlterEnumStmt:
		name, err := parseTypeName(n.TypeName)
		if err != nil {
			return nil, err
		}
		if n.OldVal != nil {
			return &ast.AlterTypeRenameValueStmt{
				Type:     name,
				OldValue: n.OldVal,
				NewValue: n.NewVal,
			}, nil
		} else {
			return &ast.AlterTypeAddValueStmt{
				Type:               name,
				NewValue:           n.NewVal,
				SkipIfNewValExists: n.SkipIfNewValExists,
			}, nil
		}

	case nodes.AlterObjectSchemaStmt:
		switch n.ObjectType {

		case nodes.OBJECT_TABLE:
			tbl, err := parseTableName(*n.Relation)
			if err != nil {
				return nil, err
			}
			return &ast.AlterTableSetSchemaStmt{
				Table:     tbl,
				NewSchema: n.Newschema,
			}, nil
		}
		return nil, errSkip

	case nodes.AlterTableStmt:
		name, err := parseTableName(*n.Relation)
		if err != nil {
			return nil, err
		}
		at := &ast.AlterTableStmt{
			Table: name,
			Cmds:  &ast.List{},
		}
		for _, cmd := range n.Cmds.Items {
			switch cmd := cmd.(type) {
			case nodes.AlterTableCmd:
				item := &ast.AlterTableCmd{Name: cmd.Name, MissingOk: cmd.MissingOk}

				switch cmd.Subtype {
				case nodes.AT_AddColumn:
					d := cmd.Def.(nodes.ColumnDef)
					tn, err := parseTypeName(d.TypeName)
					if err != nil {
						return nil, err
					}
					item.Subtype = ast.AT_AddColumn
					item.Def = &ast.ColumnDef{
						Colname:   *d.Colname,
						TypeName:  tn,
						IsNotNull: isNotNull(d),
						IsArray:   isArray(d.TypeName),
					}

				case nodes.AT_AlterColumnType:
					d := cmd.Def.(nodes.ColumnDef)
					col := ""
					if cmd.Name != nil {
						col = *cmd.Name
					} else if d.Colname != nil {
						col = *d.Colname
					} else {
						return nil, fmt.Errorf("unknown name for alter column type")
					}
					tn, err := parseTypeName(d.TypeName)
					if err != nil {
						return nil, err
					}
					item.Subtype = ast.AT_AlterColumnType
					item.Def = &ast.ColumnDef{
						Colname:   col,
						TypeName:  tn,
						IsNotNull: isNotNull(d),
						IsArray:   isArray(d.TypeName),
					}

				case nodes.AT_DropColumn:
					item.Subtype = ast.AT_DropColumn

				case nodes.AT_DropNotNull:
					item.Subtype = ast.AT_DropNotNull

				case nodes.AT_SetNotNull:
					item.Subtype = ast.AT_SetNotNull

				default:
					continue
				}

				at.Cmds.Items = append(at.Cmds.Items, item)
			}
		}
		return at, nil

	case nodes.CommentStmt:
		switch n.Objtype {

		case nodes.OBJECT_COLUMN:
			col, tbl, err := parseColName(n.Object)
			if err != nil {
				return nil, fmt.Errorf("COMMENT ON COLUMN: %w", err)
			}
			return &ast.CommentOnColumnStmt{
				Col:     col,
				Table:   tbl,
				Comment: n.Comment,
			}, nil

		case nodes.OBJECT_SCHEMA:
			o, ok := n.Object.(nodes.String)
			if !ok {
				return nil, fmt.Errorf("COMMENT ON SCHEMA: unexpected node type: %T", n.Object)
			}
			return &ast.CommentOnSchemaStmt{
				Schema:  &ast.String{Str: o.Str},
				Comment: n.Comment,
			}, nil

		case nodes.OBJECT_TABLE:
			name, err := parseTableName(n.Object)
			if err != nil {
				return nil, fmt.Errorf("COMMENT ON TABLE: %w", err)
			}
			return &ast.CommentOnTableStmt{
				Table:   name,
				Comment: n.Comment,
			}, nil

		case nodes.OBJECT_TYPE:
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

	case nodes.CompositeTypeStmt:
		name, err := parseTypeName(n.Typevar)
		if err != nil {
			return nil, err
		}
		return &ast.CompositeTypeStmt{
			TypeName: name,
		}, nil

	case nodes.CreateStmt:
		name, err := parseTableName(*n.Relation)
		if err != nil {
			return nil, err
		}
		create := &ast.CreateTableStmt{
			Name:        name,
			IfNotExists: n.IfNotExists,
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
					IsNotNull: isNotNull(n),
					IsArray:   isArray(n.TypeName),
				})
			}
		}
		return create, nil

	case nodes.CreateEnumStmt:
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

	case nodes.CreateFunctionStmt:
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

	case nodes.CreateSchemaStmt:
		return &ast.CreateSchemaStmt{
			Name:        n.Schemaname,
			IfNotExists: n.IfNotExists,
		}, nil

	case nodes.DropStmt:
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

		case nodes.OBJECT_SCHEMA:
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

		case nodes.OBJECT_TABLE:
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

		case nodes.OBJECT_TYPE:
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

	case nodes.RenameStmt:
		switch n.RenameType {

		case nodes.OBJECT_COLUMN:
			tbl, err := parseTableName(*n.Relation)
			if err != nil {
				return nil, fmt.Errorf("nodes.RenameType: COLUMN: %w", err)
			}
			return &ast.RenameColumnStmt{
				Table:   tbl,
				Col:     &ast.ColumnRef{Name: *n.Subname},
				NewName: n.Newname,
			}, nil

		case nodes.OBJECT_TABLE:
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
