package migrations

import (
	"bufio"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type PreprocessMode int

const (
	// PreprocessModeParse keeps schema text usable for sqlc parsing/codegen and
	// warns when semantic psql commands are stripped.
	PreprocessModeParse PreprocessMode = iota
	// PreprocessModeApply rejects semantic psql commands before sqlc applies
	// schema text to a live database.
	PreprocessModeApply
)

// RemoveRollbackStatements returns the up-migration portion of a migration file
// by discarding everything from the first recognized rollback marker onward.
//
// The supported markers match the migration formats sqlc already understands
// during schema loading. The function preserves original line ordering up to the
// rollback boundary so downstream parsing and error locations remain stable.
//
// goose:       -- +goose Down
// sql-migrate: -- +migrate Down
// tern:        ---- create above / drop below ----
// dbmate:      -- migrate:down
func RemoveRollbackStatements(contents string) string {
	s := bufio.NewScanner(strings.NewReader(contents))
	var lines []string
	for s.Scan() {
		statement := strings.ToLower(s.Text())
		if strings.HasPrefix(statement, "-- +goose down") {
			break
		}
		if strings.HasPrefix(statement, "-- +migrate down") {
			break
		}
		if strings.HasPrefix(statement, "---- create above / drop below ----") {
			break
		}
		if strings.HasPrefix(statement, "-- migrate:down") {
			break
		}
		lines = append(lines, s.Text())
	}
	return strings.Join(lines, "\n")
}

// PreprocessSchema normalizes schema text for sqlc before parsing or applying
// it to managed databases. Rollback sections are always removed; PostgreSQL
// schemas preserve server-side SQL, including PL/pgSQL bodies and
// extension/language DDL, while additionally stripping top-level psql
// meta-commands that are not valid SQL. The returned warnings describe
// semantic psql commands that were stripped and any best-effort approximation
// of psql session semantics used during parsing-oriented preprocessing.
func PreprocessSchema(contents, engine string) (string, []string, error) {
	return preprocessSchema(contents, engine, PreprocessModeParse)
}

// PreprocessSchemaForApply normalizes schema text for commands that will apply
// the resulting DDL to a live database. Semantic psql commands are rejected in
// this mode because sqlc cannot reproduce their effects during execution.
func PreprocessSchemaForApply(contents, engine string) (string, []string, error) {
	return preprocessSchema(contents, engine, PreprocessModeApply)
}

// preprocessSchema applies the shared engine-aware preprocessing pipeline using
// the caller's requested strictness for semantic psql commands.
func preprocessSchema(contents, engine string, mode PreprocessMode) (string, []string, error) {
	contents = RemoveRollbackStatements(contents)
	if engine == "postgresql" {
		var (
			err      error
			warnings []string
		)
		contents, warnings, err = removePsqlMetaCommands(contents, mode)
		if err != nil {
			return "", nil, err
		}
		return contents, warnings, nil
	}
	return contents, nil, nil
}

// RemovePsqlMetaCommands strips top-level psql meta-command lines from SQL
// input while preserving valid SQL text, literal contents, comments, and line
// structure. LF and CRLF line endings are preserved; lone CR line endings are
// normalized to LF so downstream line-number accounting stays correct.
//
// pg_dump can emit client-only commands such as `\restrict KEY`,
// `\unrestrict KEY`, and `\connect foo`. These are meaningful to the psql
// client but invalid for sqlc's SQL parsers. The implementation uses a
// single-pass state machine so it only removes backslash directives that
// appear at true statement line starts, not backslashes embedded in string
// literals, double-quoted identifiers, dollar-quoted bodies, or comments.
//
// The line-start matcher is intentionally broader than the currently
// documented psql command set: sqlc strips unknown future backslash directives
// and dump-tool variants the same way it strips known psql meta-commands,
// because none of them are valid SQL input to the parser.
func RemovePsqlMetaCommands(contents string) (string, []string, error) {
	return removePsqlMetaCommands(contents, PreprocessModeParse)
}

func removePsqlMetaCommands(contents string, mode PreprocessMode) (string, []string, error) {
	if contents == "" {
		return contents, nil, nil
	}

	var out strings.Builder
	out.Grow(len(contents))
	warningsByToken := map[string]struct{}{}

	lineStart := true
	inSingle := false
	inDouble := false
	inDollar := false
	singleAllowsBackslash := false
	standardConformingStringsOff := false
	warnedApproximateSessionSemantics := false
	var dollarTag string
	blockDepth := 0
	n := len(contents)
	var statement strings.Builder

	for i := 0; ; {
		// Only top-level line starts are eligible for psql meta-command removal.
		// Leading horizontal whitespace is preserved so stripped lines keep their
		// original indentation and line count for downstream position accounting.
		if lineStart && !inSingle && !inDouble && blockDepth == 0 && !inDollar {
			start := i
			for i < n {
				c := contents[i]
				if c == ' ' || c == '\t' {
					i++
					continue
				}
				break
			}
			if i < n && contents[i] == '\\' && i+1 < n && isMetaCommandStart(contents[i+1]) {
				lineEnd := i
				for lineEnd < n && contents[lineEnd] != '\r' && contents[lineEnd] != '\n' {
					lineEnd++
				}
				sep := findMetaCommandSeparator(contents, i, lineEnd)
				if token := metaCommandToken(contents[i:lineEnd]); token != "" {
					if sep < 0 && hasInvalidSeparatorCandidate(contents, i, lineEnd) {
						if start < i {
							out.WriteString(contents[start:i])
						}
						goto normalLexing
					}
					switch {
					case isUnsupportedConditionalMetaCommand(token):
						return "", nil, fmt.Errorf("psql conditional directives (%s) are not supported in schema preprocessing", token)
					case isSemanticMetaCommand(token):
						if mode == PreprocessModeApply {
							return "", nil, fmt.Errorf("psql meta-command %s is not supported when applying schema preprocessing to a live database", token)
						}
						warningsByToken[token] = struct{}{}
						if token == `\copy` && readsPsqlCopyData(contents, i, lineEnd) {
							var ok bool
							i, ok = stripPsqlCopyData(&out, contents, lineEnd)
							if !ok {
								return "", nil, fmt.Errorf(`psql meta-command \copy ... from stdin requires a terminating \. line during schema preprocessing`)
							}
							lineStart = true
							continue
						}
					}
				}
				if sep >= 0 {
					// Resume normal SQL lexing after a valid `\\` separator so
					// multiline literals/comments in the preserved tail keep parser
					// state for following lines.
					i = sep + 2
					lineStart = false
					continue
				}
				// Keep LF and CRLF intact when removing a command line. Bare CR is
				// normalized by writeLineEnding() so stripped output still matches
				// downstream line-number accounting.
				i = writeLineEnding(&out, contents, lineEnd)
				lineStart = true
				continue
			}
			if start < i {
				out.WriteString(contents[start:i])
			}
			if i >= n {
				break
			}
		}
		if i >= n {
			break
		}

	normalLexing:
		c := contents[i]
		if inSingle {
			// Inside string literals, only quote termination rules matter. Escape
			// strings (`E'...'`) and standard_conforming_strings=off additionally
			// allow backslash-escaped bytes.
			if singleAllowsBackslash && c == '\\' && i+1 < n {
				out.WriteByte(c)
				out.WriteByte(contents[i+1])
				consumeStatementFragment(contents[i:i+2], &statement, &standardConformingStringsOff, &warnedApproximateSessionSemantics, false)
				lineStart = false
				i += 2
				continue
			}
			if isLineBreak(c) {
				i = writeLineEnding(&out, contents, i)
				lineStart = true
				consumeStatementFragment(" ", &statement, &standardConformingStringsOff, &warnedApproximateSessionSemantics, false)
				continue
			}
			out.WriteByte(c)
			consumeStatementFragment(contents[i:i+1], &statement, &standardConformingStringsOff, &warnedApproximateSessionSemantics, false)
			if c == '\'' {
				if i+1 < n && contents[i+1] == '\'' {
					out.WriteByte(contents[i+1])
					consumeStatementFragment(contents[i+1:i+2], &statement, &standardConformingStringsOff, &warnedApproximateSessionSemantics, false)
					i += 2
					lineStart = false
					continue
				}
				inSingle = false
				singleAllowsBackslash = false
			}
			lineStart = false
			i++
			continue
		}

		if inDouble {
			// Double-quoted identifiers may span lines, so line-start backslashes
			// inside them must never be treated as top-level meta-commands.
			if isLineBreak(c) {
				i = writeLineEnding(&out, contents, i)
				lineStart = true
				consumeStatementFragment(" ", &statement, &standardConformingStringsOff, &warnedApproximateSessionSemantics, false)
				continue
			}
			out.WriteByte(c)
			consumeStatementFragment(contents[i:i+1], &statement, &standardConformingStringsOff, &warnedApproximateSessionSemantics, false)
			if c == '"' {
				if i+1 < n && contents[i+1] == '"' {
					out.WriteByte(contents[i+1])
					consumeStatementFragment(contents[i+1:i+2], &statement, &standardConformingStringsOff, &warnedApproximateSessionSemantics, false)
					i += 2
					lineStart = false
					continue
				}
				inDouble = false
			}
			lineStart = false
			i++
			continue
		}

		if inDollar {
			// Dollar-quoted bodies are opaque until their exact tag reappears.
			if strings.HasPrefix(contents[i:], dollarTag) {
				out.WriteString(dollarTag)
				i += len(dollarTag)
				inDollar = false
				lineStart = false
				continue
			}
			if isLineBreak(c) {
				i = writeLineEnding(&out, contents, i)
				lineStart = true
				continue
			}
			out.WriteByte(c)
			lineStart = false
			i++
			continue
		}

		if blockDepth > 0 {
			// Block comments may nest in PostgreSQL, so maintain explicit depth.
			if c == '/' && i+1 < n && contents[i+1] == '*' {
				blockDepth++
				out.WriteString("/*")
				consumeStatementFragment(" ", &statement, &standardConformingStringsOff, &warnedApproximateSessionSemantics, false)
				i += 2
				lineStart = false
				continue
			}
			if c == '*' && i+1 < n && contents[i+1] == '/' {
				blockDepth--
				out.WriteString("*/")
				consumeStatementFragment(" ", &statement, &standardConformingStringsOff, &warnedApproximateSessionSemantics, false)
				i += 2
				lineStart = false
				continue
			}
			if isLineBreak(c) {
				i = writeLineEnding(&out, contents, i)
				lineStart = true
				consumeStatementFragment(" ", &statement, &standardConformingStringsOff, &warnedApproximateSessionSemantics, false)
				continue
			}
			out.WriteByte(c)
			lineStart = false
			i++
			continue
		}

		switch c {
		case '\'':
			inSingle = true
			singleAllowsBackslash = standardConformingStringsOff || isEscapeStringPrefix(contents, i)
			out.WriteByte(c)
			consumeStatementFragment(contents[i:i+1], &statement, &standardConformingStringsOff, &warnedApproximateSessionSemantics, false)
			lineStart = false
			i++
			continue
		case '"':
			inDouble = true
			out.WriteByte(c)
			consumeStatementFragment(contents[i:i+1], &statement, &standardConformingStringsOff, &warnedApproximateSessionSemantics, false)
			lineStart = false
			i++
			continue
		case '$':
			if tag := matchDollarTagStart(contents, i); tag != "" {
				dollarTag = tag
				inDollar = true
				out.WriteString(dollarTag)
				i += len(dollarTag)
				lineStart = false
				continue
			}
		case '/':
			if i+1 < n && contents[i+1] == '*' {
				blockDepth = 1
				out.WriteString("/*")
				consumeStatementFragment(" ", &statement, &standardConformingStringsOff, &warnedApproximateSessionSemantics, false)
				i += 2
				lineStart = false
				continue
			}
		case '-':
			if i+1 < n && contents[i+1] == '-' {
				// Line comments are copied through verbatim and treated as inert
				// text so quote-like markers inside comments cannot perturb state.
				for i < n {
					if contents[i] == '\r' || contents[i] == '\n' {
						lineStart = true
						i = writeLineEnding(&out, contents, i)
						consumeStatementFragment(" ", &statement, &standardConformingStringsOff, &warnedApproximateSessionSemantics, false)
						goto nextChar
					}
					out.WriteByte(contents[i])
					i++
				}
				consumeStatementFragment(" ", &statement, &standardConformingStringsOff, &warnedApproximateSessionSemantics, false)
				goto nextChar
			}
		}

		if isLineBreak(c) {
			i = writeLineEnding(&out, contents, i)
			lineStart = true
			consumeStatementFragment(" ", &statement, &standardConformingStringsOff, &warnedApproximateSessionSemantics, false)
			continue
		}
		out.WriteByte(c)
		consumeStatementFragment(contents[i:i+1], &statement, &standardConformingStringsOff, &warnedApproximateSessionSemantics, true)
		lineStart = false
		i++
	nextChar:
	}

	return out.String(), schemaPreprocessWarnings(warningsByToken, warnedApproximateSessionSemantics), nil
}

// hasInvalidSeparatorCandidate reports whether a directive line contains an
// unquoted `\\` pair that is not shaped like the validated separator token.
// Those lines are preserved verbatim so sqlc does not normalize invalid psql
// syntax into different executable SQL.
func hasInvalidSeparatorCandidate(contents string, start, lineEnd int) bool {
	inSingle := false
	inDouble := false
	for i := start; i+1 < lineEnd; i++ {
		switch contents[i] {
		case '\'':
			if inSingle {
				if i+1 < lineEnd && contents[i+1] == '\'' {
					i++
					continue
				}
				inSingle = false
				continue
			}
			if !inDouble {
				inSingle = true
			}
		case '"':
			if inDouble {
				if i+1 < lineEnd && contents[i+1] == '"' {
					i++
					continue
				}
				inDouble = false
				continue
			}
			if !inSingle {
				inDouble = true
			}
		}
		if inSingle || inDouble {
			continue
		}
		if contents[i] != '\\' || contents[i+1] != '\\' {
			continue
		}
		if i == start || !isHorizontalSpace(contents[i-1]) {
			return true
		}
	}
	return false
}

// isMetaCommandStart reports whether b can begin a top-level backslash
// directive that sqlc should strip before SQL parsing. The accepted set is
// intentionally broader than today's documented psql commands so unknown or
// future client-side directives are treated the same as known meta-commands.
func isMetaCommandStart(b byte) bool {
	return isDollarTagChar(b) || b == '!' || b == '?' || b == ';' || b == '\\'
}

// isDollarTagChar reports whether b is valid inside a PostgreSQL dollar-quote
// tag. The opening and closing delimiters must match exactly, so the stripper
// uses the same character class when detecting tagged literals.
func isDollarTagChar(b byte) bool {
	return b == '_' || (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
}

// isEscapeStringPrefix reports whether the quote at quotePos begins a
// PostgreSQL escape string literal (`E'...'`) rather than a plain string or a
// longer identifier token that merely ends with `e`.
func isEscapeStringPrefix(contents string, quotePos int) bool {
	if quotePos == 0 {
		return false
	}
	prev, prevSize := utf8.DecodeLastRuneInString(contents[:quotePos])
	if prev != 'E' && prev != 'e' {
		return false
	}
	if quotePos == 1 {
		return true
	}
	return !isIdentifierContinuationRune(lastRuneBefore(contents, quotePos-prevSize))
}

// isIdentifierContinuationRune reports whether r can continue an unquoted
// PostgreSQL identifier. Continuation characters include `$`, so callers can
// distinguish `E'...'` from a longer token like `fooE'...'`.
func isIdentifierContinuationRune(r rune) bool {
	return r == '_' || r == '$' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func isLineBreak(b byte) bool {
	return b == '\n' || b == '\r'
}

// isHorizontalSpace reports whether b is a line-preserving indentation byte
// that can surround a top-level psql `\\` separator token.
func isHorizontalSpace(b byte) bool {
	return b == ' ' || b == '\t'
}

// matchDollarTagStart returns the opening dollar-quote delimiter at i when the
// `$tag$` sequence starts a real PostgreSQL dollar-quoted literal, not when it
// merely appears inside an ordinary identifier such as `foo$bar$baz`.
func matchDollarTagStart(contents string, i int) string {
	if contents[i] != '$' {
		return ""
	}
	if i > 0 && isIdentifierContinuationRune(lastRuneBefore(contents, i)) {
		return ""
	}
	if i+1 >= len(contents) {
		return ""
	}
	if contents[i+1] == '$' {
		return "$$"
	}
	r, size := utf8.DecodeRuneInString(contents[i+1:])
	if r == utf8.RuneError && size == 1 {
		return ""
	}
	if !isDollarTagStartRune(r) {
		return ""
	}
	tagEnd := i + 1 + size
	for tagEnd < len(contents) {
		r, size = utf8.DecodeRuneInString(contents[tagEnd:])
		if r == utf8.RuneError && size == 1 {
			return ""
		}
		if r == '$' {
			return contents[i : tagEnd+size]
		}
		if !isDollarTagRune(r) {
			return ""
		}
		tagEnd += size
	}
	return ""
}

// isDollarTagStartRune reports whether r is allowed as the first rune in a
// PostgreSQL dollar-quote tag, following unquoted identifier start rules.
func isDollarTagStartRune(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

// isDollarTagRune reports whether r is allowed after the first rune in a
// PostgreSQL dollar-quote tag, following unquoted identifier continuation
// rules but excluding `$`, which terminates the tag.
func isDollarTagRune(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

// lastRuneBefore returns the UTF-8 rune immediately preceding end in contents,
// or utf8.RuneError when there is no preceding rune.
func lastRuneBefore(contents string, end int) rune {
	if end <= 0 {
		return utf8.RuneError
	}
	r, _ := utf8.DecodeLastRuneInString(contents[:end])
	return r
}

// findMetaCommandSeparator returns the index of a `\\` separator token inside
// a top-level backslash directive line. The separator is recognized only when
// it is delimited by horizontal whitespace or line boundaries so ordinary
// doubled backslashes inside arguments are not mistaken for psql's
// meta-command/SQL separator.
//
// This rule is intentionally narrower than "any `\\` pair". In tested
// `psql` 17.9, glued forms such as `\x\\SELECT 1;`, `\echo hi\\SELECT 1;`, or
// `SELECT 1; \echo hi` are rejected with "invalid command \", and leading
// `\\`, `\\ SELECT 1;`, and `\\SELECT 1;` are also rejected when no preceding
// meta-command is present on the line. sqlc therefore treats `\\` as a
// separator only after an actual meta-command token on the same line.
func findMetaCommandSeparator(contents string, start, lineEnd int) int {
	for i := start; i+1 < lineEnd; i++ {
		if contents[i] != '\\' || contents[i+1] != '\\' {
			continue
		}
		if i <= start {
			continue
		}
		leftOK := i == start || isHorizontalSpace(contents[i-1])
		rightOK := i+2 == lineEnd || isHorizontalSpace(contents[i+2])
		if leftOK && rightOK {
			return i
		}
	}
	return -1
}

// readsPsqlCopyData reports whether the directive line is a `\copy ... from
// stdin` command that owns following data rows through an exact `\.` line.
func readsPsqlCopyData(contents string, start, lineEnd int) bool {
	fields := strings.Fields(strings.ToLower(contents[start:lineEnd]))
	if len(fields) == 0 || fields[0] != `\copy` {
		return false
	}
	for i := 1; i+1 < len(fields); i++ {
		if fields[i] == "from" && fields[i+1] == "stdin" {
			return true
		}
	}
	return false
}

// metaCommandToken returns the first backslash-command token from a top-level
// directive line, preserving psql's command-case distinctions and dropping any
// trailing arguments.
func metaCommandToken(line string) string {
	end := len(line)
	for i := 1; i < len(line); i++ {
		if isHorizontalSpace(line[i]) {
			end = i
			break
		}
	}
	return line[:end]
}

// isUnsupportedConditionalMetaCommand reports whether token is a psql
// conditional directive that sqlc currently rejects instead of flattening,
// because dropping only the meta-command lines would leave inactive branch SQL
// behind and change semantics.
func isUnsupportedConditionalMetaCommand(token string) bool {
	switch token {
	case `\if`, `\elif`, `\else`, `\endif`:
		return true
	default:
		return false
	}
}

// isSemanticMetaCommand reports whether token is a psql command whose effects
// are not reproduced by schema preprocessing even though the line is stripped.
// These commands can change connection context, include external SQL, execute
// generated SQL, or stream data, so callers should surface a warning when they
// are removed.
func isSemanticMetaCommand(token string) bool {
	switch token {
	case `\connect`, `\c`, `\i`, `\include`, `\ir`, `\include_relative`, `\copy`, `\gexec`:
		return true
	default:
		return false
	}
}

// schemaPreprocessWarnings renders stripping and approximation warnings in a
// stable order so callers see deterministic output.
func schemaPreprocessWarnings(tokens map[string]struct{}, warnedApproximateSessionSemantics bool) []string {
	if len(tokens) == 0 {
		if !warnedApproximateSessionSemantics {
			return nil
		}
		return []string{approximateSessionSemanticsWarning()}
	}

	order := []string{`\c`, `\connect`, `\i`, `\include`, `\ir`, `\include_relative`, `\copy`, `\gexec`}
	var warnings []string
	for _, token := range order {
		if _, ok := tokens[token]; !ok {
			continue
		}
		warnings = append(warnings, fmt.Sprintf("warning: stripped psql meta-command %s during schema preprocessing; sqlc does not execute its semantic effects", token))
	}
	if warnedApproximateSessionSemantics {
		warnings = append(warnings, approximateSessionSemanticsWarning())
	}
	return warnings
}

// approximateSessionSemanticsWarning explains that preprocessing intentionally
// does not emulate full psql session or transaction semantics.
func approximateSessionSemanticsWarning() string {
	return "warning: schema preprocessing only approximates psql session semantics after standard_conforming_strings or transaction-scoped script changes"
}

// stripPsqlCopyData removes the payload rows for a `\copy ... from stdin`
// command until the exact `\.` terminator line. It returns false if no
// terminator was found.
func stripPsqlCopyData(out *strings.Builder, contents string, lineEnd int) (int, bool) {
	i := writeLineEnding(out, contents, lineEnd)
	for i < len(contents) {
		dataLineEnd := i
		for dataLineEnd < len(contents) && contents[dataLineEnd] != '\r' && contents[dataLineEnd] != '\n' {
			dataLineEnd++
		}
		if contents[i:dataLineEnd] == `\.` {
			return writeLineEnding(out, contents, dataLineEnd), true
		}
		i = writeLineEnding(out, contents, dataLineEnd)
		if dataLineEnd == len(contents) {
			return dataLineEnd, false
		}
	}
	return i, false
}

// consumeStatementFragment appends normalized SQL text to the current
// statement buffer and optionally treats top-level semicolons as statement
// terminators for tracking best-effort standard_conforming_strings changes.
func consumeStatementFragment(text string, statement *strings.Builder, standardConformingStringsOff *bool, warnedApproximateSessionSemantics *bool, allowTerminator bool) {
	for i := 0; i < len(text); i++ {
		b := text[i]
		if b == '\r' || b == '\n' {
			statement.WriteByte(' ')
			continue
		}
		statement.WriteByte(b)
		if !allowTerminator || b != ';' {
			continue
		}
		applyStandardConformingStringsStatement(statement.String(), standardConformingStringsOff, warnedApproximateSessionSemantics)
		statement.Reset()
	}
}

// applyStandardConformingStringsStatement updates the best-effort
// standard_conforming_strings state from a completed top-level statement and
// flags constructs whose full psql session semantics are intentionally not
// modeled by preprocessing. Transaction-scoped and savepoint-scoped behavior
// is deliberately treated as an approximation rather than emulating psql's
// full session state machine.
func applyStandardConformingStringsStatement(stmt string, standardConformingStringsOff *bool, warnedApproximateSessionSemantics *bool) {
	if update, ok := parseStandardConformingStringsSetting(stmt); ok {
		if update.local || (update.value != nil && *update.value) {
			*warnedApproximateSessionSemantics = true
		}
		switch {
		case update.scopeDefault, update.value == nil:
			*standardConformingStringsOff = false
		default:
			*standardConformingStringsOff = *update.value
		}
		return
	}

	fields := strings.Fields(strings.NewReplacer(
		";", " ",
		"\n", " ",
		"\r", " ",
		"\t", " ",
	).Replace(strings.ToLower(stmt)))
	if len(fields) == 0 {
		return
	}
	switch fields[0] {
	case "start":
		if len(fields) > 1 && fields[1] == "transaction" {
			*warnedApproximateSessionSemantics = true
		}
	case "begin", "commit", "rollback", "abort", "end", "savepoint", "release":
		*warnedApproximateSessionSemantics = true
	case "reset":
		if len(fields) > 1 && fields[1] == "standard_conforming_strings" {
			*standardConformingStringsOff = false
			*warnedApproximateSessionSemantics = true
		}
		if len(fields) > 1 && fields[1] == "all" {
			*standardConformingStringsOff = false
			*warnedApproximateSessionSemantics = true
		}
	}
}

type standardConformingStringsUpdate struct {
	local        bool
	scopeDefault bool
	value        *bool
}

// parseStandardConformingStringsSetting extracts the standard_conforming_strings
// update encoded by a SET statement. The second return value reports whether
// stmt matched that setting at all.
func parseStandardConformingStringsSetting(stmt string) (standardConformingStringsUpdate, bool) {
	fields := strings.Fields(strings.NewReplacer(
		"=", " = ",
		";", " ",
		"\n", " ",
		"\r", " ",
		"\t", " ",
	).Replace(strings.ToLower(stmt)))
	if len(fields) == 0 || fields[0] != "set" {
		return standardConformingStringsUpdate{}, false
	}
	idx := 1
	update := standardConformingStringsUpdate{}
	if idx < len(fields) && (fields[idx] == "local" || fields[idx] == "session") {
		update.local = fields[idx] == "local"
		idx++
	}
	if idx+2 >= len(fields) || fields[idx] != "standard_conforming_strings" {
		return standardConformingStringsUpdate{}, false
	}
	if fields[idx+1] != "=" && fields[idx+1] != "to" {
		return standardConformingStringsUpdate{}, false
	}
	value := strings.Trim(fields[idx+2], `'`)
	switch value {
	case "off", "false":
		v := true
		update.value = &v
		return update, true
	case "on", "true":
		v := false
		update.value = &v
		return update, true
	case "default":
		if update.local {
			update.value = nil
		} else {
			update.scopeDefault = true
		}
		return update, true
	default:
		return standardConformingStringsUpdate{}, false
	}
}

// writeLineEnding copies the exact line terminator sequence at i, preserving
// either LF or CRLF. Bare CR is normalized to LF so downstream line counting,
// which keys off `\n`, stays consistent with stripped SQL.
func writeLineEnding(out *strings.Builder, contents string, i int) int {
	if i >= len(contents) {
		return i
	}
	if contents[i] == '\r' {
		i++
		if i < len(contents) && contents[i] == '\n' {
			out.WriteString("\r\n")
			i++
		} else {
			out.WriteByte('\n')
		}
		return i
	}
	if contents[i] == '\n' {
		out.WriteByte('\n')
		return i + 1
	}
	return i
}

// IsDown reports whether filename is a golang-migrate rollback migration.
func IsDown(filename string) bool {
	// Remove golang-migrate rollback files.
	return strings.HasSuffix(filename, ".down.sql")
}
