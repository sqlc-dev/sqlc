package migrations

import (
	"bufio"
	"strings"
)

// Remove all lines that should be ignored by sqlc, such as rollback
// comments or explicit "sqlc:ignore" lines.
//
// goose:       -- +goose Down
// sql-migrate: -- +migrate Down
// tern:        ---- create above / drop below ----
// dbmate:      -- migrate:down
// generic:     `-- sqlc:ignore` until `-- sqlc:ignore end`
func RemoveIgnoredStatements(contents string) string {
	s := bufio.NewScanner(strings.NewReader(contents))
	var lines []string
	var ignoring bool
	for s.Scan() {
		line := s.Text()

		if strings.HasPrefix(line, "-- +goose Down") {
			break
		}
		if strings.HasPrefix(line, "-- +migrate Down") {
			break
		}
		if strings.HasPrefix(line, "---- create above / drop below ----") {
			break
		}
		if strings.HasPrefix(line, "-- migrate:down") {
			break
		}

		if strings.HasPrefix(line, "-- sqlc:ignore end") {
			ignoring = false
			// no need to keep this line in result
			line = ""
		} else if strings.HasPrefix(line, "-- sqlc:ignore") {
			ignoring = true
		}

		if ignoring {
			// make this line empty, so that errors are still reported on the
			// correct line
			line = ""
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func IsDown(filename string) bool {
	// Remove golang-migrate rollback files.
	return strings.HasSuffix(filename, ".down.sql")
}
