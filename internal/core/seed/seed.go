// Package seed builds a core catalog for a SQL dialect from a declarative
// description of its type system: the types it defines, the operators and
// casts between them, and the functions it ships with.
//
// A dialect is a directory which each engine package embeds and hands to
// Dialect:
//
//	dialect.json     what the dialect is called and the rules it is built by
//	types.jsonl      the types it defines
//	operators.jsonl  operator overloads beyond the ones the rules generate
//	casts.jsonl      casts beyond the ones the rules generate
//	functions.jsonl  the functions it ships with
//	relations.jsonl  the system tables and views it ships with
//
// The lists are JSONL — one record per line — and are applied as they are
// read, so a dialect whose function list runs to thousands of entries is never
// held in memory as a whole. Any of the lists may be left out.
//
// A dialect may also hold an extensions/ directory with one directory per
// extension, each a smaller bundle of the same files, applied when a schema
// says CREATE EXTENSION.
package seed

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"slices"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

// The files a dialect directory is made of.
const (
	SettingsFile  = "dialect.json"
	TypesFile     = "types.jsonl"
	OperatorsFile = "operators.jsonl"
	CastsFile     = "casts.jsonl"
	FunctionsFile = "functions.jsonl"
	RelationsFile = "relations.jsonl"
)

// ExtensionsDir is the directory under a dialect holding one directory per
// extension the dialect knows.
const ExtensionsDir = "extensions"

// Settings is dialect.json: what the dialect is called and the rules that
// generate its operators and casts.
type Settings struct {
	// Dialect is the name recorded in the catalog, e.g. "postgresql".
	Dialect string `json:"dialect"`

	// Const names the type each kind of literal takes on, keyed by the
	// core.Const* constants. Kinds left out fall back to PostgreSQL's names.
	Const map[string]string `json:"const,omitempty"`

	// Bool names the type comparisons return.
	Bool string `json:"bool,omitempty"`

	// Limit names the type a LIMIT or OFFSET count has, which is what a
	// placeholder in one is typed as. Unset, it is the integer literal's.
	Limit string `json:"limit,omitempty"`

	// Untyped names the type a placeholder takes when nothing in the query
	// constrains it. Unset, such a placeholder stays untyped.
	Untyped string `json:"untyped,omitempty"`

	// PropagateNullable makes a function's result nullable whenever one of
	// its arguments is, the way ClickHouse's ordinary functions behave,
	// unless the function is seeded as never null.
	PropagateNullable bool `json:"propagate_nullable,omitempty"`

	// QualifyDuplicateColumns names a result column after its relation when
	// an earlier result column from another relation has the same name, as
	// ClickHouse names the second id of a join e.id.
	QualifyDuplicateColumns bool `json:"qualify_duplicate_columns,omitempty"`

	// Comparison operators are registered as (T, T) -> Bool for every type in
	// ComparisonCategories.
	Comparison           []string `json:"comparison,omitempty"`
	ComparisonCategories string   `json:"comparison_categories,omitempty"`

	// Arithmetic operators are registered as (T, T) -> T for every type in
	// ArithmeticCategories.
	Arithmetic           []string `json:"arithmetic,omitempty"`
	ArithmeticCategories string   `json:"arithmetic_categories,omitempty"`

	// CastCategories lists the categories whose types are all implicitly
	// castable to one another, so that an operator on two spellings of the
	// same kind of value resolves. "*" makes every seeded type implicitly
	// castable to every other, for dialects that compare across categories.
	CastCategories string `json:"cast_categories,omitempty"`
}

// Type is a type the dialect defines. Aliases are spellings of the same type
// that a schema may use in a column definition; each becomes its own catalog
// type, with implicit casts registered between them.
type Type struct {
	Name     string   `json:"name"`
	Category string   `json:"category"`
	Aliases  []string `json:"aliases,omitempty"`
}

// Operator is a single operator overload.
type Operator struct {
	Name   string `json:"name"`
	Left   string `json:"left"`
	Right  string `json:"right"`
	Result string `json:"result"`
}

// Cast is a single cast between two types. Context is 'i'mplicit, 'a'ssignment
// or 'e'xplicit.
type Cast struct {
	Source  string `json:"source"`
	Target  string `json:"target"`
	Context string `json:"context,omitempty"`
}

// Function is a function the dialect ships with. Kind is 'f'unction,
// 'a'ggregate, 'w'indow or 'p'rocedure.
type Function struct {
	Name string `json:"name,omitempty"`
	Kind string `json:"kind,omitempty"`
	Args []Arg  `json:"args,omitempty"`
	// Returns names the result type, or "$1", "$2"... for the type of that
	// argument.
	Returns  string `json:"returns"`
	Nullable bool   `json:"nullable,omitempty"`
	// NeverNull marks a result that is never NULL even when an argument
	// is, in a dialect that propagates nullability.
	NeverNull bool `json:"never_null,omitempty"`
}

// Relation is a table or view the dialect ships with, such as one of
// PostgreSQL's system catalogs. Kind is 'r' for a table or 'v' for a view, and
// defaults to a table.
type Relation struct {
	// Catalog is the database the relation belongs to, which PostgreSQL
	// reports for its own schemas and which the legacy catalog carries.
	Catalog string   `json:"catalog,omitempty"`
	Schema  string   `json:"schema"`
	Name    string   `json:"name"`
	Kind    string   `json:"kind,omitempty"`
	Columns []Column `json:"columns"`
}

// Column is one of a relation's columns.
type Column struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	NotNull bool   `json:"not_null,omitempty"`
	Array   bool   `json:"array,omitempty"`
	Length  int    `json:"length,omitempty"`
}

// Arg is one of a function's parameters. Mode is 'i'n, 'o'ut, 'b'oth, 't'able
// or 'v'ariadic, and defaults to in.
type Arg struct {
	Name       string `json:"name,omitempty"`
	Type       string `json:"type"`
	Mode       string `json:"mode,omitempty"`
	HasDefault bool   `json:"has_default,omitempty"`
}

// Dialect returns the catalog option that seeds the dialect described by the
// JSONL files in dir. A dialect that ships extension data — a directory per
// extension under extensions/ — also has CREATE EXTENSION wired up to it.
func Dialect(fsys fs.FS, dir string) core.Option {
	return core.WithSeed(func(cat *core.Catalog) error {
		sub, err := fs.Sub(fsys, dir)
		if err != nil {
			return fmt.Errorf("seed: %s: %w", dir, err)
		}
		if err := apply(cat, sub); err != nil {
			return err
		}
		cat.SetExtensionLoader(func(name string) error {
			return applyExtension(cat, sub, name)
		})
		return nil
	})
}

func apply(cat *core.Catalog, fsys fs.FS) error {
	settings, err := loadSettings(fsys)
	if err != nil {
		return err
	}
	b := &builder{
		cat:        cat,
		settings:   settings,
		oids:       map[string]int64{},
		namespaces: map[string]int64{},
		seenCasts:  map[[2]int64]bool{},
	}
	if b.dialectOID, err = cat.CreateDialect(settings.Dialect); err != nil {
		return err
	}

	// Types come first: everything else names them.
	if err := stream(fsys, TypesFile, b.addType); err != nil {
		return err
	}
	if err := b.consts(); err != nil {
		return err
	}
	if err := b.categoryOperators(); err != nil {
		return err
	}
	if err := stream(fsys, OperatorsFile, b.addOperator); err != nil {
		return err
	}
	if err := stream(fsys, CastsFile, b.addCast); err != nil {
		return err
	}
	if err := b.categoryCasts(); err != nil {
		return err
	}
	if err := stream(fsys, FunctionsFile, b.addFunction); err != nil {
		return err
	}
	return stream(fsys, RelationsFile, b.addRelation)
}

func loadSettings(fsys fs.FS) (Settings, error) {
	f, err := fsys.Open(SettingsFile)
	if err != nil {
		return Settings{}, fmt.Errorf("seed: %w", err)
	}
	defer f.Close()

	dec := json.NewDecoder(bufio.NewReader(f))
	dec.DisallowUnknownFields()
	var settings Settings
	if err := dec.Decode(&settings); err != nil {
		return Settings{}, fmt.Errorf("seed: %s: %w", SettingsFile, err)
	}
	if settings.Dialect == "" {
		return Settings{}, fmt.Errorf("seed: %s: dialect has no name", SettingsFile)
	}
	return settings, nil
}

// stream decodes name a record at a time, handing each to fn as it is read. A
// file that is not there contributes nothing, so a dialect only writes the
// files it needs. An unknown field is an error rather than something to
// ignore, so that a misspelled key is caught where it is written.
func stream[T any](fsys fs.FS, name string, fn func(T) error) error {
	f, err := fsys.Open(name)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("seed: %w", err)
	}
	defer f.Close()

	dec := json.NewDecoder(bufio.NewReader(f))
	dec.DisallowUnknownFields()
	for record := 1; ; record++ {
		var rec T
		if err := dec.Decode(&rec); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("seed: %s: record %d: %w", name, record, err)
		}
		if err := fn(rec); err != nil {
			return fmt.Errorf("seed: %s: record %d: %w", name, record, err)
		}
	}
}

// Functions streams a dialect's functions into the form the engine catalogs
// use. It is how an engine whose standard library lives in functions.jsonl
// hands that library to the rest of sqlc.
func Functions(fsys fs.FS, dir string) ([]*catalog.Function, error) {
	sub, err := fs.Sub(fsys, dir)
	if err != nil {
		return nil, fmt.Errorf("seed: %s: %w", dir, err)
	}
	var out []*catalog.Function
	err = stream(sub, FunctionsFile, func(fn Function) error {
		args := make([]*catalog.Argument, 0, len(fn.Args))
		for _, a := range fn.Args {
			args = append(args, &catalog.Argument{
				Name:       a.Name,
				Type:       &ast.TypeName{Name: a.Type},
				HasDefault: a.HasDefault,
				Mode:       argMode(a.Mode),
			})
		}
		out = append(out, &catalog.Function{
			Name:               fn.Name,
			Args:               args,
			ReturnType:         &ast.TypeName{Name: fn.Returns},
			ReturnTypeNullable: fn.Nullable,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Relations streams the relations a dialect ships with in the named schema
// into the form the engine catalogs use.
func Relations(fsys fs.FS, dir, schema string) ([]*catalog.Table, error) {
	sub, err := fs.Sub(fsys, dir)
	if err != nil {
		return nil, fmt.Errorf("seed: %s: %w", dir, err)
	}
	var out []*catalog.Table
	err = stream(sub, RelationsFile, func(rel Relation) error {
		if rel.Schema != schema {
			return nil
		}
		table := &catalog.Table{
			Rel:     &ast.TableName{Catalog: rel.Catalog, Schema: rel.Schema, Name: rel.Name},
			Columns: make([]*catalog.Column, 0, len(rel.Columns)),
		}
		for _, col := range rel.Columns {
			column := &catalog.Column{
				Name:      col.Name,
				Type:      ast.TypeName{Name: col.Type},
				IsNotNull: col.NotNull,
				IsArray:   col.Array,
			}
			if col.Length > 0 {
				length := col.Length
				column.Length = &length
			}
			table.Columns = append(table.Columns, column)
		}
		out = append(out, table)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// argModes maps the mode a record carries to the one the catalog uses. The
// letters are PostgreSQL's, which is where the generated function lists come
// from.
var argModes = map[string]ast.FuncParamMode{
	"i": ast.FuncParamIn,
	"o": ast.FuncParamOut,
	"b": ast.FuncParamInOut,
	"v": ast.FuncParamVariadic,
	"t": ast.FuncParamTable,
	"d": ast.FuncParamDefault,
}

func argMode(mode string) ast.FuncParamMode {
	if m, ok := argModes[mode]; ok {
		return m
	}
	return ast.FuncParamIn
}

// ArgMode is the letter a record carries for a catalog argument mode, the
// inverse of the mapping applied when reading one.
func ArgMode(mode ast.FuncParamMode) string {
	for letter, m := range argModes {
		if m == mode {
			return letter
		}
	}
	return "i"
}

type builder struct {
	cat        *core.Catalog
	settings   Settings
	dialectOID int64

	// oids maps a lowercased type name to its OID, and categories records the
	// category each was seeded under, in the order they were read.
	oids       map[string]int64
	categories []categorized

	// namespaces maps a schema name to its OID, so that a run of relations in
	// the same schema resolves it once.
	namespaces map[string]int64

	// seenCasts records the pairs already registered. The alias rules and the
	// category rules overlap, and a cast pair is unique in the catalog.
	seenCasts map[[2]int64]bool
}

type categorized struct {
	name     string
	category string
}

func (b *builder) addType(t Type) error {
	for _, name := range append([]string{t.Name}, t.Aliases...) {
		if _, err := b.createType(name, t.Category); err != nil {
			return fmt.Errorf("type %q: %w", name, err)
		}
	}
	return b.aliasCasts(t)
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
			if err := b.addCast(Cast{Source: src, Target: tgt, Context: "i"}); err != nil {
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

// constKinds are the kinds of literal a dialect may name a type for.
var constKinds = []string{core.ConstInteger, core.ConstFloat, core.ConstString, core.ConstBool}

func (b *builder) consts() error {
	for kind := range b.settings.Const {
		if !slices.Contains(constKinds, kind) {
			return fmt.Errorf("seed %s: unknown constant kind %q, want one of %s",
				b.settings.Dialect, kind, strings.Join(constKinds, ", "))
		}
	}
	// Written in a fixed order rather than the map's: a catalog built from the
	// same dialect twice has to come out byte for byte the same, so that the
	// cache stores one copy of it rather than one per run.
	for _, kind := range constKinds {
		name, ok := b.settings.Const[kind]
		if !ok {
			continue
		}
		if _, ok := b.oids[strings.ToLower(name)]; !ok {
			return fmt.Errorf("seed %s: constant %s names unknown type %q", b.settings.Dialect, kind, name)
		}
		if err := b.cat.SetConstType(b.dialectOID, kind, name); err != nil {
			return err
		}
	}
	for key, name := range map[string]string{
		core.FlagBoolType:    b.settings.Bool,
		core.FlagLimitType:   b.settings.Limit,
		core.FlagUntypedType: b.settings.Untyped,
	} {
		if name == "" {
			continue
		}
		if _, ok := b.oids[strings.ToLower(name)]; !ok {
			return fmt.Errorf("seed %s: %s names unknown type %q", b.settings.Dialect, key, name)
		}
		if err := b.cat.SetDialectFlag(b.dialectOID, key, name); err != nil {
			return err
		}
	}
	if b.settings.PropagateNullable {
		if err := b.cat.SetDialectFlag(b.dialectOID, core.FlagPropagateNullable, "true"); err != nil {
			return err
		}
	}
	if b.settings.QualifyDuplicateColumns {
		if err := b.cat.SetDialectFlag(b.dialectOID, core.FlagQualifyDuplicateColumns, "true"); err != nil {
			return err
		}
	}
	// A schema declares types the seed knows nothing about — enums, domains,
	// arrays, a SQLite column typed whatever the author felt like. Recording
	// the comparison operators lets the catalog give those types the same ones.
	if len(b.settings.Comparison) > 0 {
		if err := b.cat.SetDialectFlag(b.dialectOID, core.FlagComparisonOperators,
			strings.Join(b.settings.Comparison, ",")); err != nil {
			return err
		}
	}
	// The cast categories are recorded for the same reason: a type an
	// extension adds joins its category's casts.
	if b.settings.CastCategories != "" {
		return b.cat.SetDialectFlag(b.dialectOID, core.FlagCastCategories, b.settings.CastCategories)
	}
	return nil
}

func (b *builder) categoryOperators() error {
	s := b.settings
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

func (b *builder) addOperator(op Operator) error {
	return b.operator(op.Name, op.Left, op.Right, op.Result, 0)
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

// categoryCasts applies the rule that makes the types in a category mutually
// castable, once every type has been read.
func (b *builder) categoryCasts() error {
	categories := b.settings.CastCategories
	if categories == "" {
		return nil
	}
	anyCategory := categories == "*"
	for _, src := range b.categories {
		if !anyCategory && !strings.Contains(categories, src.category) {
			continue
		}
		for _, tgt := range b.categories {
			if src.name == tgt.name {
				continue
			}
			if !anyCategory && src.category != tgt.category {
				continue
			}
			if err := b.addCast(Cast{Source: src.name, Target: tgt.name, Context: "i"}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *builder) addCast(c Cast) error {
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

func (b *builder) addFunction(fn Function) error {
	returnOID, err := b.funcType(fn.Returns)
	if err != nil {
		return fmt.Errorf("function %q: %w", fn.Name, err)
	}
	args := make([]core.ProcArg, 0, len(fn.Args))
	for _, a := range fn.Args {
		// An out or table parameter is part of the result, not the call.
		switch a.Mode {
		case "o", "t":
			continue
		}
		argOID, err := b.funcType(a.Type)
		if err != nil {
			return fmt.Errorf("function %q: %w", fn.Name, err)
		}
		args = append(args, core.ProcArg{
			Name:       a.Name,
			TypeOID:    argOID,
			Mode:       a.Mode,
			HasDefault: a.HasDefault,
		})
	}
	_, err = b.cat.CreateProc(core.ProcSpec{
		Name:           fn.Name,
		DialectOID:     b.dialectOID,
		Kind:           fn.Kind,
		ReturnTypeOID:  returnOID,
		ReturnNullable: fn.Nullable,
		NeverNull:      fn.NeverNull,
		Args:           args,
	})
	if err != nil {
		return fmt.Errorf("function %q: %w", fn.Name, err)
	}
	return nil
}

func (b *builder) addRelation(rel Relation) error {
	if rel.Name == "" {
		return errors.New("relation has no name")
	}
	nsOID, err := b.namespace(rel.Schema)
	if err != nil {
		return err
	}
	kind := rel.Kind
	if kind == "" {
		kind = "r"
	}
	classOID, err := b.cat.CreateClass(nsOID, rel.Name, kind)
	if err != nil {
		return err
	}
	for i, col := range rel.Columns {
		name := col.Type
		if col.Array {
			name += core.ArraySuffix
		}
		typeOID, err := b.columnType(name)
		if err != nil {
			return fmt.Errorf("relation %q column %q: %w", rel.Name, col.Name, err)
		}
		if err := b.cat.CreateAttributeSpec(core.AttributeSpec{
			ClassOID:   classOID,
			Name:       col.Name,
			TypeOID:    typeOID,
			Num:        i + 1,
			NotNull:    col.NotNull,
			DeclType:   col.Type,
			TypeLength: col.Length,
		}); err != nil {
			return fmt.Errorf("relation %q column %q: %w", rel.Name, col.Name, err)
		}
	}
	return nil
}

// namespace resolves the schema a relation belongs to, creating it the first
// time it is named. A relation with no schema belongs to the default one.
func (b *builder) namespace(schema string) (int64, error) {
	if schema == "" {
		schema = "public"
	}
	if oid, ok := b.namespaces[schema]; ok {
		return oid, nil
	}
	oid, err := b.cat.NamespaceOID(schema)
	if err != nil {
		if oid, err = b.cat.CreateNamespace(schema); err != nil {
			return 0, err
		}
	}
	b.namespaces[schema] = oid
	return oid, nil
}

// columnType resolves a column's type, which unlike a function signature may
// name an array.
func (b *builder) columnType(name string) (int64, error) {
	key := strings.ToLower(name)
	if oid, ok := b.oids[key]; ok {
		return oid, nil
	}
	oid, err := b.cat.ResolveTypeName(key)
	if err != nil {
		return 0, err
	}
	b.oids[key] = oid
	return oid, nil
}

// funcType resolves the type a function signature names. Signatures reference
// pseudo types ("any", "record") and types no dialect bothers to list, so an
// unknown name is registered rather than rejected.
func (b *builder) funcType(name string) (int64, error) {
	if name == "" {
		return 0, nil
	}
	return b.createType(name, "U")
}
