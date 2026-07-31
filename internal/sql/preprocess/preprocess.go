// Package preprocess rewrites sqlc's query syntax into native SQL before any
// engine parser sees a query.
//
// sqlc.arg(), sqlc.narg(), sqlc.slice(), sqlc.embed() and @name are sqlc
// syntax, not SQL. Historically they survived all the way through each engine's
// parser and were replaced by walking the resulting AST, which meant every
// engine had to round-trip a schema-qualified function call it did not
// understand. This package removes that requirement: it scans the query text
// with a dialect-parameterized lexer, replaces each construct with the engine's
// native placeholder, and records what it replaced. Engines then only ever see
// valid SQL.
package preprocess

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/named"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

// Embed is a rewritten sqlc.embed(table) call.
type Embed struct {
	Table *ast.TableName

	// Location is the offset of the rewritten "table.*" in the preprocessed
	// text. The compiler uses it to tell a rewritten star reference apart from
	// one the user wrote.
	Location int

	param string
}

// Orig returns the sqlc syntax this embed was rewritten from.
func (e Embed) Orig() string {
	return fmt.Sprintf("sqlc.embed(%s)", e.param)
}

// EmbedSet is the set of embeds in a single statement.
type EmbedSet []*Embed

// Find returns the embed recorded at any of the given locations. Engines do not
// all populate the same location fields, so callers pass every location that
// could identify the node.
func (es EmbedSet) Find(locations ...int) (*Embed, bool) {
	for _, e := range es {
		for _, loc := range locations {
			if loc == e.Location {
				return e, true
			}
		}
	}
	return nil, false
}

// Statement holds everything the preprocessor learned about a single statement.
type Statement struct {
	// Start and End delimit the statement in the rewritten text.
	Start, End int

	// Params maps each placeholder number to its sqlc metadata.
	Params *named.ParamSet

	// Embeds records the sqlc.embed() calls that were rewritten to "table.*".
	Embeds EmbedSet

	// Slices holds the offset, in the rewritten text, of every placeholder that
	// came from sqlc.slice(), mapped to the parameter name.
	Slices map[int]string

	// Numbers maps the offset of every placeholder in the rewritten text to the
	// number assigned to it. Engines number bind parameters in whatever order
	// they happen to convert the AST, so the compiler uses this to restore
	// source order.
	Numbers map[int]int

	// Dollar reports whether the statement's own placeholders are numbered.
	Dollar bool

	// ParamErr is set when the placeholders the user wrote are inconsistent:
	// numbered and unnumbered styles mixed, or a gap in the numbering.
	ParamErr error

	// Err is a validation error found while rewriting. The statement text is
	// left untouched when this is set, so the engine still sees the original
	// query and the error can be reported against it.
	Err error
}

// Result is the outcome of preprocessing a SQL file.
type Result struct {
	// Text is the rewritten SQL. It has the same number of statements and the
	// same line structure as the input.
	Text string

	stmts []*Statement
	edits []edit
}

// edit records one rewritten span, in both coordinate systems.
type edit struct {
	oldStart, oldEnd int
	newStart, newEnd int
}

// Statement returns the preprocessing results for the statement covering the
// given offset in the rewritten text. It never returns nil.
func (r *Result) Statement(location int) *Statement {
	for _, s := range r.stmts {
		if location < s.End {
			return s
		}
	}
	return &Statement{Params: named.NewParamSet(nil, true)}
}

// Origin maps an offset in the rewritten text back to the offset in the
// original source, so errors can be reported against what the user wrote.
func (r *Result) Origin(location int) int {
	delta := 0
	for _, e := range r.edits {
		if location < e.newStart {
			break
		}
		if location < e.newEnd {
			return e.oldStart
		}
		delta += (e.oldEnd - e.oldStart) - (e.newEnd - e.newStart)
	}
	return location + delta
}

type kind int

const (
	kindArg kind = iota
	kindNarg
	kindSlice
	kindEmbed
	kindNative
)

// occurrence is one construct found in the query text.
type occurrence struct {
	kind  kind
	start int // offset in the original source
	end   int
	name  string
	// ident is the argument exactly as written, used to rebuild "table.*" for
	// embeds so quoting is preserved.
	ident string
	// number is the placeholder number assigned to this occurrence.
	number int
	// explicit reports whether a native placeholder carried its own number.
	explicit bool
	err      *sqlerr.Error
}

// File rewrites every sqlc construct in src to native SQL for the given engine.
func File(engine config.Engine, src string) (*Result, error) {
	d, ok := DialectFor(engine)
	if !ok {
		return nil, fmt.Errorf("preprocess: unknown engine: %s", engine)
	}
	return Dialected(d, src)
}

// Dialected rewrites src using an explicit dialect.
func Dialected(d Dialect, src string) (*Result, error) {
	l := &lexer{d: d, src: src}
	res := &Result{}

	var b strings.Builder
	b.Grow(len(src))

	for _, span := range l.statements() {
		start, end := span[0], span[1]
		occs := l.scan(start, end)
		stmt := &Statement{Start: b.Len()}

		// A statement that fails validation is copied through unchanged so the
		// engine parses what the user actually wrote and the error keeps its
		// original shape.
		if err := firstError(occs); err != nil {
			b.WriteString(src[start:end])
			stmt.End = b.Len()
			err.Location = stmt.Start + (err.Location - start)
			stmt.Err = err
			stmt.Params = named.NewParamSet(nil, d.hasNamedSupport())
			stmt.Dollar = true
			res.stmts = append(res.stmts, stmt)
			continue
		}

		stmt.Dollar, stmt.ParamErr = validatePlaceholders(occs)
		stmt.Params = number(d, occs)
		stmt.Numbers = map[int]int{}

		prev := start
		for i := range occs {
			occ := &occs[i]
			b.WriteString(src[prev:occ.start])
			newStart := b.Len()
			replacement := d.replacement(occ)
			b.WriteString(replacement)
			newEnd := b.Len()
			prev = occ.end

			if replacement != src[occ.start:occ.end] {
				res.edits = append(res.edits, edit{
					oldStart: occ.start, oldEnd: occ.end,
					newStart: newStart, newEnd: newEnd,
				})
			}

			if occ.kind != kindEmbed {
				stmt.Numbers[newStart+sliceMarkerLen(d, occ)] = occ.number
			}

			switch occ.kind {
			case kindEmbed:
				stmt.Embeds = append(stmt.Embeds, &Embed{
					Table:    &ast.TableName{Name: occ.name},
					Location: newStart,
					param:    occ.ident,
				})
			case kindSlice:
				if stmt.Slices == nil {
					stmt.Slices = map[int]string{}
				}
				// The placeholder follows the /*SLICE:name*/ marker, so record
				// the offset the engine will report for the parameter node.
				stmt.Slices[newStart+sliceMarkerLen(d, occ)] = occ.name
			}
		}
		b.WriteString(src[prev:end])
		stmt.End = b.Len()
		res.stmts = append(res.stmts, stmt)
	}

	res.Text = b.String()
	return res, nil
}

func firstError(occs []occurrence) *sqlerr.Error {
	for _, occ := range occs {
		if occ.err != nil {
			return occ.err
		}
	}
	return nil
}

// sliceMarkerLen returns the length of the /*SLICE:name*/ prefix a slice
// placeholder is written with, or 0 when there is none.
func sliceMarkerLen(d Dialect, occ *occurrence) int {
	if occ.kind != kindSlice || d.Style == StyleDollar {
		return 0
	}
	return len("/*SLICE:") + len(occ.name) + len("*/")
}

// validatePlaceholders checks the placeholders the user wrote, reproducing the
// rules the AST-based validator applied before sqlc syntax was rewritten: a
// query may not mix numbered and unnumbered placeholders, and the numbers it
// does use may not have gaps.
func validatePlaceholders(occs []occurrence) (dollar bool, err error) {
	var numbered, unnumbered bool
	seen := map[int]bool{}
	for _, occ := range occs {
		if occ.kind != kindNative {
			continue
		}
		if occ.explicit {
			numbered = true
			seen[occ.number] = true
		} else {
			unnumbered = true
		}
	}
	if numbered && unnumbered {
		return false, errors.New("can not mix $1 format with ? format")
	}
	for i := 1; i <= len(seen); i++ {
		if !seen[i] {
			return !unnumbered, &sqlerr.Error{
				Code:    "42P18",
				Message: fmt.Sprintf("could not determine data type of parameter $%d", i),
			}
		}
	}
	return !unnumbered, nil
}

// number assigns a placeholder number to every occurrence and builds the
// parameter set for the statement.
func number(d Dialect, occs []occurrence) *named.ParamSet {
	if d.Style == StyleQuestion {
		// Every "?" is sent as its own argument, so existing and rewritten
		// placeholders are numbered together in source order.
		ps := named.NewParamSet(nil, false)
		for i := range occs {
			occ := &occs[i]
			if occ.kind == kindEmbed {
				continue
			}
			occ.number = ps.Add(param(occ))
		}
		return ps
	}

	// Numbered dialects keep the numbers the user wrote and fill in the gaps.
	numbs := map[int]bool{}
	count := 0
	for i := range occs {
		occ := &occs[i]
		if occ.kind != kindNative {
			continue
		}
		count++
		if !occ.explicit {
			occ.number = count
		}
		numbs[occ.number] = true
	}
	ps := named.NewParamSet(numbs, d.hasNamedSupport())
	for i := range occs {
		occ := &occs[i]
		if occ.kind == kindEmbed || occ.kind == kindNative {
			continue
		}
		occ.number = ps.Add(param(occ))
	}
	return ps
}

func param(occ *occurrence) named.Param {
	switch occ.kind {
	case kindNarg:
		return named.NewUserNullableParam(occ.name)
	case kindSlice:
		return named.NewSqlcSlice(occ.name)
	default:
		return named.NewParam(occ.name)
	}
}

// replacement returns the native SQL that stands in for an occurrence.
func (d Dialect) replacement(occ *occurrence) string {
	switch occ.kind {
	case kindEmbed:
		return occ.ident + ".*"
	case kindNative:
		return occ.ident
	}
	switch d.Style {
	case StyleAtName:
		// GoogleSQL names its parameters, so "@name" is already native and a
		// user-written one is rewritten to itself.
		if occ.kind == kindSlice {
			return fmt.Sprintf("/*SLICE:%s*/@%s", occ.name, occ.name)
		}
		return "@" + occ.name
	case StyleDollar:
		// PostgreSQL passes a slice as a single array parameter, so it needs no
		// expansion marker.
		return fmt.Sprintf("$%d", occ.number)
	case StyleOrdinal:
		if occ.kind == kindSlice {
			// This sequence is also replicated in
			// internal/codegen/golang.Field since it's needed during template
			// generation for replacement.
			return fmt.Sprintf("/*SLICE:%s*/?", occ.name)
		}
		return fmt.Sprintf("?%d", occ.number)
	default:
		if occ.kind == kindSlice {
			return fmt.Sprintf("/*SLICE:%s*/?", occ.name)
		}
		return "?"
	}
}
