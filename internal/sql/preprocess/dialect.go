package preprocess

import (
	"github.com/sqlc-dev/sqlc/internal/config"
)

// Style is the native placeholder syntax a dialect uses for bind parameters.
type Style int

const (
	// StyleDollar numbers parameters as $1, $2, ... (PostgreSQL)
	StyleDollar Style = iota
	// StyleQuestion uses an unnumbered ? for every parameter (MySQL, ClickHouse)
	StyleQuestion
	// StyleOrdinal numbers parameters as ?1, ?2, ... (SQLite)
	StyleOrdinal
	// StyleAtName keeps parameters named as @name (GoogleSQL)
	StyleAtName
)

// Dialect describes the lexical rules the preprocessor needs to walk a query
// without parsing it: how the dialect quotes strings and identifiers, how it
// writes comments and which placeholder syntax it accepts.
//
// The preprocessor never interprets SQL. It only needs to know enough to skip
// the regions of a query where sqlc syntax must not be rewritten.
type Dialect struct {
	// Style is the native placeholder syntax used for rewritten parameters.
	Style Style

	// AtSign reports whether a bare "@name" is sqlc named-parameter syntax.
	// It is false for MySQL, where "@name" is a user variable.
	AtSign bool

	// DollarQuote reports whether $tag$ ... $tag$ string literals are valid.
	DollarQuote bool

	// DollarNumber reports whether $1 is a bind parameter.
	DollarNumber bool

	// Question reports whether ? is a bind parameter.
	Question bool

	// Backtick reports whether `ident` is a quoted identifier.
	Backtick bool

	// HashComment reports whether # starts a line comment.
	HashComment bool

	// DoubleQuoteString reports whether "..." is a string literal rather than
	// a quoted identifier. MySQL treats it as a string unless ANSI_QUOTES is
	// enabled.
	DoubleQuoteString bool

	// NestedBlockComment reports whether /* ... /* ... */ ... */ nests.
	NestedBlockComment bool

	// Backslash reports whether a backslash escapes the next character inside
	// a string literal.
	Backslash bool

	// FoldIdentifier reports whether the dialect lowercases unquoted
	// identifiers. It decides the case of a parameter named by a bare
	// reference, e.g. sqlc.arg(FooBar).
	FoldIdentifier bool
}

var dialects = map[config.Engine]Dialect{
	config.EnginePostgreSQL: {
		Style:              StyleDollar,
		AtSign:             true,
		DollarQuote:        true,
		DollarNumber:       true,
		NestedBlockComment: true,
		FoldIdentifier:     true,
	},
	config.EngineMySQL: {
		Style:             StyleQuestion,
		Question:          true,
		Backtick:          true,
		HashComment:       true,
		DoubleQuoteString: true,
		Backslash:         true,
		FoldIdentifier:    true,
	},
	config.EngineSQLite: {
		Style:     StyleOrdinal,
		AtSign:    true,
		Question:  true,
		Backtick:  true,
		Backslash: true,
	},
	config.EngineGoogleSQL: {
		Style:     StyleAtName,
		AtSign:    true,
		Question:  true,
		Backtick:  true,
		Backslash: true,
	},
	config.EngineClickHouse: {
		Style:       StyleQuestion,
		Question:    true,
		Backtick:    true,
		HashComment: true,
		Backslash:   true,
	},
}

// DialectFor returns the lexical rules for an engine.
func DialectFor(engine config.Engine) (Dialect, bool) {
	d, ok := dialects[engine]
	return d, ok
}

// hasNamedSupport reports whether repeated uses of the same parameter name can
// share a single placeholder. MySQL and ClickHouse send an argument per "?", so
// every occurrence needs its own number.
func (d Dialect) hasNamedSupport() bool {
	return d.Style != StyleQuestion
}
