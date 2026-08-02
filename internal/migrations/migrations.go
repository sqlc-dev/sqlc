package migrations

import (
	"bufio"
	"regexp"
	"strings"
)

// psqlMetaCommand matches a psql meta-command (a line that begins with a
// backslash followed by a command name). pg_dump emits these starting with
// PostgreSQL 17.6 / 16.10 / 15.14 / 14.19 / 13.22 (e.g. `\restrict KEY` and
// `\unrestrict KEY`), and sqlc's SQL parsers cannot handle them.
var psqlMetaCommand = regexp.MustCompile(`^\\[A-Za-z!?;][^\n]*$`)

// Remove all lines after a rollback comment.
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

// RemovePsqlMetaCommands strips psql meta-command lines (e.g. `\restrict KEY`,
// `\unrestrict KEY`, `\connect foo`) from SQL input. These are emitted by
// pg_dump but are not valid SQL, so they must be removed before parsing.
func RemovePsqlMetaCommands(contents string) string {
	s := bufio.NewScanner(strings.NewReader(contents))
	var lines []string
	for s.Scan() {
		line := s.Text()
		if psqlMetaCommand.MatchString(line) {
			continue
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func IsDown(filename string) bool {
	// Remove golang-migrate rollback files.
	return strings.HasSuffix(filename, ".down.sql")
}
