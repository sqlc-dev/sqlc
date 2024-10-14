package compiler

import (
	"fmt"
	"math/rand"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
	"github.com/sqlc-dev/sqlc/internal/sql/rewrite"
)

type QueryCatalog struct {
	catalog     *catalog.Catalog
	ctes        map[string]*Table
	fromClauses map[string]*Table
	embeds      rewrite.EmbedSet
}

func (comp *Compiler) buildQueryCatalog(c *catalog.Catalog, node ast.Node, embeds rewrite.EmbedSet) (*QueryCatalog, error) {
	var with *ast.WithClause
	var from *ast.List
	switch n := node.(type) {
	case *ast.DeleteStmt:
		with = n.WithClause
	case *ast.InsertStmt:
		with = n.WithClause
	case *ast.UpdateStmt:
		with = n.WithClause
		from = n.FromClause
	case *ast.SelectStmt:
		with = n.WithClause
		from = n.FromClause
	default:
		with = nil
		from = nil
	}
	qc := &QueryCatalog{
		catalog:     c,
		ctes:        map[string]*Table{},
		fromClauses: map[string]*Table{},
		embeds:      embeds,
	}
	if with != nil {
		for _, item := range with.Ctes.Items {
			if cte, ok := item.(*ast.CommonTableExpr); ok {
				cols, err := comp.outputColumns(qc, cte.Ctequery)
				if err != nil {
					return nil, err
				}
				var names []string
				if cte.Aliascolnames != nil {
					for _, item := range cte.Aliascolnames.Items {
						if val, ok := item.(*ast.String); ok {
							names = append(names, val.Str)
						} else {
							names = append(names, "")
						}
					}
				}
				rel := &ast.TableName{Name: *cte.Ctename}
				for i := range cols {
					cols[i].Table = rel
					if len(names) > i {
						cols[i].Name = names[i]
					}
				}
				qc.ctes[*cte.Ctename] = &Table{
					Rel:     rel,
					Columns: cols,
				}
			}
		}
	}
	if from != nil {
		for _, item := range from.Items {
			if rs, ok := item.(*ast.RangeSubselect); ok {
				cols, err := comp.outputColumns(qc, rs.Subquery)
				if err != nil {
					return nil, err
				}
				var names []string
				if rs.Alias != nil && rs.Alias.Colnames != nil {
					for _, item := range rs.Alias.Colnames.Items {
						if val, ok := item.(*ast.String); ok {
							names = append(names, val.Str)
						} else {
							names = append(names, "")
						}
					}
				}
				rel := &ast.TableName{}
				if rs.Alias != nil && rs.Alias.Aliasname != nil {
					rel.Name = *rs.Alias.Aliasname
				} else {
					rel.Name = fmt.Sprintf("unaliased_table_%d", rand.Int63())
				}
				for i := range cols {
					cols[i].Table = rel
					if len(names) > i {
						cols[i].Name = names[i]
					}
				}
				qc.fromClauses[rel.Name] = &Table{
					Rel:     rel,
					Columns: cols,
				}
			}
		}
	}
	return qc, nil
}

func ConvertColumn(rel *ast.TableName, c *catalog.Column) *Column {
	return &Column{
		Table:     rel,
		Name:      c.Name,
		DataType:  dataType(&c.Type),
		NotNull:   c.IsNotNull,
		Unsigned:  c.IsUnsigned,
		IsArray:   c.IsArray,
		ArrayDims: c.ArrayDims,
		Type:      &c.Type,
		Length:    c.Length,
	}
}

func (qc QueryCatalog) GetTable(rel *ast.TableName) (*Table, error) {
	cte, exists := qc.ctes[rel.Name]
	if exists {
		return &Table{Rel: rel, Columns: cte.Columns}, nil
	}
	src, err := qc.catalog.GetTable(rel)
	if err != nil {
		return nil, err
	}
	var cols []*Column
	for _, c := range src.Columns {
		cols = append(cols, ConvertColumn(rel, c))
	}
	return &Table{Rel: rel, Columns: cols}, nil
}

func (qc QueryCatalog) GetFunc(rel *ast.FuncName) (*Function, error) {
	funcs, err := qc.catalog.ListFuncsByName(rel)
	if err != nil {
		return nil, err
	}
	if len(funcs) == 0 {
		return nil, fmt.Errorf("function not found: %s", rel.Name)
	}
	return &Function{
		Rel:        rel,
		Outs:       funcs[0].OutArgs(),
		ReturnType: funcs[0].ReturnType,
	}, nil
}
