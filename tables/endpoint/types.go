package endpoint

import (
	"strings"

	sq "github.com/Masterminds/squirrel"
)

type Column interface {
	name() string
}

type column string

func (c column) name() string {
	return string(c)
}

const (
	ID        column = "id"
	AccountID column = "account_id"
	Settings  column = "settings"
)

var Arg = column("?")

type Columns []Column

type Comparable interface {
	pred() interface{}
}
type And []Comparable

func (a And) pred() interface{} {
	eq := sq.Eq{}
	for _, c := range a {
		switch cmp := c.(type) {
		case Eq:
			eq[cmp.First.name()] = cmp.Second.name()
		}
	}
	return eq
}

type Or []Comparable

func (o Or) pred() interface{} {
	return nil
}

type Eq struct {
	First  Column
	Second Column
}

func (e Eq) pred() interface{} {
	eq := sq.Eq{}
	eq[e.First.name()] = e.Second.name()
	return eq
}

// type NotEq struct {
// 	First  Column
// 	Second Column
// }
//
// func (e NotEq) pred() interface{} {
// }

type Select struct {
	Columns []Column
	Where   Comparable
}

func (s Select) String() string {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	cs := []string{}
	cols := s.Columns
	if len(cols) == 0 {
		cols = Star
	}
	for _, c := range cols {
		cs = append(cs, c.name())
	}
	from := psql.Select(cs...).From("endpoint")
	sql, _, _ := from.Where(s.Where.pred()).ToSql()
	return sql
}

type Update struct {
	Set       []Column
	Where     Comparable
	Returning []Column
}

type Insert struct {
	Columns   []Column
	Returning []Column
}

func (i Insert) String() string {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	cols := i.Columns
	if len(cols) == 0 {
		cols = Star
	}

	names := []string{}
	values := []interface{}{}
	for _, c := range cols {
		names = append(names, c.name())
		values = append(values, "?")
	}

	ins := psql.Insert("endpoint").Columns(names...).Values(values...)

	if len(i.Returning) != 0 {
		cs := []string{}
		for _, c := range i.Returning {
			cs = append(cs, c.name())
		}
		ins = ins.Suffix("RETURNING " + strings.Join(cs, ", "))
	}

	sql, _, _ := ins.ToSql()
	return sql
}
