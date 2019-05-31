package main

import (
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	pg "github.com/lfittl/pg_query_go"
	nodes "github.com/lfittl/pg_query_go/nodes"
)

func columnNames(tree pg.ParsetreeList) (string, []string) {
	table := ""
	cols := []string{}
	for _, stmt := range tree.Statements {
		raw, ok := stmt.(nodes.RawStmt)
		if !ok {
			continue
		}
		switch n := raw.Stmt.(type) {
		case nodes.CreateStmt:
			table = *n.Relation.Relname
			for _, elt := range n.TableElts.Items {
				switch n := elt.(type) {
				case nodes.ColumnDef:
					cols = append(cols, *n.Colname)
				}
			}
			return table, cols
		}
	}
	return table, cols
}

func tableName(stmt nodes.Node) string {
	switch n := stmt.(type) {
	case nodes.SelectStmt:
		for _, item := range n.FromClause.Items {
			switch i := item.(type) {
			case nodes.RangeVar:
				return *i.Relname
			}
		}
	case nodes.InsertStmt:
		return *n.Relation.Relname
	}
	return ""
}

const create = `
CREATE TABLE credentials (
  id         SERIAL       UNIQUE NOT NULL,
  sid        varchar(64)  UNIQUE NOT NULL,
  created    timestamp    DEFAULT NOW(),
  accountid  bigint       NOT NULL,
  tokenhash  varchar(255) NOT NULL
)
`

const load = `
SELECT *
FROM credentials
WHERE created = $1
`

const insert = `
INSERT INTO credentials (
  accountid,
  tokenhash,
  sid
) VALUES (
  $1,
  $2,
  $3
)
RETURNING *
`

func main() {
	tree, err := pg.Parse(create)
	if err != nil {
		panic(err)
	}
	table, columns := columnNames(tree)

	// select
	{
		tree, err := pg.Parse(load)
		if err != nil {
			panic(err)
		}
		for _, stmt := range tree.Statements {
			raw, ok := stmt.(nodes.RawStmt)
			if !ok {
				continue
			}
			spew.Dump(raw)
			if tableName(raw.Stmt) == table {
				fmt.Println(strings.Replace(load, "*", strings.Join(columns, ", "), 1))
			}
		}
	}

	// insert
	{
		tree, err := pg.Parse(insert)
		if err != nil {
			panic(err)
		}
		if false {
			spew.Dump(tree)
		}
		for _, stmt := range tree.Statements {
			raw, ok := stmt.(nodes.RawStmt)
			if !ok {
				continue
			}
			if tableName(raw.Stmt) == table {
				fmt.Println(strings.Replace(insert, "*", strings.Join(columns, ", "), 1))
			}
		}
	}
}
