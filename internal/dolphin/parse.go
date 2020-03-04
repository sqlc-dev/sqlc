package dolphin

import (
	"io"
	"io/ioutil"

	"github.com/kyleconroy/sqlc/internal/sql/ast"

	"github.com/pingcap/parser"
	pcast "github.com/pingcap/parser/ast"
	_ "github.com/pingcap/tidb/types/parser_driver"
)

func NewParser() *Parser {
	return &Parser{parser.New()}
}

type Parser struct {
	pingcap *parser.Parser
}

type nodeSearch struct {
	list  []pcast.Node
	check func(pcast.Node) bool
}

func (s *nodeSearch) Enter(n pcast.Node) (pcast.Node, bool) {
	if s.check(n) {
		s.list = append(s.list, n)
	}
	return n, false // skipChildren
}

func (s *nodeSearch) Leave(n pcast.Node) (pcast.Node, bool) {
	return n, true // ok
}

func collect(root pcast.Node, f func(pcast.Node) bool) []pcast.Node {
	if root == nil {
		return nil
	}
	ns := &nodeSearch{check: f}
	root.Accept(ns)
	return ns.list
}

type nodeVisit struct {
	fn func(pcast.Node)
}

func (s *nodeVisit) Enter(n pcast.Node) (pcast.Node, bool) {
	s.fn(n)
	return n, false // skipChildren
}

func (s *nodeVisit) Leave(n pcast.Node) (pcast.Node, bool) {
	return n, true // ok
}

func visit(root pcast.Node, f func(pcast.Node)) {
	if root == nil {
		return
	}
	ns := &nodeVisit{fn: f}
	root.Accept(ns)
}

// Maybe not useful?
func text(nodes []pcast.Node) []string {
	str := make([]string, len(nodes))
	for i := range nodes {
		if nodes[i] == nil {
			continue
		}
		str[i] = nodes[i].Text()
	}
	return str
}

func parseTableName(n *pcast.TableName) *ast.TableName {
	return &ast.TableName{
		Schema: n.Schema.String(),
		Name:   n.Name.String(),
	}
}

func (p *Parser) Parse(r io.Reader) ([]ast.Statement, error) {
	blob, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	stmtNodes, _, err := p.pingcap.Parse(string(blob), "", "")
	if err != nil {
		return nil, err
	}
	var stmts []ast.Statement
	for i := range stmtNodes {
		var stmt ast.Node
		switch n := stmtNodes[i].(type) {

		case *pcast.AlterTableStmt:
			alt := &ast.AlterTableStmt{
				Table: parseTableName(n.Table),
				Cmds:  &ast.List{},
			}
			for _, spec := range n.Specs {
				switch spec.Tp {
				case pcast.AlterTableAddColumns:
					for _, def := range spec.NewColumns {
						name := def.Name.String()
						alt.Cmds.Items = append(alt.Cmds.Items, &ast.AlterTableCmd{
							Name:    &name,
							Subtype: ast.AT_AddColumn,
							Def: &ast.ColumnDef{
								Colname: def.Name.String(),
								// TODO: Use def.Tp to generate type name
								TypeName: &ast.TypeName{Name: "text"},
							},
						})
					}

				case pcast.AlterTableDropColumn:
					name := spec.OldColumnName.String()
					alt.Cmds.Items = append(alt.Cmds.Items, &ast.AlterTableCmd{
						Name:      &name,
						Subtype:   ast.AT_DropColumn,
						MissingOk: spec.IfExists,
					})

				case pcast.AlterTableChangeColumn:
					// 	spew.Dump("change column", spec)

				case pcast.AlterTableModifyColumn:
					// 	spew.Dump("modify column", spec)

				case pcast.AlterTableAlterColumn:
					// 	spew.Dump("alter column", spec)

				case pcast.AlterTableAddConstraint:
					// 	spew.Dump("add const", spec)

				default:
					continue
				}
			}
			stmt = alt

		case *pcast.CreateTableStmt:
			create := &ast.CreateTableStmt{
				Name:        parseTableName(n.Table),
				IfNotExists: n.IfNotExists,
			}
			for _, def := range n.Cols {
				create.Cols = append(create.Cols, &ast.ColumnDef{
					Colname: def.Name.String(),
					// TODO: Use n.Tp to generate type name
					TypeName: &ast.TypeName{Name: "text"},
				})
			}
			stmt = create

		case *pcast.DropTableStmt:
			drop := &ast.DropTableStmt{IfExists: n.IfExists}
			for _, name := range n.Tables {
				drop.Tables = append(drop.Tables, parseTableName(name))
			}
			stmt = drop

		case *pcast.SelectStmt:
			sel := &ast.SelectStmt{}
			var tables []ast.Node
			visit(n.From, func(n pcast.Node) {
				name, ok := n.(*pcast.TableName)
				if !ok {
					return
				}
				tables = append(tables, parseTableName(name))
			})
			var cols []ast.Node
			visit(n.Fields, func(n pcast.Node) {
				col, ok := n.(*pcast.ColumnName)
				if !ok {
					return
				}
				cols = append(cols, &ast.ResTarget{
					Val: &ast.ColumnRef{
						Name: col.Name.String(),
					},
				})
			})
			sel.From = &ast.List{Items: tables}
			sel.Fields = &ast.List{Items: cols}
			stmt = sel

		default:
			// spew.Dump(n)

		}

		if stmt != nil {
			stmts = append(stmts, ast.Statement{
				Raw: &ast.RawStmt{Stmt: stmt},
			})
		}
	}
	return stmts, nil
}
