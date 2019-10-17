package postgres

import (
	"strings"
)

type Schema struct {
	Tables []Table
	Enums  []Enum
}

type Enum struct {
	GoName string
	Name   string
	Vals   []string
}

type Constant struct {
	Name  string
	Type  string
	Value string
}

func (e Enum) Constants() []Constant {
	var c []Constant
	for _, v := range e.Vals {
		name := ""
		for _, part := range strings.Split(strings.Replace(v, "-", "_", -1), "_") {
			name += strings.Title(part)
		}
		c = append(c, Constant{
			Name:  e.GoName + name,
			Value: v,
			Type:  e.GoName,
		})
	}
	return c
}

type Table struct {
	GoName  string
	Name    string
	Columns []Column
}

type Column struct {
	Table   string
	GoName  string
	GoType  string
	Name    string
	Type    string
	NotNull bool
}
