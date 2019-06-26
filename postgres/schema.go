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
		c = append(c, Constant{
			Name:  strings.Title(e.Name) + strings.Title(v),
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
	GoName  string
	GoType  string
	Name    string
	Type    string
	NotNull bool
}
