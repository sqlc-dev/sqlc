package schema

import (
	"fmt"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/core/analyzer"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func Apply(cat *core.Catalog, n ast.Node) error {
	switch v := n.(type) {
	case nil:
		return nil
	case *ast.RawStmt:
		return Apply(cat, v.Stmt)
	case *ast.List:
		for _, it := range v.Items {
			if err := Apply(cat, it); err != nil {
				return err
			}
		}
		return nil
	case *ast.CreateTableStmt:
		return applyCreateTable(cat, v)
	case *ast.DropTableStmt:
		return applyDropTable(cat, v)
	case *ast.CreateEnumStmt:
		return applyCreateEnum(cat, v)
	case *ast.CreateExtensionStmt:
		if v.Extname == nil {
			return nil
		}
		return cat.LoadExtension(*v.Extname)
	case *ast.CreateFunctionStmt:
		return applyCreateFunction(cat, v)
	case *ast.AlterTableStmt:
		return applyAlterTable(cat, v)
	case *ast.RenameColumnStmt:
		return applyRenameColumn(cat, v)
	case *ast.RenameTableStmt:
		return applyRenameTable(cat, v)
	case *ast.ViewStmt:
		return applyView(cat, v.View, v.Aliases, v.Query, v.Replace)
	case *ast.CreateTableAsStmt:
		if v.Into == nil {
			return nil
		}
		return applyView(cat, v.Into.Rel, v.Into.ColNames, v.Query, false)
	}
	return nil
}

// applyView records a view, or a table created from a query, as a relation
// whose columns are the ones its query selects.
func applyView(cat *core.Catalog, rel *ast.RangeVar, aliases *ast.List, query ast.Node, replace bool) error {
	if rel == nil || rel.Relname == nil {
		return fmt.Errorf("create view with nil name")
	}
	name := *rel.Relname
	schema := ""
	if rel.Schemaname != nil {
		schema = *rel.Schemaname
	}
	nsOID, err := resolveOrCreateNamespace(cat, schema)
	if err != nil {
		return err
	}
	if existing, err := cat.ClassOID(nsOID, name); err == nil {
		if !replace {
			return fmt.Errorf("relation %q already exists", name)
		}
		if err := cat.DropClass(existing); err != nil {
			return err
		}
	}

	sel, ok := query.(*ast.SelectStmt)
	if !ok {
		return fmt.Errorf("view %q: unsupported query %T", name, query)
	}
	res, err := analyzer.Prepare(cat, sel)
	if err != nil {
		return fmt.Errorf("view %q: %w", name, err)
	}

	classOID, err := cat.CreateClass(nsOID, name, "v")
	if err != nil {
		return err
	}
	names := listStrings(aliases)
	for i, col := range res.Columns {
		colName := col.Name
		if i < len(names) {
			colName = names[i]
		}
		typeOID := col.TypeOID
		if typeOID == 0 {
			typeOID, err = cat.ResolveTypeName("any")
			if err != nil {
				return err
			}
		}
		if err := cat.CreateAttributeSpec(core.AttributeSpec{
			ClassOID: classOID,
			Name:     colName,
			TypeOID:  typeOID,
			Num:      i + 1,
			NotNull:  col.NotNull,
			DeclType: col.DataType,
		}); err != nil {
			return fmt.Errorf("view %s.%s: %w", name, colName, err)
		}
	}
	return nil
}

func listStrings(l *ast.List) []string {
	if l == nil {
		return nil
	}
	out := make([]string, 0, len(l.Items))
	for _, item := range l.Items {
		if s, ok := item.(*ast.String); ok {
			out = append(out, s.Str)
		}
	}
	return out
}

func applyCreateTable(cat *core.Catalog, stmt *ast.CreateTableStmt) error {
	if stmt.Name == nil {
		return fmt.Errorf("create table with nil name")
	}
	nsOID, err := resolveOrCreateNamespace(cat, stmt.Name.Schema)
	if err != nil {
		return err
	}
	if _, err := cat.ClassOID(nsOID, stmt.Name.Name); err == nil {
		if stmt.IfNotExists {
			return nil
		}
		return fmt.Errorf("relation %q already exists", stmt.Name.Name)
	}
	classOID, err := cat.CreateClass(nsOID, stmt.Name.Name, "r")
	if err != nil {
		return err
	}
	for i, col := range stmt.Cols {
		if col == nil || col.TypeName == nil {
			return fmt.Errorf("column %d on %q: missing type", i+1, stmt.Name.Name)
		}
		typeOID, err := columnTypeOID(cat, col)
		if err != nil {
			return fmt.Errorf("column %s.%s: %w", stmt.Name.Name, col.Colname, err)
		}
		if err := cat.CreateAttributeSpec(core.AttributeSpec{
			ClassOID:     classOID,
			Name:         col.Colname,
			TypeOID:      typeOID,
			Num:          i + 1,
			NotNull:      col.IsNotNull || col.PrimaryKey,
			IsPrimaryKey: col.PrimaryKey,
			DeclType:     col.TypeName.Name,
		}); err != nil {
			return fmt.Errorf("attr %s.%s: %w", stmt.Name.Name, col.Colname, err)
		}
	}
	return nil
}

func applyDropTable(cat *core.Catalog, stmt *ast.DropTableStmt) error {
	for _, tn := range stmt.Tables {
		if tn == nil {
			continue
		}
		nsOID, err := cat.NamespaceOID(nsName(tn.Schema))
		if err != nil {
			if stmt.IfExists {
				continue
			}
			return err
		}
		classOID, err := cat.ClassOID(nsOID, tn.Name)
		if err != nil {
			if stmt.IfExists {
				continue
			}
			return fmt.Errorf("drop table %q: %w", tn.Name, err)
		}
		if err := cat.DropClass(classOID); err != nil {
			return fmt.Errorf("drop table %q: %w", tn.Name, err)
		}
	}
	return nil
}

func applyAlterTable(cat *core.Catalog, stmt *ast.AlterTableStmt) error {
	table := stmt.Table
	if table == nil && stmt.Relation != nil {
		table = rangeVarTableName(stmt.Relation)
	}
	classOID, err := lookupClass(cat, table)
	if err != nil {
		if stmt.MissingOk {
			return nil
		}
		return err
	}
	for _, item := range listItems(stmt.Cmds) {
		cmd, ok := item.(*ast.AlterTableCmd)
		if !ok {
			continue
		}
		switch cmd.Subtype {
		case ast.AT_AddColumn:
			if cmd.Def == nil {
				continue
			}
			typeOID, err := columnTypeOID(cat, cmd.Def)
			if err != nil {
				return err
			}
			num, err := cat.NextAttributeNum(classOID)
			if err != nil {
				return err
			}
			if err := cat.CreateAttributeSpec(core.AttributeSpec{
				ClassOID:     classOID,
				Name:         cmd.Def.Colname,
				TypeOID:      typeOID,
				Num:          num,
				NotNull:      cmd.Def.IsNotNull || cmd.Def.PrimaryKey,
				IsPrimaryKey: cmd.Def.PrimaryKey,
				DeclType:     cmd.Def.TypeName.Name,
			}); err != nil {
				return err
			}
		case ast.AT_DropColumn:
			if cmd.Name == nil {
				continue
			}
			if err := cat.DropAttribute(classOID, *cmd.Name); err != nil {
				return err
			}
		case ast.AT_AlterColumnType:
			if cmd.Def == nil {
				continue
			}
			name := cmd.Def.Colname
			if cmd.Name != nil && *cmd.Name != "" {
				name = *cmd.Name
			}
			typeOID, err := columnTypeOID(cat, cmd.Def)
			if err != nil {
				return err
			}
			if err := cat.SetAttributeType(classOID, name, typeOID, cmd.Def.TypeName.Name); err != nil {
				return err
			}
			// An engine that reports a column's whole new definition also
			// reports whether it still accepts NULL.
			if cmd.Def.IsNotNull {
				if err := cat.SetAttributeNotNull(classOID, name, true); err != nil {
					return err
				}
			}
		case ast.AT_SetNotNull, ast.AT_DropNotNull:
			if cmd.Name == nil {
				continue
			}
			if err := cat.SetAttributeNotNull(classOID, *cmd.Name, cmd.Subtype == ast.AT_SetNotNull); err != nil {
				return err
			}
		}
	}
	return nil
}

func applyRenameColumn(cat *core.Catalog, stmt *ast.RenameColumnStmt) error {
	if stmt.Col == nil || stmt.NewName == nil {
		return nil
	}
	classOID, err := lookupClass(cat, stmt.Table)
	if err != nil {
		if stmt.MissingOk {
			return nil
		}
		return err
	}
	// Engines report the column either as a field list or as a bare name.
	name := stmt.Col.Name
	if names := listStrings(stmt.Col.Fields); len(names) > 0 {
		name = names[len(names)-1]
	}
	if name == "" {
		return nil
	}
	return cat.RenameAttribute(classOID, name, *stmt.NewName)
}

func applyRenameTable(cat *core.Catalog, stmt *ast.RenameTableStmt) error {
	if stmt.NewName == nil {
		return nil
	}
	classOID, err := lookupClass(cat, stmt.Table)
	if err != nil {
		if stmt.MissingOk {
			return nil
		}
		return err
	}
	return cat.RenameClass(classOID, *stmt.NewName)
}

func lookupClass(cat *core.Catalog, table *ast.TableName) (int64, error) {
	if table == nil {
		return 0, fmt.Errorf("missing table name")
	}
	nsOID, err := cat.NamespaceOID(nsName(table.Schema))
	if err != nil {
		return 0, err
	}
	return cat.ClassOID(nsOID, table.Name)
}

func rangeVarTableName(rv *ast.RangeVar) *ast.TableName {
	tn := &ast.TableName{}
	if rv.Schemaname != nil {
		tn.Schema = *rv.Schemaname
	}
	if rv.Relname != nil {
		tn.Name = *rv.Relname
	}
	return tn
}

func listItems(l *ast.List) []ast.Node {
	if l == nil {
		return nil
	}
	return l.Items
}

func applyCreateEnum(cat *core.Catalog, stmt *ast.CreateEnumStmt) error {
	if stmt.TypeName == nil {
		return fmt.Errorf("create type with nil name")
	}
	name := core.TypeNameString(stmt.TypeName)
	if name == "" {
		return fmt.Errorf("create type with empty name")
	}
	if _, err := cat.TypeOID(name); err == nil {
		return nil
	}
	_, err := cat.CreateUserType(name, "E")
	return err
}

func applyCreateFunction(cat *core.Catalog, stmt *ast.CreateFunctionStmt) error {
	// A procedure returns nothing, so there is no result for a query to
	// select and nothing worth recording.
	if stmt.Func == nil || stmt.Func.Name == "" || stmt.ReturnType == nil {
		return nil
	}
	returnOID, err := cat.ResolveType(stmt.ReturnType)
	if err != nil {
		return fmt.Errorf("function %q: %w", stmt.Func.Name, err)
	}
	var args []core.ProcArg
	if stmt.Params != nil {
		for _, item := range stmt.Params.Items {
			p, ok := item.(*ast.FuncParam)
			if !ok {
				continue
			}
			switch p.Mode {
			case ast.FuncParamOut, ast.FuncParamTable:
				continue
			}
			argOID, err := cat.ResolveType(p.Type)
			if err != nil {
				return fmt.Errorf("function %q: %w", stmt.Func.Name, err)
			}
			arg := core.ProcArg{TypeOID: argOID, HasDefault: p.DefExpr != nil}
			if p.Name != nil {
				arg.Name = *p.Name
			}
			args = append(args, arg)
		}
	}
	_, err = cat.CreateProc(core.ProcSpec{
		Name:          stmt.Func.Name,
		ReturnTypeOID: returnOID,
		ReturnSet:     stmt.ReturnType != nil && stmt.ReturnType.Setof,
		Args:          args,
	})
	return err
}

func nsName(schema string) string {
	if schema == "" {
		return "public"
	}
	return schema
}

func resolveOrCreateNamespace(cat *core.Catalog, schema string) (int64, error) {
	name := nsName(schema)
	if oid, err := cat.NamespaceOID(name); err == nil {
		return oid, nil
	}
	return cat.CreateNamespace(name)
}

// columnTypeOID resolves a column's type. Engines report an array column
// either on the type name or on the column itself.
func columnTypeOID(cat *core.Catalog, col *ast.ColumnDef) (int64, error) {
	name := core.TypeNameString(col.TypeName)
	if name == "" {
		return 0, fmt.Errorf("missing type name")
	}
	if (col.IsArray || col.ArrayDims > 0) && !strings.HasSuffix(name, core.ArraySuffix) {
		name += core.ArraySuffix
	}
	return cat.ResolveTypeName(name)
}
