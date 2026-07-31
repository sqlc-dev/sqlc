// Package seed builds a core catalog for a SQL dialect from a declarative
// description of its type system: the types it defines, the operators and
// casts between them, and the functions it ships with.
package seed

import (
	"fmt"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

// Type is a type the dialect defines. Aliases are spellings of the same type
// that a schema may use in a column definition; each becomes its own catalog
// type, with implicit casts registered between them.
type Type struct {
	Name     string
	Category string
	Aliases  []string
}

// Operator is a single operator overload.
type Operator struct {
	Name   string
	Left   string
	Right  string
	Result string
}

// Cast is a single cast between two types. Context is 'i'mplicit, 'a'ssignment
// or 'e'xplicit.
type Cast struct {
	Source  string
	Target  string
	Context string
}

// Spec describes a dialect. Apply turns it into catalog rows.
type Spec struct {
	// Dialect is the name recorded in the catalog, e.g. "postgresql".
	Dialect string

	// Types are the dialect's types.
	Types []Type

	// Const names the type each kind of literal takes on, keyed by the
	// core.Const* constants. Kinds left out fall back to PostgreSQL's names.
	Const map[string]string

	// Comparison operators are registered as (T, T) -> Bool for every type in
	// ComparisonCategories.
	Comparison           []string
	ComparisonCategories string

	// Arithmetic operators are registered as (T, T) -> T for every type in
	// ArithmeticCategories.
	Arithmetic           []string
	ArithmeticCategories string

	// Bool names the type comparisons return.
	Bool string

	// Operators and Casts are registered as given, on top of the ones the
	// category rules above generate.
	Operators []Operator
	Casts     []Cast

	// CastCategories lists the categories whose types are all implicitly
	// castable to one another, so that an operator on two spellings of the
	// same kind of value resolves. "*" makes every seeded type implicitly
	// castable to every other, for dialects that compare across categories.
	CastCategories string

	// Funcs are the dialect's built-in functions, in the form the engine
	// packages already describe them.
	Funcs []*catalog.Function
}

// Apply registers everything the spec describes on cat.
func (s Spec) Apply(cat *core.Catalog) error {
	b := &builder{cat: cat, oids: map[string]int64{}, seenCasts: map[[2]int64]bool{}}

	dialectOID, err := cat.CreateDialect(s.Dialect)
	if err != nil {
		return err
	}
	b.dialectOID = dialectOID

	if err := b.types(s); err != nil {
		return err
	}
	if err := b.consts(s); err != nil {
		return err
	}
	if err := b.categoryOperators(s); err != nil {
		return err
	}
	if err := b.operators(s); err != nil {
		return err
	}
	if err := b.casts(s); err != nil {
		return err
	}
	return b.funcs(s)
}

type builder struct {
	cat        *core.Catalog
	dialectOID int64

	// oids maps a lowercased type name to its OID, and categories a type name
	// to the category it was seeded under.
	oids       map[string]int64
	categories []categorized

	// seenCasts records the pairs already registered. The alias rules and the
	// category rules overlap, and a cast pair is unique in the catalog.
	seenCasts map[[2]int64]bool
}

type categorized struct {
	name     string
	category string
}

func (b *builder) types(s Spec) error {
	for _, t := range s.Types {
		for _, name := range append([]string{t.Name}, t.Aliases...) {
			if _, err := b.createType(name, t.Category); err != nil {
				return fmt.Errorf("seed %s type %q: %w", s.Dialect, name, err)
			}
		}
		if err := b.aliasCasts(t); err != nil {
			return err
		}
	}
	return nil
}

// aliasCasts makes every spelling of a type implicitly castable to every other,
// so that a column declared "integer" and one declared "int4" compare.
func (b *builder) aliasCasts(t Type) error {
	names := append([]string{t.Name}, t.Aliases...)
	for _, src := range names {
		for _, tgt := range names {
			if src == tgt {
				continue
			}
			if err := b.cast(Cast{Source: src, Target: tgt, Context: "i"}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *builder) createType(name, category string) (int64, error) {
	key := strings.ToLower(name)
	if oid, ok := b.oids[key]; ok {
		return oid, nil
	}
	oid, err := b.cat.CreateTypeSpec(core.TypeSpec{
		Name:       key,
		Typtype:    "b",
		Category:   category,
		DialectOID: b.dialectOID,
	})
	if err != nil {
		return 0, err
	}
	b.oids[key] = oid
	b.categories = append(b.categories, categorized{name: key, category: category})
	return oid, nil
}

func (b *builder) consts(s Spec) error {
	for kind, name := range s.Const {
		if _, ok := b.oids[strings.ToLower(name)]; !ok {
			return fmt.Errorf("seed %s: constant %s names unknown type %q", s.Dialect, kind, name)
		}
		if err := b.cat.SetConstType(b.dialectOID, kind, name); err != nil {
			return err
		}
	}
	// A schema declares types the seed knows nothing about — enums, domains,
	// arrays, a SQLite column typed whatever the author felt like. Recording
	// the comparison operators lets the catalog give those types the same ones.
	if len(s.Comparison) > 0 {
		if err := b.cat.SetDialectFlag(b.dialectOID, core.FlagComparisonOperators, strings.Join(s.Comparison, ",")); err != nil {
			return err
		}
	}
	return nil
}

func (b *builder) categoryOperators(s Spec) error {
	boolOID, ok := b.oids[strings.ToLower(s.Bool)]
	if len(s.Comparison) > 0 && !ok {
		return fmt.Errorf("seed %s: comparison operators need a bool type", s.Dialect)
	}
	for _, t := range b.categories {
		if strings.Contains(s.ComparisonCategories, t.category) {
			for _, op := range s.Comparison {
				if err := b.operator(op, t.name, t.name, s.Bool, boolOID); err != nil {
					return err
				}
			}
		}
		if strings.Contains(s.ArithmeticCategories, t.category) {
			for _, op := range s.Arithmetic {
				if err := b.operator(op, t.name, t.name, t.name, b.oids[t.name]); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (b *builder) operator(name, left, right, result string, resultOID int64) error {
	leftOID, ok := b.oids[strings.ToLower(left)]
	if !ok {
		return fmt.Errorf("operator %q: unknown left type %q", name, left)
	}
	rightOID, ok := b.oids[strings.ToLower(right)]
	if !ok {
		return fmt.Errorf("operator %q: unknown right type %q", name, right)
	}
	if resultOID == 0 {
		resultOID, ok = b.oids[strings.ToLower(result)]
		if !ok {
			return fmt.Errorf("operator %q: unknown result type %q", name, result)
		}
	}
	_, err := b.cat.CreateOperator(core.OperatorSpec{
		Name:          name,
		DialectOID:    b.dialectOID,
		LeftTypeOID:   leftOID,
		RightTypeOID:  rightOID,
		ResultTypeOID: resultOID,
	})
	return err
}

func (b *builder) operators(s Spec) error {
	for _, op := range s.Operators {
		if err := b.operator(op.Name, op.Left, op.Right, op.Result, 0); err != nil {
			return fmt.Errorf("seed %s: %w", s.Dialect, err)
		}
	}
	return nil
}

func (b *builder) casts(s Spec) error {
	for _, c := range s.Casts {
		if err := b.cast(c); err != nil {
			return fmt.Errorf("seed %s: %w", s.Dialect, err)
		}
	}
	if s.CastCategories == "" {
		return nil
	}
	anyCategory := s.CastCategories == "*"
	for _, src := range b.categories {
		if !anyCategory && !strings.Contains(s.CastCategories, src.category) {
			continue
		}
		for _, tgt := range b.categories {
			if src.name == tgt.name {
				continue
			}
			if !anyCategory && src.category != tgt.category {
				continue
			}
			if err := b.cast(Cast{Source: src.name, Target: tgt.name, Context: "i"}); err != nil {
				return fmt.Errorf("seed %s: %w", s.Dialect, err)
			}
		}
	}
	return nil
}

func (b *builder) cast(c Cast) error {
	srcOID, ok := b.oids[strings.ToLower(c.Source)]
	if !ok {
		return fmt.Errorf("cast: unknown source type %q", c.Source)
	}
	tgtOID, ok := b.oids[strings.ToLower(c.Target)]
	if !ok {
		return fmt.Errorf("cast: unknown target type %q", c.Target)
	}
	if srcOID == tgtOID || b.seenCasts[[2]int64{srcOID, tgtOID}] {
		return nil
	}
	b.seenCasts[[2]int64{srcOID, tgtOID}] = true
	context := c.Context
	if context == "" {
		context = "i"
	}
	return b.cat.CreateCast(core.CastSpec{
		SourceTypeOID: srcOID,
		TargetTypeOID: tgtOID,
		Context:       context,
		DialectOID:    b.dialectOID,
	})
}

func (b *builder) funcs(s Spec) error {
	for _, fn := range s.Funcs {
		returnOID, err := b.funcType(fn.ReturnType)
		if err != nil {
			return fmt.Errorf("seed %s function %q: %w", s.Dialect, fn.Name, err)
		}
		args := make([]core.ProcArg, 0, len(fn.InArgs()))
		for _, arg := range fn.InArgs() {
			argOID, err := b.funcType(arg.Type)
			if err != nil {
				return fmt.Errorf("seed %s function %q: %w", s.Dialect, fn.Name, err)
			}
			args = append(args, core.ProcArg{
				Name:       arg.Name,
				TypeOID:    argOID,
				HasDefault: arg.HasDefault,
			})
		}
		if _, err := b.cat.CreateProc(core.ProcSpec{
			Name:           fn.Name,
			DialectOID:     b.dialectOID,
			ReturnTypeOID:  returnOID,
			ReturnNullable: fn.ReturnTypeNullable,
			Args:           args,
		}); err != nil {
			return fmt.Errorf("seed %s function %q: %w", s.Dialect, fn.Name, err)
		}
	}
	return nil
}

// funcType resolves the type a function signature names. Signatures reference
// pseudo types ("any", "record") and types no dialect spec bothers to list, so
// an unknown name is registered rather than rejected.
func (b *builder) funcType(tn *ast.TypeName) (int64, error) {
	if tn == nil || tn.Name == "" {
		return 0, nil
	}
	return b.createType(tn.Name, "U")
}
