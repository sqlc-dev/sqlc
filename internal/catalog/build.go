package catalog

import (
	"fmt"
	"strings"

	"github.com/kyleconroy/sqlc/internal/pg"

	"github.com/davecgh/go-spew/spew"
	nodes "github.com/lfittl/pg_query_go/nodes"
)

func ParseRange(rv *nodes.RangeVar) (pg.FQN, error) {
	fqn := pg.FQN{
		Schema: "public",
	}
	if rv.Catalogname != nil {
		fqn.Catalog = *rv.Catalogname
	}
	if rv.Schemaname != nil {
		fqn.Schema = *rv.Schemaname
	}
	if rv.Relname != nil {
		fqn.Rel = *rv.Relname
	} else {
		return fqn, fmt.Errorf("range has empty relname")
	}
	return fqn, nil
}

func ParseList(list nodes.List) (pg.FQN, error) {
	parts := stringSlice(list)
	var fqn pg.FQN
	switch len(parts) {
	case 1:
		fqn = pg.FQN{
			Catalog: "",
			Schema:  "public",
			Rel:     parts[0],
		}
	case 2:
		fqn = pg.FQN{
			Catalog: "",
			Schema:  parts[0],
			Rel:     parts[1],
		}
	case 3:
		fqn = pg.FQN{
			Catalog: parts[0],
			Schema:  parts[1],
			Rel:     parts[2],
		}
	default:
		return fqn, fmt.Errorf("Invalid FQN: %s", join(list, "."))
	}
	return fqn, nil
}

func wrap(e pg.Error, loc int) pg.Error {
	return e
}

func Update(c *pg.Catalog, stmt nodes.Node) error {
	if false {
		spew.Dump(stmt)
	}
	raw, ok := stmt.(nodes.RawStmt)
	if !ok {
		return fmt.Errorf("expected RawStmt; got %T", stmt)
	}

	switch n := raw.Stmt.(type) {

	case nodes.AlterObjectSchemaStmt:
		switch n.ObjectType {

		case nodes.OBJECT_TABLE:
			fqn, err := ParseRange(n.Relation)
			if err != nil {
				return err
			}
			from, exists := c.Schemas[fqn.Schema]
			if !exists {
				return wrap(pg.ErrorSchemaDoesNotExist(fqn.Schema), raw.StmtLocation)
			}
			table, exists := from.Tables[fqn.Rel]
			if !exists {
				return wrap(pg.ErrorRelationDoesNotExist(fqn.Rel), raw.StmtLocation)
			}
			to, exists := c.Schemas[*n.Newschema]
			if !exists {
				return wrap(pg.ErrorSchemaDoesNotExist(*n.Newschema), raw.StmtLocation)
			}
			// Move the table
			delete(from.Tables, fqn.Rel)
			to.Tables[fqn.Rel] = table

		}

	case nodes.AlterTableStmt:
		var implemented bool
		for _, item := range n.Cmds.Items {
			switch cmd := item.(type) {
			case nodes.AlterTableCmd:
				switch cmd.Subtype {
				case nodes.AT_AddColumn:
					implemented = true
				case nodes.AT_AlterColumnType:
					implemented = true
				case nodes.AT_DropColumn:
					implemented = true
				case nodes.AT_DropNotNull:
					implemented = true
				case nodes.AT_SetNotNull:
					implemented = true
				}
			}
		}

		if !implemented {
			return nil
		}
		fqn, err := ParseRange(n.Relation)
		if err != nil {
			return err
		}
		schema, exists := c.Schemas[fqn.Schema]
		if !exists {
			return wrap(pg.ErrorSchemaDoesNotExist(fqn.Schema), raw.StmtLocation)
		}
		table, exists := schema.Tables[fqn.Rel]
		if !exists {
			return wrap(pg.ErrorRelationDoesNotExist(fqn.Rel), raw.StmtLocation)
		}

		for _, cmd := range n.Cmds.Items {
			switch cmd := cmd.(type) {
			case nodes.AlterTableCmd:
				idx := -1

				// Lookup column names for column-related commands
				switch cmd.Subtype {
				case nodes.AT_AlterColumnType,
					nodes.AT_DropColumn,
					nodes.AT_DropNotNull,
					nodes.AT_SetNotNull:

					for i, c := range table.Columns {
						if c.Name == *cmd.Name {
							idx = i
							break
						}
					}
					if idx < 0 && !cmd.MissingOk {
						return wrap(pg.ErrorColumnDoesNotExist(table.Name, *cmd.Name), raw.StmtLocation)
					}
					// If a missing column is allowed, skip this command
					if idx < 0 && cmd.MissingOk {
						continue
					}
				}

				switch cmd.Subtype {

				case nodes.AT_AddColumn:
					d := cmd.Def.(nodes.ColumnDef)
					for _, c := range table.Columns {
						if c.Name == *d.Colname {
							return wrap(pg.ErrorColumnAlreadyExists(table.Name, *d.Colname), d.Location)
						}
					}
					table.Columns = append(table.Columns, pg.Column{
						Name:     *d.Colname,
						DataType: join(d.TypeName.Names, "."),
						NotNull:  isNotNull(d),
						IsArray:  isArray(d.TypeName),
						Table:    fqn,
					})

				case nodes.AT_AlterColumnType:
					d := cmd.Def.(nodes.ColumnDef)
					table.Columns[idx].DataType = join(d.TypeName.Names, ".")
					table.Columns[idx].IsArray = isArray(d.TypeName)

				case nodes.AT_DropColumn:
					table.Columns = append(table.Columns[:idx], table.Columns[idx+1:]...)

				case nodes.AT_DropNotNull:
					table.Columns[idx].NotNull = false

				case nodes.AT_SetNotNull:
					table.Columns[idx].NotNull = true

				}

				schema.Tables[fqn.Rel] = table
			}
		}

	case nodes.CreateStmt:
		fqn, err := ParseRange(n.Relation)
		if err != nil {
			return err
		}
		schema, exists := c.Schemas[fqn.Schema]
		if !exists {
			return wrap(pg.ErrorSchemaDoesNotExist(fqn.Schema), raw.StmtLocation)
		}
		if _, exists := schema.Tables[fqn.Rel]; exists {
			return wrap(pg.ErrorRelationAlreadyExists(fqn.Rel), raw.StmtLocation)
		}
		table := pg.Table{
			Name: fqn.Rel,
		}
		for _, elt := range n.TableElts.Items {
			switch n := elt.(type) {
			case nodes.ColumnDef:
				colName := *n.Colname
				table.Columns = append(table.Columns, pg.Column{
					Name:     colName,
					DataType: join(n.TypeName.Names, "."),
					NotNull:  isNotNull(n),
					IsArray:  isArray(n.TypeName),
					Table:    fqn,
				})
			}
		}
		schema.Tables[fqn.Rel] = table

	case nodes.CreateEnumStmt:
		fqn, err := ParseList(n.TypeName)
		if err != nil {
			return err
		}
		schema, exists := c.Schemas[fqn.Schema]
		if !exists {
			return wrap(pg.ErrorSchemaDoesNotExist(fqn.Schema), raw.StmtLocation)
		}
		if _, exists := schema.Enums[fqn.Rel]; exists {
			return wrap(pg.ErrorTypeAlreadyExists(fqn.Rel), raw.StmtLocation)
		}
		schema.Enums[fqn.Rel] = pg.Enum{
			Name: fqn.Rel,
			Vals: stringSlice(n.Vals),
		}

	case nodes.CreateSchemaStmt:
		name := *n.Schemaname
		if _, exists := c.Schemas[name]; exists {
			return wrap(pg.ErrorSchemaAlreadyExists(name), raw.StmtLocation)
		}
		c.Schemas[name] = pg.NewSchema()

	case nodes.DropStmt:
		for _, obj := range n.Objects.Items {
			if n.RemoveType == nodes.OBJECT_TABLE || n.RemoveType == nodes.OBJECT_TYPE {
				var fqn pg.FQN
				var err error

				switch o := obj.(type) {
				case nodes.List:
					fqn, err = ParseList(o)
				case nodes.TypeName:
					fqn, err = ParseList(o.Names)
				default:
					return fmt.Errorf("nodes.DropStmt: unknown node in objects list: %T", o)
				}
				if err != nil {
					return err
				}

				schema, exists := c.Schemas[fqn.Schema]
				if !exists {
					return pg.ErrorSchemaDoesNotExist(fqn.Schema)
				}

				switch n.RemoveType {
				case nodes.OBJECT_TABLE:
					if _, exists := schema.Tables[fqn.Rel]; exists {
						delete(schema.Tables, fqn.Rel)
					} else if !n.MissingOk {
						return wrap(pg.ErrorRelationDoesNotExist(fqn.Rel), raw.StmtLocation)
					}

				case nodes.OBJECT_TYPE:
					if _, exists := schema.Enums[fqn.Rel]; exists {
						delete(schema.Enums, fqn.Rel)
					} else if !n.MissingOk {
						return wrap(pg.ErrorTypeDoesNotExist(fqn.Rel), raw.StmtLocation)
					}

				}

			}

			if n.RemoveType == nodes.OBJECT_SCHEMA {
				var name string
				switch o := obj.(type) {
				case nodes.String:
					name = o.Str
				default:
					return fmt.Errorf("nodes.DropStmt: unknown node in objects list: %T", o)
				}
				if _, exists := c.Schemas[name]; exists {
					delete(c.Schemas, name)
				} else if !n.MissingOk {
					return wrap(pg.ErrorSchemaDoesNotExist(name), raw.StmtLocation)
				}
			}
		}

	case nodes.RenameStmt:
		switch n.RenameType {
		case nodes.OBJECT_COLUMN:
			fqn, err := ParseRange(n.Relation)
			if err != nil {
				return err
			}
			schema, exists := c.Schemas[fqn.Schema]
			if !exists {
				return wrap(pg.ErrorSchemaDoesNotExist(fqn.Schema), raw.StmtLocation)
			}
			table, exists := schema.Tables[fqn.Rel]
			if !exists {
				return wrap(pg.ErrorRelationDoesNotExist(fqn.Rel), raw.StmtLocation)
			}
			idx := -1
			for i, c := range table.Columns {
				if c.Name == *n.Subname {
					idx = i
				}
				if c.Name == *n.Newname {
					return wrap(pg.ErrorColumnAlreadyExists(table.Name, c.Name), raw.StmtLocation)
				}
			}
			if idx < 0 {
				return wrap(pg.ErrorColumnDoesNotExist(table.Name, *n.Subname), raw.StmtLocation)
			}
			table.Columns[idx].Name = *n.Newname

		case nodes.OBJECT_TABLE:
			fqn, err := ParseRange(n.Relation)
			if err != nil {
				return err
			}
			schema, exists := c.Schemas[fqn.Schema]
			if !exists {
				return wrap(pg.ErrorSchemaDoesNotExist(fqn.Schema), raw.StmtLocation)
			}
			table, exists := schema.Tables[fqn.Rel]
			if !exists {
				return wrap(pg.ErrorRelationDoesNotExist(fqn.Rel), raw.StmtLocation)
			}
			if _, exists := schema.Tables[*n.Newname]; exists {
				return wrap(pg.ErrorRelationAlreadyExists(*n.Newname), raw.StmtLocation)
			}

			// Remove the table under the old name
			delete(schema.Tables, fqn.Rel)

			// Add the table under the new name
			table.Name = *n.Newname
			schema.Tables[*n.Newname] = table
		}

	case nodes.CreateFunctionStmt:
		fqn, err := ParseList(n.Funcname)
		if err != nil {
			return err
		}
		schema, exists := c.Schemas[fqn.Schema]
		if !exists {
			return wrap(pg.ErrorSchemaDoesNotExist(fqn.Schema), raw.StmtLocation)
		}
		arity := len(n.Parameters.Items)
		args := make([]pg.Argument, arity)
		for i, item := range n.Parameters.Items {
			arg := item.(nodes.FunctionParameter)
			var name string
			if arg.Name != nil {
				name = *arg.Name
			}
			args[i] = pg.Argument{
				Name:       name,
				DataType:   join(arg.ArgType.Names, "."),
				HasDefault: arg.Defexpr != nil,
			}
		}
		// TODO: support return parameter:
		// CREATE FUNCTION foo(bar TEXT, OUT quz bool) AS $$ SELECT true $$ LANGUAGE sql;
		schema.Funcs[fqn.Rel] = append(schema.Funcs[fqn.Rel], pg.Function{
			Name:       fqn.Rel,
			Arguments:  args,
			ReturnType: join(n.ReturnType.Names, "."),
		})

	case nodes.CommentStmt:
		switch n.Objtype {

		case nodes.OBJECT_SCHEMA:
			name := n.Object.(nodes.String).Str
			schema, exists := c.Schemas[name]
			if !exists {
				return wrap(pg.ErrorSchemaDoesNotExist(name), raw.StmtLocation)
			}
			if n.Comment != nil {
				schema.Comment = *n.Comment
			} else {
				schema.Comment = ""
			}
			c.Schemas[name] = schema

		case nodes.OBJECT_TABLE:
			fqn, err := ParseList(n.Object.(nodes.List))
			if err != nil {
				return err
			}
			schema, exists := c.Schemas[fqn.Schema]
			if !exists {
				return wrap(pg.ErrorSchemaDoesNotExist(fqn.Schema), raw.StmtLocation)
			}
			table, exists := schema.Tables[fqn.Rel]
			if !exists {
				return wrap(pg.ErrorRelationDoesNotExist(fqn.Rel), raw.StmtLocation)
			}
			if n.Comment != nil {
				table.Comment = *n.Comment
			} else {
				table.Comment = ""
			}
			schema.Tables[fqn.Rel] = table

		case nodes.OBJECT_COLUMN:
			colParts := stringSlice(n.Object.(nodes.List))
			var fqn pg.FQN
			var col string
			switch len(colParts) {
			case 2:
				col = colParts[1]
				fqn = pg.FQN{Schema: "public", Rel: colParts[0]}
			case 3:
				col = colParts[2]
				fqn = pg.FQN{Schema: colParts[0], Rel: colParts[1]}
			case 4:
				col = colParts[3]
				fqn = pg.FQN{Catalog: colParts[0], Schema: colParts[1], Rel: colParts[2]}
			default:
				return fmt.Errorf("column specifier %q is not the proper format, expected '[catalog.][schema.]colname.tablename'", strings.Join(colParts, "."))
			}
			schema, exists := c.Schemas[fqn.Schema]
			if !exists {
				return wrap(pg.ErrorSchemaDoesNotExist(fqn.Schema), raw.StmtLocation)
			}
			table, exists := schema.Tables[fqn.Rel]
			if !exists {
				return wrap(pg.ErrorRelationDoesNotExist(fqn.Rel), raw.StmtLocation)
			}
			idx := -1
			for i, c := range table.Columns {
				if c.Name == col {
					idx = i
				}
			}
			if idx < 0 {
				return wrap(pg.ErrorColumnDoesNotExist(table.Name, col), raw.StmtLocation)
			}
			if n.Comment != nil {
				table.Columns[idx].Comment = *n.Comment
			} else {
				table.Columns[idx].Comment = ""
			}

		case nodes.OBJECT_TYPE:
			fqn, err := ParseList(n.Object.(nodes.TypeName).Names)
			if err != nil {
				return err
			}
			schema, exists := c.Schemas[fqn.Schema]
			if !exists {
				return wrap(pg.ErrorSchemaDoesNotExist(fqn.Schema), raw.StmtLocation)
			}
			enum, exists := schema.Enums[fqn.Rel]
			if !exists {
				return wrap(pg.ErrorRelationDoesNotExist(fqn.Rel), raw.StmtLocation)
			}
			if n.Comment != nil {
				enum.Comment = *n.Comment
			} else {
				enum.Comment = ""
			}
			schema.Enums[fqn.Rel] = enum

		}

	}
	return nil
}

func stringSlice(list nodes.List) []string {
	items := []string{}
	for _, item := range list.Items {
		if n, ok := item.(nodes.String); ok {
			items = append(items, n.Str)
		}
	}
	return items
}

func join(list nodes.List, sep string) string {
	return strings.Join(stringSlice(list), sep)
}

func isArray(n *nodes.TypeName) bool {
	if n == nil {
		return false
	}
	return len(n.ArrayBounds.Items) > 0
}

func isNotNull(n nodes.ColumnDef) bool {
	if n.IsNotNull {
		return true
	}
	for _, c := range n.Constraints.Items {
		switch n := c.(type) {
		case nodes.Constraint:
			if n.Contype == nodes.CONSTR_NOTNULL {
				return true
			}
			if n.Contype == nodes.CONSTR_PRIMARY {
				return true
			}
		}
	}
	return false
}

func ToColumn(n *nodes.TypeName) pg.Column {
	if n == nil {
		panic("can't build column for nil type name")
	}
	return pg.Column{
		DataType: join(n.Names, "."),
		NotNull:  true, // XXX: How do we know if this should be null?
		IsArray:  isArray(n),
	}
}
