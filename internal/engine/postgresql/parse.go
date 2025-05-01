package postgresql

import (
	"errors"
	"fmt"
	"io"
	"strings"

	nodes "github.com/pganalyze/pg_query_go/v6"

	"github.com/sqlc-dev/sqlc/internal/engine/postgresql/parser"
	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

func stringSlice(list *nodes.List) []string {
	items := []string{}
	for _, item := range list.Items {
		if n, ok := item.Node.(*nodes.Node_String_); ok {
			items = append(items, n.String_.Sval)
		}
	}
	return items
}

func stringSliceFromNodes(s []*nodes.Node) []string {
	var items []string
	for _, item := range s {
		if n, ok := item.Node.(*nodes.Node_String_); ok {
			items = append(items, n.String_.Sval)
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

func parseRelation(in *nodes.Node) (*relation, error) {
	switch n := in.Node.(type) {
	case *nodes.Node_List:
		return parseRelationFromNodes(n.List.Items)
	case *nodes.Node_RangeVar:
		return parseRelationFromRangeVar(n.RangeVar), nil
	case *nodes.Node_TypeName:
		return parseRelationFromNodes(n.TypeName.Names)
	default:
		return nil, fmt.Errorf("unexpected node type: %T", n)
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
	contents, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	tree, err := Parse(string(contents))
	if err != nil {
		pErr := normalizeErr(err)
		return nil, pErr
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

func normalizeErr(err error) error {
	//TODO: errors.As complains that *parser.Error does not implement error
	if pErr, ok := err.(*parser.Error); ok {
		sErr := &sqlerr.Error{
			Message: pErr.Message,
			//Err:      pErr,
			Line:     pErr.Lineno,
			Location: pErr.Cursorpos,
		}
		return sErr
	}

	return err
}

// https://www.postgresql.org/docs/current/sql-syntax-lexical.html#SQL-SYNTAX-COMMENTS
func (p *Parser) CommentSyntax() source.CommentSyntax {
	return source.CommentSyntax{
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
				OldValue: makeString(n.OldVal),
				NewValue: makeString(n.NewVal),
			}, nil
		} else {
			return &ast.AlterTypeAddValueStmt{
				Type:               rel.TypeName(),
				NewValue:           makeString(n.NewVal),
				NewValHasNeighbor:  len(n.NewValNeighbor) > 0,
				NewValNeighbor:     makeString(n.NewValNeighbor),
				NewValIsAfter:      n.NewValIsAfter,
				SkipIfNewValExists: n.SkipIfNewValExists,
			}, nil
		}

	case *nodes.Node_AlterObjectSchemaStmt:
		n := inner.AlterObjectSchemaStmt
		switch n.ObjectType {

		case nodes.ObjectType_OBJECT_TABLE, nodes.ObjectType_OBJECT_VIEW, nodes.ObjectType_OBJECT_MATVIEW:
			rel := parseRelationFromRangeVar(n.Relation)
			return &ast.AlterTableSetSchemaStmt{
				Table:     rel.TableName(),
				NewSchema: makeString(n.Newschema),
				MissingOk: n.MissingOk,
			}, nil

		case nodes.ObjectType_OBJECT_TYPE:
			rel, err := parseRelation(n.Object)
			if err != nil {
				return nil, err
			}
			return &ast.AlterTypeSetSchemaStmt{
				Type:      rel.TypeName(),
				NewSchema: makeString(n.Newschema),
			}, nil
		}
		return nil, errSkip

	case *nodes.Node_AlterTableStmt:
		n := inner.AlterTableStmt
		rel := parseRelationFromRangeVar(n.Relation)
		at := &ast.AlterTableStmt{
			Table:     rel.TableName(),
			Cmds:      &ast.List{},
			MissingOk: n.MissingOk,
		}
		for _, cmd := range n.Cmds {
			switch cmdOneOf := cmd.Node.(type) {
			case *nodes.Node_AlterTableCmd:
				altercmd := cmdOneOf.AlterTableCmd
				item := &ast.AlterTableCmd{Name: &altercmd.Name, MissingOk: altercmd.MissingOk}

				switch altercmd.Subtype {
				case nodes.AlterTableType_AT_AddColumn:
					d, ok := altercmd.Def.Node.(*nodes.Node_ColumnDef)
					if !ok {
						return nil, fmt.Errorf("expected alter table definition to be a ColumnDef")
					}

					rel, err := parseRelationFromNodes(d.ColumnDef.TypeName.Names)
					if err != nil {
						return nil, err
					}
					item.Subtype = ast.AT_AddColumn
					item.Def = &ast.ColumnDef{
						Colname:   d.ColumnDef.Colname,
						TypeName:  rel.TypeName(),
						IsNotNull: isNotNull(d.ColumnDef),
						IsArray:   isArray(d.ColumnDef.TypeName),
						ArrayDims: len(d.ColumnDef.TypeName.ArrayBounds),
					}

				case nodes.AlterTableType_AT_AlterColumnType:
					d, ok := altercmd.Def.Node.(*nodes.Node_ColumnDef)
					if !ok {
						return nil, fmt.Errorf("expected alter table definition to be a ColumnDef")
					}
					col := ""
					if altercmd.Name != "" {
						col = altercmd.Name
					} else if d.ColumnDef.Colname != "" {
						col = d.ColumnDef.Colname
					} else {
						return nil, fmt.Errorf("unknown name for alter column type")
					}
					rel, err := parseRelationFromNodes(d.ColumnDef.TypeName.Names)
					if err != nil {
						return nil, err
					}
					item.Subtype = ast.AT_AlterColumnType
					item.Def = &ast.ColumnDef{
						Colname:   col,
						TypeName:  rel.TypeName(),
						IsNotNull: isNotNull(d.ColumnDef),
						IsArray:   isArray(d.ColumnDef.TypeName),
						ArrayDims: len(d.ColumnDef.TypeName.ArrayBounds),
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
				Comment: makeString(n.Comment),
			}, nil

		case nodes.ObjectType_OBJECT_SCHEMA:
			o, ok := n.Object.Node.(*nodes.Node_String_)
			if !ok {
				return nil, fmt.Errorf("COMMENT ON SCHEMA: unexpected node type: %T", n.Object)
			}
			return &ast.CommentOnSchemaStmt{
				Schema:  &ast.String{Str: o.String_.Sval},
				Comment: makeString(n.Comment),
			}, nil

		case nodes.ObjectType_OBJECT_TABLE:
			rel, err := parseRelation(n.Object)
			if err != nil {
				return nil, fmt.Errorf("COMMENT ON TABLE: %w", err)
			}
			return &ast.CommentOnTableStmt{
				Table:   rel.TableName(),
				Comment: makeString(n.Comment),
			}, nil

		case nodes.ObjectType_OBJECT_TYPE:
			rel, err := parseRelation(n.Object)
			if err != nil {
				return nil, err
			}
			return &ast.CommentOnTypeStmt{
				Type:    rel.TypeName(),
				Comment: makeString(n.Comment),
			}, nil

		case nodes.ObjectType_OBJECT_VIEW:
			rel, err := parseRelation(n.Object)
			if err != nil {
				return nil, fmt.Errorf("COMMENT ON VIEW: %w", err)
			}
			return &ast.CommentOnViewStmt{
				View:    rel.TableName(),
				Comment: makeString(n.Comment),
			}, nil

		}
		return nil, errSkip

	case *nodes.Node_CompositeTypeStmt:
		n := inner.CompositeTypeStmt
		rel := parseRelationFromRangeVar(n.Typevar)
		return &ast.CompositeTypeStmt{
			TypeName: rel.TypeName(),
		}, nil

	case *nodes.Node_CreateStmt:
		n := inner.CreateStmt
		rel := parseRelationFromRangeVar(n.Relation)
		create := &ast.CreateTableStmt{
			Name:        rel.TableName(),
			IfNotExists: n.IfNotExists,
		}
		for _, node := range n.InhRelations {
			switch item := node.Node.(type) {
			case *nodes.Node_RangeVar:
				if item.RangeVar.Inh {
					rel := parseRelationFromRangeVar(item.RangeVar)
					create.Inherits = append(create.Inherits, rel.TableName())
				}
			}
		}
		primaryKey := make(map[string]bool)
		for _, elt := range n.TableElts {
			switch item := elt.Node.(type) {
			case *nodes.Node_Constraint:
				if item.Constraint.Contype == nodes.ConstrType_CONSTR_PRIMARY {
					for _, key := range item.Constraint.Keys {
						// FIXME: Possible nil pointer dereference
						primaryKey[key.Node.(*nodes.Node_String_).String_.Sval] = true
					}
				}

			case *nodes.Node_TableLikeClause:
				rel := parseRelationFromRangeVar(item.TableLikeClause.Relation)
				create.ReferTable = rel.TableName()
			}
		}
		for _, elt := range n.TableElts {
			switch item := elt.Node.(type) {
			case *nodes.Node_ColumnDef:
				rel, err := parseRelationFromNodes(item.ColumnDef.TypeName.Names)
				if err != nil {
					return nil, err
				}

				primary := false
				for _, con := range item.ColumnDef.Constraints {
					if constraint, ok := con.Node.(*nodes.Node_Constraint); ok {
						primary = constraint.Constraint.Contype == nodes.ConstrType_CONSTR_PRIMARY
					}
				}

				create.Cols = append(create.Cols, &ast.ColumnDef{
					Colname:    item.ColumnDef.Colname,
					TypeName:   rel.TypeName(),
					IsNotNull:  isNotNull(item.ColumnDef) || primaryKey[item.ColumnDef.Colname],
					IsArray:    isArray(item.ColumnDef.TypeName),
					ArrayDims:  len(item.ColumnDef.TypeName.ArrayBounds),
					PrimaryKey: primary,
				})
			}
		}
		return create, nil

	case *nodes.Node_CreateEnumStmt:
		n := inner.CreateEnumStmt
		rel, err := parseRelationFromNodes(n.TypeName)
		if err != nil {
			return nil, err
		}
		stmt := &ast.CreateEnumStmt{
			TypeName: rel.TypeName(),
			Vals:     &ast.List{},
		}
		for _, val := range n.Vals {
			switch v := val.Node.(type) {
			case *nodes.Node_String_:
				stmt.Vals.Items = append(stmt.Vals.Items, &ast.String{
					Str: v.String_.Sval,
				})
			}
		}
		return stmt, nil

	case *nodes.Node_CreateFunctionStmt:
		n := inner.CreateFunctionStmt
		fn, err := parseRelationFromNodes(n.Funcname)
		if err != nil {
			return nil, err
		}
		var rt *ast.TypeName
		if n.ReturnType != nil {
			rel, err := parseRelationFromNodes(n.ReturnType.Names)
			if err != nil {
				return nil, err
			}
			rt = rel.TypeName()
		}
		stmt := &ast.CreateFunctionStmt{
			Func:       fn.FuncName(),
			ReturnType: rt,
			Replace:    n.Replace,
			Params:     &ast.List{},
		}
		for _, item := range n.Parameters {
			arg := item.Node.(*nodes.Node_FunctionParameter).FunctionParameter
			rel, err := parseRelationFromNodes(arg.ArgType.Names)
			if err != nil {
				return nil, err
			}
			mode, err := convertFuncParamMode(arg.Mode)
			if err != nil {
				return nil, err
			}
			fp := &ast.FuncParam{
				Name: &arg.Name,
				Type: rel.TypeName(),
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
			Name:        makeString(n.Schemaname),
			IfNotExists: n.IfNotExists,
		}, nil

	case *nodes.Node_DropStmt:
		n := inner.DropStmt
		switch n.RemoveType {

		case nodes.ObjectType_OBJECT_FUNCTION:
			drop := &ast.DropFunctionStmt{
				MissingOk: n.MissingOk,
			}
			for _, obj := range n.Objects {
				nowa, ok := obj.Node.(*nodes.Node_ObjectWithArgs)
				if !ok {
					return nil, fmt.Errorf("nodes.DropStmt: FUNCTION: unknown type in objects list: %T", obj)
				}
				owa := nowa.ObjectWithArgs
				fn, err := parseRelationFromNodes(owa.Objname)
				if err != nil {
					return nil, fmt.Errorf("nodes.DropStmt: FUNCTION: %w", err)
				}
				args := make([]*ast.TypeName, len(owa.Objargs))
				for i, objarg := range owa.Objargs {
					tn, ok := objarg.Node.(*nodes.Node_TypeName)
					if !ok {
						return nil, fmt.Errorf("nodes.DropStmt: FUNCTION: unknown type in objargs list: %T", objarg)
					}
					at, err := parseRelationFromNodes(tn.TypeName.Names)
					if err != nil {
						return nil, fmt.Errorf("nodes.DropStmt: FUNCTION: %w", err)
					}
					args[i] = at.TypeName()
				}
				drop.Funcs = append(drop.Funcs, &ast.FuncSpec{
					Name:    fn.FuncName(),
					Args:    args,
					HasArgs: !owa.ArgsUnspecified,
				})
			}
			return drop, nil

		case nodes.ObjectType_OBJECT_SCHEMA:
			drop := &ast.DropSchemaStmt{
				MissingOk: n.MissingOk,
			}
			for _, obj := range n.Objects {
				val, ok := obj.Node.(*nodes.Node_String_)
				if !ok {
					return nil, fmt.Errorf("nodes.DropStmt: SCHEMA: unknown type in objects list: %T", obj)
				}
				drop.Schemas = append(drop.Schemas, &ast.String{Str: val.String_.Sval})
			}
			return drop, nil

		case nodes.ObjectType_OBJECT_TABLE, nodes.ObjectType_OBJECT_VIEW, nodes.ObjectType_OBJECT_MATVIEW:
			drop := &ast.DropTableStmt{
				IfExists: n.MissingOk,
			}
			for _, obj := range n.Objects {
				name, err := parseRelation(obj)
				if err != nil {
					return nil, fmt.Errorf("nodes.DropStmt: TABLE: %w", err)
				}
				drop.Tables = append(drop.Tables, name.TableName())
			}
			return drop, nil

		case nodes.ObjectType_OBJECT_TYPE:
			drop := &ast.DropTypeStmt{
				IfExists: n.MissingOk,
			}
			for _, obj := range n.Objects {
				name, err := parseRelation(obj)
				if err != nil {
					return nil, fmt.Errorf("nodes.DropStmt: TYPE: %w", err)
				}
				drop.Types = append(drop.Types, name.TypeName())
			}
			return drop, nil

		}
		return nil, errSkip

	case *nodes.Node_RenameStmt:
		n := inner.RenameStmt
		switch n.RenameType {

		case nodes.ObjectType_OBJECT_COLUMN:
			rel := parseRelationFromRangeVar(n.Relation)
			return &ast.RenameColumnStmt{
				Table:     rel.TableName(),
				Col:       &ast.ColumnRef{Name: n.Subname},
				NewName:   makeString(n.Newname),
				MissingOk: n.MissingOk,
			}, nil

		case nodes.ObjectType_OBJECT_TABLE, nodes.ObjectType_OBJECT_MATVIEW, nodes.ObjectType_OBJECT_VIEW:
			rel := parseRelationFromRangeVar(n.Relation)
			return &ast.RenameTableStmt{
				Table:     rel.TableName(),
				NewName:   makeString(n.Newname),
				MissingOk: n.MissingOk,
			}, nil

		case nodes.ObjectType_OBJECT_TYPE:
			rel, err := parseRelation(n.Object)
			if err != nil {
				return nil, fmt.Errorf("nodes.RenameStmt: TYPE: %w", err)
			}
			return &ast.RenameTypeStmt{
				Type:    rel.TypeName(),
				NewName: makeString(n.Newname),
			}, nil

		}
		return nil, errSkip

	default:
		return convert(node)
	}
}
