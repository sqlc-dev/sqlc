package postgresql

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/kyleconroy/sqlc/internal/sql/ast"

	pg "github.com/lfittl/pg_query_go"
	nodes "github.com/lfittl/pg_query_go/nodes"
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

func parseTypeName(node nodes.Node) (*ast.TypeName, error) {
	switch n := node.(type) {

	case nodes.TypeName:
		return parseTypeName(n.Names)

	case nodes.List:
		parts := stringSlice(n)
		switch len(parts) {
		case 1:
			return &ast.TypeName{
				Name: parts[0],
			}, nil
		case 2:
			return &ast.TypeName{
				Schema: parts[0],
				Name:   parts[1],
			}, nil
		default:
			return nil, fmt.Errorf("invalid type name: %s", join(n, "."))
		}

	default:
		return nil, fmt.Errorf("parseTypeName: unexpected node type: %T", n)
	}
}

func parseTableName(node nodes.Node) (*ast.TableName, error) {
	switch n := node.(type) {

	case nodes.List:
		parts := stringSlice(n)
		switch len(parts) {
		case 1:
			return &ast.TableName{
				Name: parts[0],
			}, nil
		case 2:
			return &ast.TableName{
				Schema: parts[0],
				Name:   parts[1],
			}, nil
		case 3:
			return &ast.TableName{
				Catalog: parts[0],
				Schema:  parts[1],
				Name:    parts[2],
			}, nil
		default:
			return nil, fmt.Errorf("invalid table name: %s", join(n, "."))
		}

	case nodes.RangeVar:
		name := ast.TableName{}
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

	default:
		return nil, fmt.Errorf("parseTableName: unexpected node type: %T", n)
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
			Raw: &ast.RawStmt{Stmt: n},
		})
	}
	return stmts, nil
}

func translate(node nodes.Node) (ast.Node, error) {
	switch n := node.(type) {

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
					item.Subtype = ast.AT_AddColumn
					item.Def = &ast.ColumnDef{
						Colname:   *d.Colname,
						TypeName:  &ast.TypeName{Name: join(d.TypeName.Names, ".")},
						IsNotNull: isNotNull(d),
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
					item.Subtype = ast.AT_AlterColumnType
					item.Def = &ast.ColumnDef{
						Colname:   col,
						TypeName:  &ast.TypeName{Name: join(d.TypeName.Names, ".")},
						IsNotNull: isNotNull(d),
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
				create.Cols = append(create.Cols, &ast.ColumnDef{
					Colname:   *n.Colname,
					TypeName:  &ast.TypeName{Name: join(n.TypeName.Names, ".")},
					IsNotNull: isNotNull(n),
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

	case nodes.CreateSchemaStmt:
		return &ast.CreateSchemaStmt{
			Name:        n.Schemaname,
			IfNotExists: n.IfNotExists,
		}, nil

	case nodes.DropStmt:
		switch n.RemoveType {

		case nodes.OBJECT_SCHEMA:
			drop := &ast.DropSchemaStmt{
				MissingOk: n.MissingOk,
			}
			for _, obj := range n.Objects.Items {
				val, ok := obj.(nodes.String)
				if !ok {
					return nil, fmt.Errorf("nodes.DropStmt: unknown type in objects list: %T", obj)
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
					return nil, err
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
					return nil, err
				}
				drop.Types = append(drop.Types, name)
			}
			return drop, nil

		}
		return nil, errSkip

	case nodes.RenameStmt:
		switch n.RenameType {

		case nodes.OBJECT_TABLE:
			tbl, err := parseTableName(*n.Relation)
			if err != nil {
				return nil, err
			}
			return &ast.RenameTableStmt{
				Table:   tbl,
				NewName: n.Newname,
			}, nil

		}
		return nil, errSkip

	default:
		return nil, errSkip
	}
}
