package sqlite

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/goldeneye/analysis"
)

// catalog is what the database says about a schema: its tables and their
// columns, its indexes, and the b-tree each is stored in, which is how the
// bytecode names them.
type catalog struct {
	tables map[string]*table
	roots  map[int]object // by root page
}

type table struct {
	name         string
	withoutRowid bool
	cols         []*tableColumn // by cid
	// stored lists the columns in the order a row's record holds them,
	// which is cid order less the generated columns that are not stored.
	stored []*tableColumn
	// pk is the PRIMARY KEY index of a WITHOUT ROWID table, whose order is
	// the order the table's own b-tree stores a row in.
	pk *index
}

type tableColumn struct {
	table    *table
	cid      int
	name     string
	declType string
	notNull  bool
	pk       int // position in the primary key, 0 for none
}

type index struct {
	name  string
	table *table
	cids  []int // by seqno: -1 for the rowid, -2 for an expression
}

// object is what a cursor is open on: a table, an index, or, with neither
// set, something the catalog does not describe — an ephemeral table, a
// sorter, a virtual table.
type object struct {
	table *table
	index *index
}

// The catalog is read with these, each in a section of its own, in json
// mode. A query that finds nothing prints nothing, so a section may be
// empty.
var catalogQueries = []struct{ section, sql string }{
	{"tables", `SELECT name, type, wr FROM pragma_table_list WHERE schema = 'main' ORDER BY name`},
	{"columns", `SELECT t.name AS tbl, c.cid, c.name, c.type, c."notnull" AS "notnull", c.pk, c.hidden
	 FROM pragma_table_list t, pragma_table_xinfo(t.name) c
	 WHERE t.schema = 'main' AND t.type IN ('table', 'virtual') ORDER BY t.name, c.cid`},
	{"indexes", `SELECT t.name AS tbl, i.name AS idx, i.origin, x.seqno, x.cid
	 FROM pragma_table_list t, pragma_index_list(t.name) i, pragma_index_xinfo(i.name) x
	 WHERE t.schema = 'main' ORDER BY t.name, i.name, x.seqno`},
	{"roots", `SELECT type, name, tbl_name, rootpage FROM sqlite_schema WHERE rootpage > 0`},
}

// readCatalog builds the catalog from the sections catalogQueries printed.
func readCatalog(out *output) (*catalog, error) {
	var tables []struct {
		Name string `json:"name"`
		Type string `json:"type"`
		WR   int    `json:"wr"`
	}
	var columns []struct {
		Table   string `json:"tbl"`
		CID     int    `json:"cid"`
		Name    string `json:"name"`
		Type    string `json:"type"`
		NotNull int    `json:"notnull"`
		PK      int    `json:"pk"`
		Hidden  int    `json:"hidden"`
	}
	var indexes []struct {
		Table  string `json:"tbl"`
		Index  string `json:"idx"`
		Origin string `json:"origin"`
		Seqno  int    `json:"seqno"`
		CID    int    `json:"cid"`
	}
	var roots []struct {
		Type     string `json:"type"`
		Name     string `json:"name"`
		Table    string `json:"tbl_name"`
		Rootpage int    `json:"rootpage"`
	}
	for name, v := range map[string]any{"tables": &tables, "columns": &columns, "indexes": &indexes, "roots": &roots} {
		s := out.sections[name]
		if s == nil {
			return nil, fmt.Errorf("reading the catalog: no %s section in the output", name)
		}
		if len(s.blocks) > 0 {
			if err := s.decode(0, v); err != nil {
				return nil, fmt.Errorf("reading the catalog's %s: %w", name, err)
			}
		}
	}

	c := &catalog{tables: map[string]*table{}, roots: map[int]object{}}
	for _, t := range tables {
		if t.Type == "table" || t.Type == "virtual" {
			c.tables[t.Name] = &table{name: t.Name, withoutRowid: t.WR == 1}
		}
	}
	for _, col := range columns {
		t := c.tables[col.Table]
		if t == nil {
			continue
		}
		tc := &tableColumn{table: t, cid: col.CID, name: col.Name, declType: col.Type, notNull: col.NotNull == 1, pk: col.PK}
		for len(t.cols) <= col.CID {
			t.cols = append(t.cols, nil)
		}
		t.cols[col.CID] = tc
		// A generated column that is not stored — hidden is 2 — has no
		// place in the record.
		if col.Hidden != 2 {
			t.stored = append(t.stored, tc)
		}
	}
	byName := map[string]*index{}
	for _, ic := range indexes {
		t := c.tables[ic.Table]
		if t == nil {
			continue
		}
		ix := byName[ic.Index]
		if ix == nil {
			ix = &index{name: ic.Index, table: t}
			byName[ic.Index] = ix
			if ic.Origin == "pk" && t.withoutRowid {
				t.pk = ix
			}
		}
		for len(ix.cids) <= ic.Seqno {
			ix.cids = append(ix.cids, -2)
		}
		ix.cids[ic.Seqno] = ic.CID
	}
	for _, r := range roots {
		switch r.Type {
		case "table":
			if t := c.tables[r.Name]; t != nil {
				c.roots[r.Rootpage] = object{table: t}
			}
		case "index":
			if ix := byName[r.Name]; ix != nil {
				c.roots[r.Rootpage] = object{index: ix}
			}
		}
	}
	return c, nil
}

// root returns what is stored in a b-tree of the main database.
func (c *catalog) root(db, page int) object {
	if db != 0 {
		return object{}
	}
	return c.roots[page]
}

// lookup finds a table column by the names `.stats stmt` prints.
func (c *catalog) lookup(tbl, col string) *tableColumn {
	t := c.tables[tbl]
	if t == nil {
		return nil
	}
	for _, tc := range t.cols {
		if tc != nil && tc.name == col {
			return tc
		}
	}
	return nil
}

// rowidAlias returns the column that is another name for the rowid: the
// lone INTEGER PRIMARY KEY of a table that has a rowid.
func (t *table) rowidAlias() *tableColumn {
	if t.withoutRowid {
		return nil
	}
	var alias *tableColumn
	for _, tc := range t.cols {
		if tc == nil || tc.pk == 0 {
			continue
		}
		if alias != nil || !strings.EqualFold(tc.declType, "INTEGER") {
			return nil
		}
		alias = tc
	}
	return alias
}

// describe is the column as sqlc analyze describes one: its declared type,
// nullable unless declared NOT NULL or the rowid, which is never NULL.
func (tc *tableColumn) describe() analysis.Column {
	typ := parseType(tc.declType)
	typ.Nullable = !tc.notNull && tc.table.rowidAlias() != tc
	return analysis.Column{Name: tc.name, Type: typ, Table: tc.table.name}
}

// rowid describes the table's rowid: the column that aliases it, or the
// rowid itself.
func (t *table) rowid() analysis.Column {
	if alias := t.rowidAlias(); alias != nil {
		return alias.describe()
	}
	return analysis.Column{Name: "rowid", Type: &analysis.TypeExpr{Name: "integer"}, Table: t.name}
}

// column returns what the object's ith stored column is: a table column,
// or the rowid of a table. An index stores the columns it is on and then
// the rowid; a WITHOUT ROWID table is stored in the order of its primary
// key; any other table is stored in column order.
func (o object) column(i int) (*tableColumn, *table, bool) {
	ix := o.index
	if o.table != nil && o.table.withoutRowid {
		ix = o.table.pk
	}
	switch {
	case ix != nil:
		if i < 0 || i >= len(ix.cids) {
			return nil, nil, false
		}
		cid := ix.cids[i]
		switch {
		case cid == -1:
			return nil, ix.table, true
		case cid >= 0 && cid < len(ix.table.cols) && ix.table.cols[cid] != nil:
			return ix.table.cols[cid], nil, true
		}
	case o.table != nil:
		if i >= 0 && i < len(o.table.stored) {
			return o.table.stored[i], nil, true
		}
	}
	return nil, nil, false
}

// owner is the table the object stores rows of.
func (o object) owner() *table {
	if o.index != nil {
		return o.index.table
	}
	return o.table
}

// parseType reads a declared type the way sqlc's catalog does: the name
// lowercased, with whatever is in parentheses after it as arguments. A
// column declared with no type at all can hold anything.
func parseType(decl string) *analysis.TypeExpr {
	decl = strings.TrimSpace(decl)
	if decl == "" {
		return &analysis.TypeExpr{Name: "any"}
	}
	name, args := decl, ""
	if open := strings.IndexByte(decl, '('); open >= 0 && strings.HasSuffix(decl, ")") {
		name, args = decl[:open], decl[open+1:len(decl)-1]
	}
	t := &analysis.TypeExpr{Name: strings.ToLower(strings.TrimSpace(name))}
	if strings.TrimSpace(args) == "" {
		return t
	}
	for _, a := range strings.Split(args, ",") {
		a = strings.TrimSpace(a)
		switch {
		case strings.HasPrefix(a, "'") && strings.HasSuffix(a, "'") && len(a) >= 2:
			s := strings.ReplaceAll(a[1:len(a)-1], "''", "'")
			t.Args = append(t.Args, analysis.TypeArg{String: &s})
		case strings.EqualFold(a, "true") || strings.EqualFold(a, "false"):
			b := strings.EqualFold(a, "true")
			t.Args = append(t.Args, analysis.TypeArg{Bool: &b})
		default:
			if n, err := strconv.ParseInt(a, 10, 64); err == nil {
				t.Args = append(t.Args, analysis.TypeArg{Int: &n})
			} else {
				t.Args = append(t.Args, analysis.TypeArg{Type: &analysis.TypeExpr{Name: strings.ToLower(a)}})
			}
		}
	}
	return t
}
