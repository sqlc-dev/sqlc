package googlesql

import (
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/core"
)

func Dialect() core.Option {
	return core.WithSeed(Seed)
}

type gsqlType struct {
	name     string
	category string
}

var googlesqlTypes = []gsqlType{
	{"bool", "B"},
	{"int64", "N"}, {"float32", "N"}, {"float64", "N"},
	{"numeric", "N"}, {"bignumeric", "N"},
	{"string", "S"}, {"bytes", "S"}, {"uuid", "S"},
	{"date", "D"}, {"datetime", "D"}, {"time", "D"}, {"timestamp", "D"},
	{"interval", "T"},
	{"json", "U"}, {"geography", "U"}, {"struct", "U"}, {"enum", "U"},
	{"tokenlist", "U"},
	{"array", "A"},
}

var comparisonOperators = []string{"=", "<>", "!=", "<", "<=", ">", ">="}

var arithmeticOperators = []string{"+", "-", "*", "/"}

var arithmeticTypes = []string{"int64", "float32", "float64", "numeric", "bignumeric"}

func Seed(cat *core.Catalog) error {
	dialectOID, err := cat.CreateDialect("googlesql")
	if err != nil {
		return err
	}

	oids := make(map[string]int64, len(googlesqlTypes))
	for _, t := range googlesqlTypes {
		oid, err := cat.CreateTypeSpec(core.TypeSpec{
			Name:       t.name,
			Typtype:    "b",
			Category:   t.category,
			DialectOID: dialectOID,
		})
		if err != nil {
			return fmt.Errorf("seed googlesql type %q: %w", t.name, err)
		}
		oids[t.name] = oid
	}

	boolOID := oids["bool"]
	for _, t := range googlesqlTypes {
		if t.name == "array" || t.name == "struct" {
			continue
		}
		for _, op := range comparisonOperators {
			if _, err := cat.CreateOperator(core.OperatorSpec{
				Name:          op,
				DialectOID:    dialectOID,
				LeftTypeOID:   oids[t.name],
				RightTypeOID:  oids[t.name],
				ResultTypeOID: boolOID,
			}); err != nil {
				return fmt.Errorf("seed googlesql operator %q(%s): %w", op, t.name, err)
			}
		}
	}

	for _, name := range arithmeticTypes {
		for _, op := range arithmeticOperators {
			if _, err := cat.CreateOperator(core.OperatorSpec{
				Name:          op,
				DialectOID:    dialectOID,
				LeftTypeOID:   oids[name],
				RightTypeOID:  oids[name],
				ResultTypeOID: oids[name],
			}); err != nil {
				return fmt.Errorf("seed googlesql operator %q(%s): %w", op, name, err)
			}
		}
	}

	if _, err := cat.CreateProc(core.ProcSpec{
		Name:          "count",
		DialectOID:    dialectOID,
		Kind:          "a",
		ReturnTypeOID: oids["int64"],
	}); err != nil {
		return fmt.Errorf("seed googlesql function %q: %w", "count", err)
	}

	if _, err := cat.CreateOperator(core.OperatorSpec{
		Name:          "||",
		DialectOID:    dialectOID,
		LeftTypeOID:   oids["string"],
		RightTypeOID:  oids["string"],
		ResultTypeOID: oids["string"],
	}); err != nil {
		return fmt.Errorf("seed googlesql operator %q: %w", "||", err)
	}

	return nil
}
